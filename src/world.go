package main

import (
	"fmt"
	"strconv"
	"strings"
)

const ALIVE byte = 0xFF
const DEAD byte = 0x00

type world struct {
	width  int
	height int
	matrix []byte
}

// Create a new empty world
func newWorld(height int, width int) world {
	matrix := make([]byte, height*width)
	world := world{
		width:  width,
		height: height,
		matrix: matrix,
	}
	return world
}

func (w *world) getCell(row int, col int) byte {
	return w.matrix[(row*w.width)+col]
}

func (w *world) setCell(row int, col int, val byte) {
	w.matrix[(row*w.width)+col] = val
}

func (w *world) iterate(closure func(row int, col int, value byte)) {
	for row := 0; row < w.height; row++ {
		for col := 0; col < w.width; col++ {
			closure(row, col, w.getCell(row, col))
		}
	}
}

// Create a clone of a world (separate memory)
func (w *world) clone() world {
	newW := newWorld(w.height, w.width)
	copy(newW.matrix, w.matrix)
	return newW
}

// Load a world from a PGM
func loadWorldFromPgm(height int, width int, d distributorChans) world {
	w := newWorld(height, width)

	// Request the io goroutine to read in the image with the given filename.
	d.io.command <- ioInput
	d.io.filename <- strings.Join([]string{strconv.Itoa(w.width), strconv.Itoa(w.height)}, "x")

	// The io goroutine sends the requested image byte by byte, in rows.
	w.iterate(func(y int, x int, _ byte) {
		val := <-d.io.inputVal
		if val != 0 {
			//fmt.Println("Alive cell at", x, y)
			w.setCell(y, x, val)
		}
	})

	return w
}

// Save a world to a PGM
// The turn number is appended to the filename
func (w *world) saveToPgm(d distributorChans, turn int) {
	d.io.command <- ioOutput
	size := strings.Join([]string{strconv.Itoa(w.width), strconv.Itoa(w.height)}, "x")
	d.io.filename <- fmt.Sprintf("%s_turn_%d", size, turn)
	w.iterate(func(y int, x int, value byte) {
		d.io.outputVal <- value
	})
}

func (w *world) calculateAlive() []cell {
	// Create an empty slice to store coordinates of cells that are still alive after p.turns are done.
	var finalAlive []cell
	// Go through the w and append the cells that are still alive.
	w.iterate(func(y int, x int, value byte) {
		if value != 0 {
			finalAlive = append(finalAlive, cell{x: x, y: y})
		}
	})
	return finalAlive
}
