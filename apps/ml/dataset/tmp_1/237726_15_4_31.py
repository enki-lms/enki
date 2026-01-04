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
from turtle import*
from random import*
a = int(input())
b = int(input())
colormode(255)
fillcolor('red')
begin_fill()
pd()
circle(((a**2+b**2)**0.5)//2)
end_fill()
pu()
lt(90)
fd(((a**2+b**2)**0.5)//2)
rt(180)
fd(b//2)
rt(90)
fd(a//2)
rt(180)
pd()
begin_fill()
for i in range(2):
    fd(a)
    fillcolor('blue')
    lt(90)
    fd(b)
    lt(90)
end_fill()

import canvasvg

canvasvg.saveall("/home/optimuseprime/Projects/enki/apps/ml/dataset/tmp_1/temp_output.svg", screen.getcanvas())