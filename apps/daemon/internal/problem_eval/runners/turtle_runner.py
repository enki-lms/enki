import sys
import json
import math
from dataclasses import dataclass
from typing import List, Tuple, Optional

# Mock the turtle module
class MockTurtle:
    def __init__(self):
        self._commands = []
        self._x = 0
        self._y = 0
        self._angle = 0 # 0 is East, 90 is North
        self._pen_down = True
        self._pen_color = "black"
        self._fill_color = "black"
        self._pen_size = 1
        self._is_filling = False
        self._fill_path = []
        self._speed = 3
        self._visible = True
        
        # Initial state
        self._record_command("init", {})

    def _record_command(self, cmd: str, args: dict):
        self._commands.append({
            "cmd": cmd,
            "args": args,
            "state": {
                "x": self._x,
                "y": self._y,
                "angle": self._angle,
                "pen_down": self._pen_down,
                "pen_color": self._pen_color,
                "fill_color": self._fill_color,
                "pen_size": self._pen_size
            }
        })
        if self._is_filling:
            self._fill_path.append((self._x, self._y))

    # Motion
    def forward(self, distance):
        rad = math.radians(self._angle)
        new_x = self._x + distance * math.cos(rad)
        new_y = self._y + distance * math.sin(rad)
        
        self._record_command("line", {"x1": self._x, "y1": self._y, "x2": new_x, "y2": new_y})
        
        self._x = new_x
        self._y = new_y
        if self._is_filling:
            self._fill_path.append((self._x, self._y))

    def fd(self, distance): self.forward(distance)

    def backward(self, distance):
        self.forward(-distance)
    
    def bk(self, distance): self.backward(distance)
    def back(self, distance): self.backward(distance)

    def right(self, angle):
        self._angle -= angle
        self._angle %= 360
    
    def rt(self, angle): self.right(angle)

    def left(self, angle):
        self._angle += angle
        self._angle %= 360
        
    def lt(self, angle): self.left(angle)

    def goto(self, x, y):
        if self._pen_down:
             self._record_command("line", {"x1": self._x, "y1": self._y, "x2": x, "y2": y})
        self._x = x
        self._y = y
        if self._is_filling:
            self._fill_path.append((self._x, self._y))

    def setpos(self, x, y): self.goto(x, y)
    def setposition(self, x, y): self.goto(x, y)

    def setx(self, x): self.goto(x, self._y)
    def sety(self, y): self.goto(self._x, y)

    def setheading(self, to_angle):
        self._angle = to_angle % 360
    
    def seth(self, to_angle): self.setheading(to_angle)

    def home(self):
        self.goto(0, 0)
        self.setheading(0)

    def circle(self, radius, extent=None, steps=None):
        if extent is None:
            extent = 360
        if steps is None:
            steps = int(math.fabs(extent)) # simplified steps calculation
            if steps == 0: steps = 1

        # Record circle command directly for better rendering, or approximate with lines
        # For simplicity in SVG rendering later, let's approximate with lines if needed, 
        # but storing as 'circle' command is better if we implement SVG arc handling.
        # However, standard turtle draws circles by moving. Let's approximate.
        
        # Actually, let's record it as a command so the renderer can handle it nicely
        # But wait, python turtle updates position after circle.
        # So we need to calculate end position.
        
        # Arc geometry
        # Center of circle is at (radius) to the left of turtle
        # This is complex to replicate exactly without full math.
        # Let's use the approximation loop which is what turtle docs say it effectively does.
        
        theta = extent / steps
        dist = 2 * radius * math.sin(math.radians(theta / 2))
        
        for _ in range(steps):
            self.left(theta / 2) # curve start
            self.forward(dist)
            self.left(theta / 2) # curve end

    def dot(self, size=None, color=None):
        if size is None:
            size = max(self._pen_size + 4, 2 * self._pen_size)
        
        self._record_command("dot", {"x": self._x, "y": self._y, "size": size, "color": color or self._pen_color})

    # Pen control
    def penup(self): self._pen_down = False
    def pu(self): self.penup()
    def up(self): self.penup()

    def pendown(self): self._pen_down = True
    def pd(self): self.pendown()
    def down(self): self.pendown()

    def pensize(self, width): self._pen_size = width
    def width(self, width): self.pensize(width)

    def pencolor(self, *args):
        if len(args) == 1:
            self._pen_color = args[0]
        elif len(args) == 3:
             self._pen_color = "#%02x%02x%02x" % (int(args[0]), int(args[1]), int(args[2]))
    
    def fillcolor(self, *args):
        if len(args) == 1:
            self._fill_color = args[0]
        elif len(args) == 3:
             self._fill_color = "#%02x%02x%02x" % (int(args[0]), int(args[1]), int(args[2]))

    def color(self, *args):
        if len(args) == 1:
            self.pencolor(args[0])
            self.fillcolor(args[0])
        elif len(args) == 2:
            self.pencolor(args[0])
            self.fillcolor(args[1])
            
    def begin_fill(self):
        self._is_filling = True
        self._fill_path = [(self._x, self._y)]
        
    def end_fill(self):
        self._is_filling = False
        if len(self._fill_path) > 2:
            self._record_command("fill", {"points": self._fill_path[:], "color": self._fill_color})

    # State
    def speed(self, speed): self._speed = speed
    def hideturtle(self): self._visible = False
    def ht(self): self.hideturtle()
    def showturtle(self): self._visible = True
    def st(self): self.showturtle()

    # Screen
    def Screen(self): return MockScreen()


class MockScreen:
    def __init__(self):
        self._bgcolor = "white"
        self._tracer = 1
        
    def bgcolor(self, color): self._bgcolor = color
    def tracer(self, n=None, delay=None): pass
    def update(self): pass
    def setup(self, width=0.5, height=0.75, startx=None, starty=None): pass
    def title(self, titlestring): pass
    def bye(self): pass
    def exitonclick(self): pass


# SVG Renderer
def commands_to_svg(commands, width=800, height=600):
    # Turtle coordinates: center is (0,0). SVG coordinates: top-left is (0,0).
    # We need to translate.
    cx, cy = width / 2, height / 2
    
    svg_elements = []
    
    # Background - usually white unless set
    bg_color = "white"
    
    # Process commands
    for cmd in commands:
        c = cmd["cmd"]
        args = cmd["args"]
        state = cmd["state"]
        
        if c == "line":
            if state["pen_down"]:
                x1, y1 = cx + args["x1"], cy - args["y1"]
                x2, y2 = cx + args["x2"], cy - args["y2"]
                color = state["pen_color"]
                width_px = state["pen_size"]
                svg_elements.append(f'<line x1="{x1}" y1="{y1}" x2="{x2}" y2="{y2}" stroke="{color}" stroke-width="{width_px}" stroke-linecap="round" />')
        
        elif c == "dot":
             x, y = cx + args["x"], cy - args["y"]
             r = args["size"] / 2
             color = args["color"]
             svg_elements.append(f'<circle cx="{x}" cy="{y}" r="{r}" fill="{color}" />')
             
        elif c == "fill":
            points = args["points"]
            color = args["color"]
            pts_str = " ".join([f"{cx + x},{cy - y}" for x, y in points])
            svg_elements.append(f'<polygon points="{pts_str}" fill="{color}" stroke="none" />')

    svg_content = "\n".join(svg_elements)
    return f'<svg xmlns="http://www.w3.org/2000/svg" width="{width}" height="{height}" viewBox="0 0 {width} {height}" style="background-color: {bg_color}">\n{svg_content}\n</svg>'

# Monkey patch turtle module
import types
turtle_mod = types.ModuleType("turtle")

# Create global instance
_t = MockTurtle()
_s = MockScreen()

# Expose methods to module
for attr in dir(_t):
    if not attr.startswith("_"):
        setattr(turtle_mod, attr, getattr(_t, attr))

turtle_mod.Screen = lambda: _s
turtle_mod.Turtle = lambda: _t
turtle_mod.bgcolor = _s.bgcolor
sys.modules["turtle"] = turtle_mod

# --- End of Mock Setup ---

# User code will be executed here
def run_user_code(user_code):
    try:
        exec(user_code, {"turtle": turtle_mod})
    except Exception as e:
        print(f"Error: {e}", file=sys.stderr)
        # We might want to still render what was drawn so far
        pass

if __name__ == "__main__":
    if len(sys.argv) > 1:
        # If code is passed as argument (or file)
        # For simplicity in this environment, let's assume code is read from stdin or a file
        # But wait, problem_eval/executor passes code as a string to python execution.
        # We should read from a file passed as arg.
        
        filename = sys.argv[1]
        with open(filename, 'r') as f:
            code = f.read()
        
        run_user_code(code)
        
        # After execution, generate SVG
        svg = commands_to_svg(_t._commands)
        
        # Print SVG to stdout with a delimiter effectively, or just print it.
        # The executor should capture stdout.
        # Let's wrap it in a JSON or just print a specific marker.
        # Actually base64 might be safer to avoid mixing with print outputs from user code.
        
        import base64
        svg_b64 = base64.b64encode(svg.encode('utf-8')).decode('utf-8')
        print(f"\n---TURTLE_RESULT_START---\n{svg_b64}\n---TURTLE_RESULT_END---")
    else:
        print("Usage: python3 turtle_runner.py <code_file>")
