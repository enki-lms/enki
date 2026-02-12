#%%
import io
from importlib.readers import FileReader
from subprocess import TimeoutExpired

from PIL import Image
import pandas as pd
import requests
import os
import subprocess
from pypdf import PdfReader
import openai
import base64
from IPython.display import HTML
from io import BytesIO
import cairosvg
#%%
def convert_exercises_to_pdf(folder_path):
    files = os.listdir(folder_path)

    for file_name in files:
        _, file_ext = os.path.splitext(file_name)
        file_path = os.path.join(folder_path, file_name)

        if file_ext != ".pdf":
            subprocess.run(["libreoffice", "--headless", "--convert-to", "pdf", "--outdir", folder_path, file_path])


convert_exercises_to_pdf("/home/optimuseprime/Downloads/zadaci")

#%%
df = pd.read_csv("./data/submissions.csv")

PATCH_CODE = """
import turtle

def patch_turtle_speed():
    # Get the main screen.
    screen = turtle.Screen()
    # Disable animation refreshes: 0 means “don’t animate”, so that drawing happens as fast as possible.
    screen.tracer(0, 0)
    # Set the delay (in milliseconds) to 0.
    screen.delay(0)

    # Define a replacement for the speed method that forces fastest drawing.
    def always_fast(self, speed=None):
        # If called with no argument, return the current speed—which we always want to be fastest.
        # If an argument is given, ignore it and return 0.
        return 0

    # Patch the Turtle class speed methods so that any attempt to change speed is ignored.
    turtle.Turtle.speed = always_fast
    # Some code might use these alternative names, so patch them too.
    # turtle.Turtle.setspeed = always_fast
    # turtle.Turtle.setSpeed = always_fast

    # Also patch the module-level speed and delay functions.
    turtle.speed = always_fast
    turtle.delay = lambda *args, **kwargs: 0
    turtle.tracer = lambda n=0, delay=0: 0

def patch_turtle_loop_functions():
    # Save the original functions.
    original_mainloop = turtle.mainloop
    original_done = turtle.done
    original_exitonclick = turtle.exitonclick

    def patched_mainloop(*args, **kwargs):
        turtle.update()
        return original_mainloop(*args, **kwargs)

    def patched_done(*args, **kwargs):
        turtle.update()
        return original_done(*args, **kwargs)

    def patched_exitonclick(*args, **kwargs):
        turtle.update()
        return original_exitonclick(*args, **kwargs)

    # Patch the functions in the turtle module.
    turtle.mainloop = patched_mainloop
    turtle.done = patched_done
    turtle.exitonclick = patched_exitonclick


# Apply the patch immediately.
patch_turtle_speed()
patch_turtle_loop_functions()

del patch_turtle_speed, patch_turtle_loop_functions

screen = turtle.Screen()
screen.setup(1000, 1000)
""".strip()

SAVE_CODE = """
import canvasvg

canvasvg.saveall("/home/optimuseprime/Projects/enki/apps/ml/dataset/tmp_1/temp_output.svg", screen.getcanvas())
""".strip()

FOLDER_PATH = "/home/optimuseprime/Downloads"


def svg_to_high_res_png(svg_path, scale_factor=1):
    png_data = cairosvg.svg2png(url=svg_path, output_width=512,
                                output_height=512)

    img = Image.open(io.BytesIO(png_data))

    return img


def encode_pdf_to_base64(pdf_path):
    with open(pdf_path, "rb") as pdf_file:
        return base64.b64encode(pdf_file.read()).decode('utf-8')


def generate_solution_examples(row):
    input_file = os.path.join(FOLDER_PATH, "zadaci", "TestInput", f"z{row['ProblemID']}_0_turtle_input.txt")
    if os.path.isfile(input_file):
        return

    client = openai.OpenAI(
        base_url="https://openrouter.ai/api/v1",
        api_key="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"
    )

    f_org_path = row["path"]
    f_name, ext = os.path.splitext(f_org_path)
    if ext != ".pdf":
        f_org_path = f_name + ".pdf"

    f_path = os.path.join(FOLDER_PATH, f_org_path)

    completion = client.chat.completions.create(
        model="google/gemini-3-flash-preview",
        reasoning_effort="high",
        messages=[
            {
                "role": "user",
                "content": [
                    {
                        "type": "text",
                        "text": """Generate an appropriate input example that would adequately test whether the solution code as described in the exercise specification is correct.
                            Output only the input example, with each part of the input on a new line and nothing else. Attempt to generate the example such that the final drawing is quite large, e.g. fits well on a 1000x1000 screen."""
                    },
                    {
                        "type": "image_url",
                        "image_url": {
                            "url": f"data:application/pdf;base64,{encode_pdf_to_base64(f_path)}"
                        }
                    }
                ]

            }
        ]
    )

    input_ex = completion.choices[0].message.content

    with open(input_file, "w") as f:
        f.write(input_ex)


def generate_turtle_images(row):
    print(f"Processing problem {row['ProblemID']} for submission {row['SubmissionID']}")
    file_name_url = row["FileName"]

    resp = requests.get(file_name_url)
    file_name = os.path.basename(file_name_url)

    original_code = resp.text

    file_path = os.path.join("./tmp_1", file_name)

    with open(file_path, "w") as f:
        f.write(PATCH_CODE + "\n")
        f.write(original_code.replace("import turtle", "") + "\n")
        f.write(SAVE_CODE)

    input_file = os.path.join(FOLDER_PATH, "zadaci", "TestInput", f"z{row['ProblemID']}_0_turtle_input.txt")

    if os.path.exists("/home/optimuseprime/Projects/enki/apps/ml/dataset/tmp_1/temp_output.svg"):
        os.remove("/home/optimuseprime/Projects/enki/apps/ml/dataset/tmp_1/temp_output.svg")

    with open(input_file, "r") as f:
        try:
            out = subprocess.run(
                ["xvfb-run", "-a", "python3", file_path],
                input=f.read(),
                text=True,
                timeout=5,
                capture_output=True,
            )

            if out.returncode != 0:
                return None

            img = svg_to_high_res_png("/home/optimuseprime/Projects/enki/apps/ml/dataset/tmp_1/temp_output.svg")
            img.load()

            print(f"Processed problem {row['ProblemID']} for submission {row['SubmissionID']}")

            return img

        except subprocess.TimeoutExpired:
            print("Submission timed out")


def image_formatter(im):
    # Convert PIL image to Base64 string
    with BytesIO() as buffer:
        im.save(buffer, format='PNG')
        return f'<img src="data:image/png;base64,{base64.b64encode(buffer.getvalue()).decode()}" width="100"/>'


df["Image"] = df.apply(generate_turtle_images, axis=1)
# df.apply(generate_solution_examples, axis=1)
#%%
import pymupdf4llm
from datasets import Dataset
from datasets import Dataset, Features, Image, Value, Sequence

def convert_to_conversation(sample):
    if sample["Image"] is None:
        return None

    file_name, _ = os.path.splitext(sample["path"])

    problem_path = os.path.join(FOLDER_PATH, file_name + ".pdf")
    problem_text = pymupdf4llm.to_markdown(problem_path)

    file_name_url = sample["FileName"]

    resp = requests.get(file_name_url)

    code = resp.text

    instruction = f"""
    Your task is to grade a student's code for drawing a particular Turtle image as described in the problem text.
    You have been given the student's Python code, the drawing that the code creates, and the problem text.
    You must grade the student's solution on a scale from 0 to 25 in intervals of 5. Output ONLY the final score and NOTHING ELSE.

    <problem_text>
    {problem_text}
    </problem_text>

    <solution_code>
    {code}
    </solution_code>
    """.strip()

    img_bytes = None
    if sample["Image"] is not None:
        with io.BytesIO() as b:
            # Convert to PNG (or JPEG) bytes
            sample["Image"].save(b, format='PNG')
            img_bytes = b.getvalue()

    conversation = [
        {
            "role": "user",
            "content": [
                {
                    "type": "text",
                    "text": instruction,
                },
                {
                    "type": "image",
                    "image": img_bytes,
                }
            ]
        },
        {
            "role": "assistant",
            "content": [
                {
                    "type": "text",
                    "text": str(sample["Score"]),
                    "image": None,
                }
            ]
        }

    ]

    return conversation

features = Features({
    "messages": Sequence({
        "role": Value("string"),
        "content": Sequence({
            "type": Value("string"),
            "text": Value("string"),
            "image": Image()
        })
    })
})

series = df.apply(convert_to_conversation, axis=1)
df_result = series.to_frame(name="messages")
train_dataset = Dataset.from_pandas(df_result)
train_dataset.save_to_disk("train_dataset_1")
#%%
train_dataset = train_dataset.filter(lambda row: row['messages'] is not None)
train_dataset.save_to_disk("train_dataset_1")