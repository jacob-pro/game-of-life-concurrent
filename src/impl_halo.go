package main

type halo struct {
	height  int
	width   int
	workers []dist
}

type haloCommand uint8

const (
	haloTick haloCommand = iota
	haloGetWorld
)

type dist struct {
	rows              int
	sendCommand       chan<- haloCommand
	receiveCompletion <-chan int
	getResult         <-chan byte
}

type worker struct {
	rows           int
	width          int
	receiveCommand <-chan haloCommand
	sendCompletion chan<- int
	receiveWorld   <-chan byte
	sendResult     chan<- byte
	reverse        bool
	workerTop      chan byte
	workerBottom   chan byte
}

// Stage 4 halo exchange implementation
func initHalo(world world, threads int) implementation {
	s := rowsForEachThread(threads, world.height)
	workers := make([]dist, threads)

	workerChans := make([]chan byte, threads+1)
	for i := range workerChans {
		workerChans[i] = make(chan byte)
	}

	offset := 0

	for i, rows := range s {
		commandChan := make(chan haloCommand)
		completeChan := make(chan int)
		worldChan := make(chan byte)
		resultChan := make(chan byte)
		workers[i] = dist{
			rows:              rows,
			sendCommand:       commandChan,
			receiveCompletion: completeChan,
			getResult:         resultChan,
		}
		reverse := false
		if i == 0 {
			reverse = true
		}
		x := worker{
			rows:           rows,
			width:          world.width,
			receiveCommand: commandChan,
			sendCompletion: completeChan,
			sendResult:     resultChan,
			reverse:        reverse,
			receiveWorld:   worldChan,
			workerTop:      workerChans[customMod(i-1, threads)], // Assign comm chan to previous thread. May not work.
			workerBottom:   workerChans[customMod(i, threads)],
		}
		go haloWorker(x)
		for i := 0; i < x.rows; i++ {
			sendRowToChannel(world.matrix[offset+i], world.width, worldChan)
		}
		offset += x.rows
	}
	return &halo{
		height:  world.height,
		width:   world.width,
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
		proceed := <-w.receiveCommand
		switch proceed {
		case haloTick:
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

			worldFragment = gameOfLifeTurn(func(i int) []byte {
				if i == -1 {
					return rowAbove
				} else if i == w.rows {
					return rowBelow
				} else {
					return worldFragment[i]
				}
			}, w.rows, w.width)

			w.sendCompletion <- 0
		case haloGetWorld:
			for i := range worldFragment {
				sendRowToChannel(worldFragment[i], w.width, w.sendResult)
			}
		}
	}
}

func (h *halo) nextTurn() {
	for _, w := range h.workers {
		w.sendCommand <- haloTick
	}
	for _, w := range h.workers {
		<-w.receiveCompletion
	}
}

func (h *halo) getWorld() world {
	world := newWorld(h.height, h.width)
	rowCounter := 0
	for _, worker := range h.workers {
		worker.sendCommand <- haloGetWorld
		for i := 0; i < worker.rows; i++ {
			for w := 0; w < h.width; w++ {
				world.matrix[rowCounter][w] = <-worker.getResult
			}
			rowCounter++
		}
	}
	return world
}

func (h *halo) close() {}
