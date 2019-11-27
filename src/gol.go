package main

import (
	"fmt"
	"sync"
	"time"
)

// The state should be locked before any read/writes to the implementation or current turn
type distributorState struct {
	currentTurn int
	impl        Implementation
	locker      sync.Locker
}

// distributor divides the work between workers and interacts with other goroutines.
// d.keyChan is nil when launched by the test framework
func distributor(p golParams, d distributorChans, alive chan []cell) {

	// The state lock is used to synchronise the Distributor, Keyboard, and Ticker threads.
	// If no user is connected (e.g. testing/benchmark) then we do not need real synchronisation
	// Because only the Distributor thread will be running
	var lock sync.Locker
	if d.keyChan != nil {
		lock = &sync.Mutex{}
	} else {
		lock = NopLocker{}
	}

	i := getInitFn(p.implementationName)(LoadWorldFromPgm(p.imageHeight, p.imageWidth, d), p.threads)

	state := distributorState{
		currentTurn: 0,
		impl:        i,
		locker:      lock,
	}

	if d.keyChan != nil {
		go handleKeyboard(&state, d)
		go ticker(&state)
	}

	turnLocal := 0
	for turnLocal < p.turns {
		state.locker.Lock()
		state.impl.NextTurn()
		turnLocal++
		state.currentTurn = turnLocal
		state.locker.Unlock()
	}

	// Make sure that the Io has finished any output before exiting (there may be a save in progress).
	d.io.command <- ioCheckIdle
	<-d.io.idle

	// Return the coordinates of cells that are still alive.
	state.locker.Lock()
	w := state.impl.GetWorld()
	alive <- w.CalculateAlive()

	state.impl.Close()
}

func handleKeyboard(state *distributorState, d distributorChans) {
	for {
		keyPress := <-d.keyChan
		switch keyPress {
		case 's':
			state.locker.Lock()
			w := state.impl.GetWorld()
			t := state.currentTurn
			state.locker.Unlock()
			w.SaveToPgm(d, t)
		case 'p':
			state.locker.Lock()
			fmt.Printf("Paused. Press p to continue...\n")
			// Wait for p to be pressed again
			// Whilst the mutex remains locked the GoL and Ticker threads will be stuck
			for {
				keyPress := <-d.keyChan
				if keyPress == 'p' {
					break
				}
			}
			fmt.Printf("Continuing...\n")
			state.locker.Unlock()
		case 'q':
			state.locker.Lock()
			w := state.impl.GetWorld()
			w.SaveToPgm(d, state.currentTurn)
			// Make sure that the Io has finished any output before exiting (so that the save completes)
			d.io.command <- ioCheckIdle
			<-d.io.idle
			exit()
		}
	}
}

func ticker(state *distributorState) {
	ticker := time.NewTicker(2 * time.Second)
	for {
		<-ticker.C
		state.locker.Lock()
		world := state.impl.GetWorld()
		turn := state.currentTurn
		state.locker.Unlock()
		fmt.Printf("On turn: %d there are %d alive\n", turn, len(world.CalculateAlive()))
	}
}

func getInitFn(name string) ImplementationInitFn {
	if name == "" {
		return ImplementationDefault.initFn()
	} else {
		impl, err := implementationFromName(name)
		check(err)
		return impl.initFn()
	}
}
