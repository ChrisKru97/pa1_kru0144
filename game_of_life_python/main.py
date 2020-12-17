from os import popen
from sys import argv
import numpy as np
from time import sleep


class Field:
    def __init__(self, height, width):
        self.h = height
        self.w = width
        self.data = np.zeros((height, width), dtype=np.bool)

    def set(self, x, y, value, use_next=False):
        if use_next:
            self.next[y, x] = value
        else:
            self.data[y, x] = value

    def set_next(self, x, y, value):
        self.set(x, y, value, True)

    def prepare_switch(self):
        self.next = np.copy(self.data)

    def switch_to_next(self):
        self.data = self.next

    def is_alive(self, x, y):
        alive_around = 0
        for i in range(-1 if x > 0 else 0, 2 if x < self.w - 1 else 1):
            for j in range(-1 if y > 0 else 0, 2 if y < self.h - 1 else 1):
                if(not (i == 0 and j == 0) and self.data[y + j, x + i]):
                    alive_around += 1
        return alive_around == 3 or (self.data[y, x] and alive_around == 2)

    def __str__(self):
        arr = np.where(self.data, 'O', ' ').transpose()
        endlines = np.array([np.full(self.h, '\n')])
        result = np.append(arr, endlines, axis=0)
        return ''.join(result.transpose().flatten())


class Life:
    def __init__(self, height, width):
        self.h = int(height)
        self.w = int(width)
        self.field = Field(self.h, self.w)
        for _ in range(round(self.h * self.w / 5)):
            x, y = np.random.rand(2)
            x *= self.w
            y *= self.h
            self.field.set(int(x), int(y), True)

    def update(self):
        self.field.prepare_switch()
        for y in range(self.h):
            for x in range(self.w):
                self.field.set_next(x, y, self.field.is_alive(x, y))
        self.field.switch_to_next()

    def start(self):
        while True:
            print(self.field)
            self.update()
            # sleep(.33)


if __name__ == "__main__":
    life = Life(*popen('stty size', 'r').read().split())
    life.start()
