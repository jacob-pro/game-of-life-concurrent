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
	matrix [][]byte
}

func newWorld(height int, width int) world {
	matrix := make([][]byte, height)
	for i := range matrix {
		matrix[i] = make([]byte, width)
	}
	world := world{
		width:  width,
		height: height,
		matrix: matrix,
	}
	return world
}

func loadWorldFromPGM(world *world, d distributorChans) {
	// Request the io goroutine to read in the image with the given filename.
	d.io.command <- ioInput
	d.io.filename <- strings.Join([]string{strconv.Itoa(world.width), strconv.Itoa(world.height)}, "x")

	// The io goroutine sends the requested image byte by byte, in rows.
	for y := 0; y < world.height; y++ {
		for x := 0; x < world.width; x++ {
			val := <-d.io.inputVal
			if val != 0 {
				fmt.Println("Alive cell at", x, y)
				world.matrix[y][x] = val
			}
		}
	}
}

func cloneWorld(original *world) world {
	newW := newWorld(original.height, original.width)
	for y := 0; y < original.height; y++ {
		for x := 0; x < original.width; x++ {
			newW.matrix[y][x] = original.matrix[y][x]
		}
	}
	return newW
}

func customMod(index int, max int) int {
	if index >= max {
		return index - max
	} else if index < 0 {
		return index + max
	} else {
		return index
	}
}

func updateWorldState(original *world) world {
	clone := cloneWorld(original)
	for y := 0; y < original.height; y++ {
		for x := 0; x < original.width; x++ {
			var neighboursAlive = 0
			//Count alive cells in 3x3 grid
			for i := y - 1; i <= y+1; i++ {
				for j := x - 1; j <= x+1; j++ {
					if original.matrix[customMod(i, original.height)][customMod(j, original.width)] == ALIVE {
						neighboursAlive++
					}
				}
			}
			if original.matrix[y][x] == ALIVE {
				neighboursAlive--
				if neighboursAlive == 2 || neighboursAlive == 3 {
					clone.matrix[y][x] = ALIVE
				} else {
					clone.matrix[y][x] = DEAD
				}
			} else if neighboursAlive == 3 {
				clone.matrix[y][x] = ALIVE
			}
		}
	}
	return clone
}

func calculateAlive(world *world) []cell {
	// Create an empty slice to store coordinates of cells that are still alive after p.turns are done.
	var finalAlive []cell
	// Go through the world and append the cells that are still alive.
	for y := 0; y < world.height; y++ {
		for x := 0; x < world.width; x++ {
			if world.matrix[y][x] != 0 {
				finalAlive = append(finalAlive, cell{x: x, y: y})
			}
		}
	}
	return finalAlive
}

// distributor divides the work between workers and interacts with other goroutines.
func distributor(p golParams, d distributorChans, alive chan []cell) {

	var world = newWorld(p.imageHeight, p.imageWidth)
	loadWorldFromPGM(&world, d)

	// Calculate the new state of Game of Life after the given number of turns.
	for turns := 0; turns < p.turns; turns++ {
		world = updateWorldState(&world)
	}

	// Make sure that the Io has finished any output before exiting.
	d.io.command <- ioCheckIdle
	<-d.io.idle

	// Return the coordinates of cells that are still alive.
	alive <- calculateAlive(&world)
}
