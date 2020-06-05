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
			sendRowToChannel(world.matrix, offset+i, world.width, worldChan)
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

	// Loads in world, fragment is 2 rows taller to have space for the rowAbove and rowBelow
	worldFragment := make([]byte, (w.rows+2)*w.width)
	for i := 0; i < (w.rows * w.width); i++ {
		worldFragment[i+w.width] = <-w.receiveWorld
	}

	firstRow := 1
	lastRow := w.rows
	rowAboveStart := 0
	rowAboveEnd := w.width
	rowBelowStart := (w.rows + 1) * w.width
	rowBelowEnd := (w.rows + 2) * w.width

	for {
		proceed := <-w.receiveCommand
		switch proceed {
		case haloTick:
			if w.reverse {
				sendRowToChannel(worldFragment, lastRow, w.width, w.workerBottom)
				for i := rowAboveStart; i < rowAboveEnd; i++ {
					worldFragment[i] = <-w.workerTop
				}

				sendRowToChannel(worldFragment, firstRow, w.width, w.workerTop)
				for i := rowBelowStart; i < rowBelowEnd; i++ {
					worldFragment[i] = <-w.workerBottom
				}
			} else {
				for i := rowAboveStart; i < rowAboveEnd; i++ {
					worldFragment[i] = <-w.workerTop
				}
				sendRowToChannel(worldFragment, lastRow, w.width, w.workerBottom)

				for i := rowBelowStart; i < rowBelowEnd; i++ {
					worldFragment[i] = <-w.workerBottom
				}
				sendRowToChannel(worldFragment, firstRow, w.width, w.workerTop)
			}

			result := gameOfLifeTurn(worldFragment, w.rows, w.width, 1)
			copy(worldFragment[rowAboveEnd:rowBelowStart], result)

			w.sendCompletion <- 0
		case haloGetWorld:
			for i := 0; i < w.rows; i++ {
				sendRowToChannel(worldFragment, i+1, w.width, w.sendResult)
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
	// Collect results from workers
	recv := 0
	for _, worker := range h.workers {
		worker.sendCommand <- haloGetWorld
		size := worker.rows * h.width
		for i := 0; i < worker.rows*h.width; i++ {
			world.matrix[i+recv] = <-worker.getResult
		}
		recv += size
	}
	return world
}

func (h *halo) close() {}
