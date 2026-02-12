from unsloth import FastVisionModel
import torch
from datasets import load_from_disk
import re
from trl import GRPOConfig, GRPOTrainer

max_seq_length = 16384 # Must be this long for VLMs
lora_rank = 32 # Larger rank = smarter, but slower

model, tokenizer = FastVisionModel.from_pretrained(
    model_name = "unsloth/Qwen3-VL-8B-Instruct-unsloth-bnb-4bit",
    max_seq_length = max_seq_length,
    load_in_4bit = True, # False for LoRA 16bit
    fast_inference = False, # Enable vLLM fast inference
    gpu_memory_utilization = 0.8, # Reduce if out of memory
)

model = FastVisionModel.get_peft_model(
    model,
    finetune_vision_layers     = False, # False if not finetuning vision layers
    finetune_language_layers   = True,  # False if not finetuning language layers
    finetune_attention_modules = True,  # False if not finetuning attention layers
    finetune_mlp_modules       = True,  # False if not finetuning MLP layers

    r = lora_rank,           # The larger, the higher the accuracy, but might overfit
    lora_alpha = 2*lora_rank,  # Recommended alpha == r at least
    lora_dropout = 0,
    bias = "none",
    random_state = 3407,
    use_rslora = False,  # We support rank stabilized LoRA
    loftq_config = None, # And LoftQ
    use_gradient_checkpointing = "unsloth", # Reduces memory usage
    # target_modules = "all-linear", # Optional now! Can specify a list if needed
)

dataset = load_from_disk("train_dataset_reasoning_0")

def resize_images(example):
    image = example["image"]
    image = image.resize((512, 512))
    example["image"] = image
    return example

dataset = dataset.map(resize_images)

def convert_to_rgb(example):
    image = example["image"]
    if image.mode != "RGB":
        image = image.convert("RGB")
    example["image"] = image
    return example

train_dataset = dataset.map(convert_to_rgb)

REASONING_START = "<REASONING>"
REASONING_END = "</REASONING>"
SOLUTION_START = "<SOLUTION>"
SOLUTION_END = "</SOLUTION>"

def formatting_reward_func(completions,**kwargs):
    import re
    thinking_pattern = f'{REASONING_START}(.*?){REASONING_END}'
    answer_pattern = f'{SOLUTION_START}(.*?){SOLUTION_END}'

    scores = []
    for completion in completions:
        if isinstance(completion, list):
            completion = completion[0]["content"] if completion else ""
        score = 0
        thinking_matches = re.findall(thinking_pattern, completion, re.DOTALL)
        answer_matches = re.findall(answer_pattern, completion, re.DOTALL)
        if len(thinking_matches) == 1:
            score += 1.0
        if len(answer_matches) == 1:
            score += 1.0

        # Fix up addCriterion issues
        # See https://unsloth.ai/docs/new/vision-reinforcement-learning-vlm-rl#qwen-2.5-vl-vision-rl-issues-and-quirks
        # Penalize on excessive addCriterion and newlines
        if len(completion) != 0:
            removal = completion.replace("addCriterion", "").replace("\n", "")
            if (len(completion)-len(removal))/len(completion) >= 0.5:
                score -= 2.0

        scores.append(score)
    return scores


def correctness_reward_func(prompts, completions, answer, **kwargs) -> list[float]:
    answer_pattern = f'{SOLUTION_START}(.*?){SOLUTION_END}'
    
    # Extract content if completion is a list of dicts (standard for some HF models)
    completions = [(c[0]["content"] if c else "") if isinstance(c, list) else c for c in completions]
    
    # Find all matches for the solution pattern
    responses = [re.findall(answer_pattern, completion, re.DOTALL) for completion in completions]
    
    scores = []
    
    # Debug print (optional)
    # print("Responses found:", responses)

    for r, a in zip(responses, answer):
        # r is a list of strings found by regex, a is the ground truth
        
        # Check if we found exactly one solution tag
        if len(r) == 1 and r[0].strip():
            try:
                # Parse the prediction
                pred_val = int(r[0].replace('\n', '').strip())
                
                # Parse the ground truth (handle cases like "25.0")
                true_val = int(float(str(a)))
                
                # Calculate score
                # This logic assumes the score decreases as distance from truth increases
                # Adjust / 5 logic if you want the score to be between specific bounds (e.g. 0 to 1)
                # Currently: Exact match = 5.0, Distance of 5 = 4.0, etc.
                score = (25 - abs(true_val - pred_val)) / 5
                
            except ValueError:
                # Handle cases where the text inside tags isn't a valid number
                score = 0.0
        else:
            # Case where tags are missing or multiple tags found
            score = 0.0

        scores.append(score)

    return scores

training_args = GRPOConfig(
    learning_rate = 5e-6,
    adam_beta1 = 0.9,
    adam_beta2 = 0.99,
    weight_decay = 0.1,
    warmup_ratio = 0.1,
    lr_scheduler_type = "cosine",
    optim = "adamw_8bit",
    logging_steps = 1,
    log_completions = False,
    per_device_train_batch_size = 1,
    gradient_accumulation_steps = 1, # Increase to 4 for smoother training
    num_generations = 8, # Decrease if out of memory
    max_prompt_length = 1024,
    max_completion_length = 1024,
    num_train_epochs = 1, # Set to 1 for a full training run
    save_steps = 10,
    # max_steps=40,
    max_grad_norm = 0.1,
#    num_train_epochs=1,
    report_to = "wandb", # Can use Weights & Biases
    output_dir = "outputs",

    # Below enables GSPO:
    importance_sampling_level = "sequence",
    mask_truncated_completions = False,
    loss_type = "dr_grpo",
)

trainer = GRPOTrainer(
    model = model,
    args = training_args,
    # Pass the processor to handle multimodal inputs
    processing_class = tokenizer,
    reward_funcs = [
        formatting_reward_func,
        correctness_reward_func,
    ],
    train_dataset = train_dataset,
)

trainer.train()

# Merge to 16bit
if True: model.save_pretrained_merged("model", tokenizer, save_method = "merged_16bit",)
if False: model.push_to_hub_merged("hf/model", tokenizer, save_method = "merged_16bit", token = "")

# Merge to 4bit
if True: model.save_pretrained_merged("model", tokenizer, save_method = "merged_4bit",)
if False: model.push_to_hub_merged("hf/model", tokenizer, save_method = "merged_4bit", token = "")

# Just LoRA adapters
if True:
    model.save_pretrained("model")
    tokenizer.save_pretrained("model")
if False:
    model.push_to_hub("hf/model", token = "")
    tokenizer.push_to_hub("hf/model", token = "")

trainer.train()
