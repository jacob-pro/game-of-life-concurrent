package main

import (
	"fmt"
	"strconv"
	"strings"
)

const ALIVE byte = 0xFF
const DEAD byte = 0x00

type World struct {
	width  int
	height int
	matrix [][]byte
}

// Create a new empty World
func NewWorld(height int, width int) World {
	matrix := make([][]byte, height)
	for i := range matrix {
		matrix[i] = make([]byte, width)
	}
	world := World{
		width:  width,
		height: height,
		matrix: matrix,
	}
	return world
}

func (w *World) Iterate(closure func(y int, x int, value byte)) {
	for y := 0; y < w.height; y++ {
		for x := 0; x < w.width; x++ {
			closure(y, x, w.matrix[y][x])
		}
	}
}

// Create a clone of a World (separate memory)
func (w *World) Clone() World {
	newW := NewWorld(w.height, w.width)
	w.Iterate(func(y int, x int, value byte) {
		newW.matrix[y][x] = value
	})
	return newW
}

// Load a World from a PGM
func LoadWorldFromPgm(height int, width int, d distributorChans) World {
	w := NewWorld(height, width)

	// Request the io goroutine to read in the image with the given filename.
	d.io.command <- ioInput
	d.io.filename <- strings.Join([]string{strconv.Itoa(w.width), strconv.Itoa(w.height)}, "x")

	// The io goroutine sends the requested image byte by byte, in rows.
	w.Iterate(func(y int, x int, _ byte) {
		val := <-d.io.inputVal
		if val != 0 {
			//fmt.Println("Alive cell at", x, y)
			w.matrix[y][x] = val
		}
	})

	return w
}

// Save a World to a PGM
// The turn number is appended to the filename
func (w *World) SaveToPgm(d distributorChans, turn int) {
	d.io.command <- ioOutput
	size := strings.Join([]string{strconv.Itoa(w.width), strconv.Itoa(w.height)}, "x")
	d.io.filename <- fmt.Sprintf("%s_turn_%d", size, turn)
	w.Iterate(func(y int, x int, value byte) {
		d.io.outputVal <- value
	})
}

func (w *World) CalculateAlive() []cell {
	// Create an empty slice to store coordinates of cells that are still alive after p.turns are done.
	var finalAlive []cell
	// Go through the w and append the cells that are still alive.
	w.Iterate(func(y int, x int, value byte) {
		if value != 0 {
			finalAlive = append(finalAlive, cell{x: x, y: y})
		}
	})
	return finalAlive
}
