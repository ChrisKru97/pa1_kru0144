package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"
)

var updateWg sync.WaitGroup
var prepareWg sync.WaitGroup

// Field : used to store Life data
type Field struct {
	w, h int
	data [][]bool
	next [][]bool
}

// NewField : Creates new Field
func NewField(width, height int) *Field {
	data := make([][]bool, height)
	next := make([][]bool, height)
	for i := 0; i < height; i++ {
		data[i] = make([]bool, width)
		next[i] = make([]bool, width)
	}
	return &Field{data: data, next: next, w: width, h: height}
}

// Set : Sets dead/alive cell at x/y
func (field *Field) Set(x, y int, alive bool, useNext bool) {
	if useNext {
		field.next[y][x] = alive
	} else {
		field.data[y][x] = alive
	}
}

// PrepareSwitch : Copy field to prepare for processing next step
func (field *Field) PrepareSwitch(index, threadCount int) {
	defer prepareWg.Done()
	for i := index; i < field.w; i += threadCount {
		for j := 0; j < field.h; j++ {
			field.Set(i, j, field.data[j][i], true)
		}
	}
}

// SwitchToNext : switch current and next data
func (field *Field) SwitchToNext() {
	field.next, field.data = field.data, field.next
}

// IsAlive : Processes next step and tells if is alive
func (field *Field) IsAlive(x, y int) bool {
	aliveAround := 0
	iMin := -1
	iMax := 1
	if x == 0 {
		iMin = 0
	}
	if x == field.w-1 {
		iMax = 0
	}
	jMin := -1
	jMax := 1
	if y == 0 {
		jMin = 0
	}
	if y == field.h-1 {
		jMax = 0
	}
	for i := iMin; i <= iMax; i++ {
		for j := jMin; j <= jMax; j++ {
			if !(j == 0 && i == 0) && field.data[y+j][x+i] {
				aliveAround++
			}
		}
	}
	return aliveAround == 3 || (field.data[y][x] && aliveAround == 2)
}

func (field *Field) String() string {
	var buffer bytes.Buffer
	cell := byte(' ')
	for j := 0; j < field.h; j++ {
		for i := 0; i < field.w; i++ {
			if field.data[j][i] {
				cell = 'O'
			} else {
				cell = ' '
			}
			buffer.WriteByte(cell)
		}
		buffer.WriteByte('\n')
	}
	return buffer.String()
}

// Life : Contains life
type Life struct {
	field             *Field
	threadCount, w, h int
}

// NewLife : Creates new life
func NewLife(width, height, threadCount int) *Life {
	field := NewField(width, height)
	for i := 0; i < (width * height / 5); i++ {
		field.Set(rand.Intn(width), rand.Intn(height), true, false)
	}
	return &Life{
		threadCount: threadCount,
		field:       field,
		w:           width,
		h:           height,
	}
}

// Update : process next step
func (life *Life) Update(index int) {
	defer updateWg.Done()
	life.field.PrepareSwitch(index, life.threadCount)
	prepareWg.Wait()
	for i := index; i < life.w; i += life.threadCount {
		for j := 0; j < life.h; j++ {
			life.field.Set(i, j, life.field.IsAlive(i, j), true)
		}
	}
}

// Start : start a life
func (life *Life) Start() {
	for true {
		fmt.Print("\x0c", life.field)
		updateWg.Add(life.threadCount)
		prepareWg.Add(life.threadCount)
		for i := 0; i < life.threadCount; i++ {
			go life.Update(i)
		}
		updateWg.Wait()
		life.field.SwitchToNext()
		time.Sleep(66000000)
	}
}

func main() {
	cmd := exec.Command("stty", "size")
	cmd.Stdin = os.Stdin
	out, _ := cmd.Output()
	size := strings.Split(strings.Trim(string(out), "\n"), " ")
	height, _ := strconv.Atoi(size[0])
	width, _ := strconv.Atoi(size[1])

	threadCount := 1
	if len(os.Args) > 1 {
		threadCount, _ = strconv.Atoi(os.Args[1])
	}

	life := NewLife(width, height, threadCount)
	life.Start()
}
