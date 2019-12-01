package main

type halo struct {
	startWorld world
	workers []dist
}

type dist struct {
	rows          int
	sendTick     chan<- bool // False causes return, True causes continue
}

type worker struct {
	rows         int
	width        int
	receiveTick  <-chan bool
	receiveWorld <-chan byte
	sendResult   chan<- byte
	reverse      bool
	workerTop    chan byte
	workerBottom chan byte
}

func initHalo(world world, threads int) implementation {
	s := rowsForEachThread(threads, world.height)
	workers := make([]dist, threads)

	workerChans := make([]chan byte, threads+1)
	for i := range workerChans {
		workerChans[i] = make(chan byte)
	}

	offset := 0

	for i, rows := range s {
		tickChan := make(chan bool)
		worldChan := make(chan byte)
		resultChan := make(chan byte)
		workers[i] = dist{
			rows:          rows,
			sendTick:     tickChan,
		}
		reverse := false
		if i == 0 {
			reverse = true
		}
		x := worker{
			rows:         rows,
			width:        world.width,
			receiveTick:  tickChan,
			sendResult:   resultChan,
			reverse:      reverse,
			receiveWorld: worldChan,
			workerTop:    workerChans[customMod(i - 1, threads)], // Assign comm chan to previous thread. May not work.
			workerBottom: workerChans[customMod(i, threads)],
		}
		go haloWorker(x)
		for i := 0; i < x.rows; i++ {
			sendRowToChannel(world.matrix[offset+i], world.width, worldChan)
		}
		offset += x.rows
	}
	return &halo {
		startWorld: world,
		workers: workers,
	}
}

func haloWorker(w worker) {

	// Loads in world
	worldFragment := make([][]byte, w.rows)
	for i, _ := range worldFragment {
		worldFragment[i] = make([]byte, w.width)
		for x := 0; x < w.width; x++ {
			worldFragment[i][x] = <-w.receiveWorld
		}
	}

	rowAbove := make([]byte, w.width)
	rowBelow := make([]byte, w.width)


	for {
		if w.reverse {
			sendRowToChannel(worldFragment[w.rows - 1], w.width, w.workerBottom)
			for i := 0; i < w.width; i++ {
				rowAbove[i] = <- w.workerTop
			}

			sendRowToChannel(worldFragment[0], w.width, w.workerTop)
			for i := 0; i < w.width; i++ {
				rowBelow[i] = <- w.workerBottom
			}
		} else {
			for i := 0; i < w.width; i++ {
				rowAbove[i] = <- w.workerTop
			}
			sendRowToChannel(worldFragment[w.rows - 1], w.width, w.workerBottom)

			for i := 0; i < w.width; i++ {
				rowBelow[i] = <- w.workerBottom
			}
			sendRowToChannel(worldFragment[0], w.width, w.workerTop)
		}

		/* Done to here */


		// Compute GoL
		result := gameOfLifeTurn(func(i int) []byte {
			return worldFragment[i+1]
		}, w.rows, w.width)

		// Send result
		for y := 0; y < w.rows; y++ {
			for x := 0; x < w.width; x++ {
				w.sendResult <- result[y][x]
			}
		}
	}
}

func (h *halo) nextTurn() {
}

func (h *halo) getWorld() world {
	return h.startWorld
}

func (h *halo) close() {}