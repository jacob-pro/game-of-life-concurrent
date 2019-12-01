package main

type halo struct {
	startWorld world
	workers []dist
}

type dist struct {
	rows         int
	sendTick     chan<- bool // False causes return, True causes continue
	getResult    <-chan byte
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
			rows:         rows,
			sendTick:     tickChan,
			getResult:    resultChan,
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
		proceed := <- w.receiveTick
		if proceed { // Process Turn
			if w.reverse {
				sendRowToChannel(worldFragment[w.rows-1], w.width, w.workerBottom)
				for i := 0; i < w.width; i++ {
					rowAbove[i] = <-w.workerTop
				}

				sendRowToChannel(worldFragment[0], w.width, w.workerTop)
				for i := 0; i < w.width; i++ {
					rowBelow[i] = <-w.workerBottom
				}
			} else {
				for i := 0; i < w.width; i++ {
					rowAbove[i] = <-w.workerTop
				}
				sendRowToChannel(worldFragment[w.rows-1], w.width, w.workerBottom)

				for i := 0; i < w.width; i++ {
					rowBelow[i] = <-w.workerBottom
				}
				sendRowToChannel(worldFragment[0], w.width, w.workerTop)
			}
		} else { // Return Data
			for i := range worldFragment {
				sendRowToChannel(worldFragment[i], w.width, w.sendResult)
			}
		}

		/* Done to here */


		// Compute GoL

	}
}

func (h *halo) nextTurn() {
	for i := range h.workers {
		h.workers[i].sendTick <- true
	}
}

func (h *halo) getWorld() world { // Gets rows from each worker in turn
	offset := 0
	for i := range h.workers {
		for j := offset; j < offset + h.workers[i].rows; j++ {
			for k := 0; k < h.startWorld.width; k++ {
				h.startWorld.matrix[offset + j][k] = <-h.workers[i].getResult
			}
		}
	}
	return h.startWorld
}

func (h *halo) close() {}