package main

import (
	"fmt"
	"sync"
	"time"
)

func customMod(index int, max int) int {
	if index >= max {
		return index - max
	} else if index < 0 {
		return index + max
	} else {
		return index
	}
}

func (w *World) computeStateUpdate() {
	clone := w.Clone()
	for y := 0; y < clone.height; y++ {
		for x := 0; x < clone.width; x++ {
			var neighboursAlive = 0
			//Count alive cells in 3x3 grid
			for i := y - 1; i <= y+1; i++ {
				for j := x - 1; j <= x+1; j++ {
					if clone.matrix[customMod(i, clone.height)][customMod(j, clone.width)] == ALIVE {
						neighboursAlive++
					}
				}
			}
			if clone.matrix[y][x] == ALIVE {
				neighboursAlive--
				if neighboursAlive == 2 || neighboursAlive == 3 {
					w.matrix[y][x] = ALIVE
				} else {
					w.matrix[y][x] = DEAD
				}
			} else if neighboursAlive == 3 {
				w.matrix[y][x] = ALIVE
			}
		}
	}
}

type distributorState struct {
	currentTurn  int
	currentWorld World
	mux          sync.Mutex
}

func handleKeyboard(state *distributorState, d distributorChans) {
	for {
		keyPress := <-d.keyChan
		switch keyPress {
		case 's':
			state.mux.Lock()
			state.currentWorld.SaveToPgm(d, state.currentTurn)
			state.mux.Unlock()
		case 'p':
			state.mux.Lock()
			fmt.Printf("Paused. Press p to continue...\n")
			for {
				keyPress := <-d.keyChan
				if keyPress == 'p' {
					break
				}
			}
			fmt.Printf("Continuing...\n")
			time.Sleep(500 * time.Millisecond)
			state.mux.Unlock()
		case 'q':
			state.mux.Lock()
			state.currentWorld.SaveToPgm(d, state.currentTurn)
			exit()
		}
	}
}

// distributor divides the work between workers and interacts with other goroutines.
func distributor(p golParams, d distributorChans, alive chan []cell) {

	state := distributorState{
		currentTurn:  0,
		currentWorld: LoadWorldFromPgm(p.imageHeight, p.imageWidth, d),
		mux:          sync.Mutex{},
	}

	go handleKeyboard(&state, d)

	// Calculate the new state of Game of Life after the given number of turns.
	turnLocal := 0
	for turnLocal < p.turns {
		fmt.Printf("Starting %d\n", turnLocal)
		//Clone the World into local memory
		state.mux.Lock()
		worldLocal := state.currentWorld.Clone()
		state.mux.Unlock()

		//Perform the computation
		worldLocal.computeStateUpdate()
		turnLocal++

		//Update to the new state
		state.mux.Lock()
		state.currentWorld = worldLocal
		state.currentTurn = turnLocal
		state.mux.Unlock()
	}

	// Make sure that the Io has finished any output before exiting.
	d.io.command <- ioCheckIdle
	<-d.io.idle

	// Return the coordinates of cells that are still alive.
	state.mux.Lock()
	alive <- state.currentWorld.CalculateAlive()
}
