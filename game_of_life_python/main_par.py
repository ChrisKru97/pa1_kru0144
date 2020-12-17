from os import popen
from sys import argv
import numpy as np
from time import sleep
from multiprocessing.sharedctypes import Array
from threading import Condition, Thread

thread_count = 1
if(len(argv) > 1):
    thread_count = int(argv[1])


class Field:
    def __init__(self, height, width):
        self.h = height
        self.w = width
        self.data = Array('b', [False] * width * height)
        self.next = Array('b', [False] * width * height)

    def set(self, x, y, value, use_next=False):
        if use_next:
            self.next[y * self.w + x] = value
        else:
            self.data[y * self.w + x] = value

    def set_next(self, x, y, value):
        self.set(x, y, value, True)

    def prepare_switch(self):
        for i in range(self.h * self.w):
            self.next[i] = self.data[i]

    def switch_to_next(self):
        self.next, self.data = self.data, self.next

    def is_alive(self, x, y):
        alive_around = 0
        for i in range(-1 if x > 0 else 0, 2 if x < self.w - 1 else 1):
            for j in range(-1 if y > 0 else 0, 2 if y < self.h - 1 else 1):
                if(not (i == 0 and j == 0) and self.data[(y + j) * self.w + x + i]):
                    alive_around += 1
        return alive_around == 3 or (self.data[y * self.w + x] and alive_around == 2)

    def __str__(self):
        arr = np.where(self.data, 'O', ' ').reshape(
            (self.h, self.w)).transpose()
        endlines = np.array([np.full(self.h, '\n')])
        result = np.append(arr, endlines, axis=0)
        return ''.join(result.transpose().flatten())


class Life:
    def __init__(self, height, width):
        self.h = height
        self.w = width
        self.field = Field(self.h, self.w)
        self.cv = Condition()
        self.count = thread_count
        for _ in range(int(self.h * self.w / 5)):
            x, y = np.random.uniform(size=2)
            x *= self.w
            y *= self.h
            self.field.set(int(x), int(y), True)

    def start(self):
        while True:
            print(self.field)

            self.field.prepare_switch()

            self.cv.acquire()
            self.count = thread_count
            self.cv.release()

            for i in range(thread_count):
                Thread(target=update, args=(i, self)).start()

            self.cv.acquire()
            while self.count > 0:
                self.cv.wait()
            self.cv.release()

            self.field.switch_to_next()


def update(i, life):
    for y in range(i, life.h, thread_count):
        for x in range(life.w):
            life.field.set_next(x, y, life.field.is_alive(x, y))
    life.cv.acquire()
    life.count -= 1
    if life.count == 0:
        life.cv.notify_all()
    life.cv.release()


if __name__ == "__main__":
    life = Life(*(int(i) for i in popen('stty size', 'r').read().split()))
    life.start()
