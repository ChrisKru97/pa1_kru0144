#include <stdio.h>
#include <stdlib.h>
#include <stdbool.h>
#include <sys/ioctl.h>
#include <unistd.h>
#include <time.h>
#include <omp.h>

#define for_i_h for (int i = 0; i < h; i++)
#define for_j_w for (int j = 0; j < w; j++)

const char life = 'O';
const char death = ' ';

void fillRandomly(bool **field, int w, int h, double ratio)
{
    srand(time(NULL));

    int x_rand, y_rand;
    int life_count = (int)(w * h * ratio);

#pragma omp parallel for private(x_rand, y_rand)
    for (int i = 0; i < life_count; i++)
    {
        x_rand = (int)(w * 1.0 * rand() / RAND_MAX);
        y_rand = (int)(h * 1.0 * rand() / RAND_MAX);
        field[y_rand][x_rand] = true;
    }
}

bool **generateField(int w, int h)
{
    bool **arr;
    arr = malloc(w * h * sizeof(bool));

#pragma omp parallel for
    for_i_h
    {
        arr[i] = malloc(w * sizeof(bool));
    }

    return arr;
}

void printField(bool **field, int w, int h)
{
    for_i_h
    {
        putc('\n', stdout);
        for_j_w
        {
            putc(field[i][j] ? life : death, stdout);
        }
    }
    fflush(stdout);
}

void copy(bool **source, bool **target, int w, int h)
{
#pragma omp parallel for collapse(2)
    for_i_h
    {
        for_j_w
        {
            target[i][j] = source[i][j];
        }
    }
}

bool checkAlive(bool **field, int x, int y, int w, int h)
{
    int alive = 0;
    int x_min = x > 0 ? -1 : 0, x_max = x < (w - 1) ? 1 : 0;
    int y_min = y > 0 ? -1 : 0, y_max = y < (h - 1) ? 1 : 0;

#pragma omp parallel for collapse(2) reduction(+ \
                                               : alive)
    for (int i = x_min; i <= x_max; i++)
    {
        for (int j = y_min; j <= y_max; j++)
        {
            if (!(i == 0 && j == 0) && field[y + j][x + i])
            {
                alive++;
            }
        }
    }
    return alive == 3 || (field[y][x] && alive == 2);
}

void processNextStep(bool **field, bool **temp_field, int w, int h)
{
#pragma omp parallel for collapse(2)
    for_i_h
    {
        for_j_w
        {
            temp_field[i][j] = checkAlive(field, j, i, w, h);
        }
    }
    copy(temp_field, field, w, h);
}

void freeFromMemory(bool **field, int w, int h)
{
    for_i_h
    {
        free(field[i]);
    }
    free(field);
}

int main()
{
    struct winsize w;
    ioctl(STDOUT_FILENO, TIOCGWINSZ, &w);

    bool **field = generateField(w.ws_col, w.ws_row);
    bool **temp_field = generateField(w.ws_col, w.ws_row);

    fillRandomly(field, w.ws_col, w.ws_row, 0.3);

    while (true)
    {
        printField(field, w.ws_col, w.ws_row);
        processNextStep(field, temp_field, w.ws_col, w.ws_row);
        usleep(66000);
    }
    freeFromMemory(field, w.ws_col, w.ws_row);
    freeFromMemory(temp_field, w.ws_col, w.ws_row);
    return 0;
}
