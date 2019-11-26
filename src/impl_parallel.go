package main

// Stage 1b parallel implementation
type Parallel struct {
	world   World
	workers []workerExternal
}

func (p *Parallel) Init(world World, threads int) {
	p.world = world

	s := rowsForEachThread(threads, p.world.height)
	p.workers = make([]workerExternal, threads)
	for i, rows := range s {
		input := make(chan byte)
		result := make(chan byte)
		p.workers[i] = workerExternal{
			rows:          rows,
			sendInput:     input,
			receiveResult: result,
		}
		x := workerInternal{
			rows:         rows,
			width:        p.world.width,
			receiveInput: input,
			sendResult:   result,
		}
		go parallelWorker(x)
	}
}

type workerExternal struct {
	rows          int
	sendInput     chan<- byte
	receiveResult <-chan byte
}

type workerInternal struct {
	rows         int
	width        int
	receiveInput <-chan byte
	sendResult   chan<- byte
}

func parallelWorker(w workerInternal) {
	for {
		// Receive the world fragment
		worldFragment := make([][]byte, w.rows+2)
		for i, _ := range worldFragment {
			worldFragment[i] = make([]byte, w.width)
			for x := 0; x < w.width; x++ {
				worldFragment[i][x] = <-w.receiveInput
			}
		}
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

func sendRowToChannel(row []byte, width int, dest chan<- byte) {
	for i := 0; i < width; i++ {
		dest <- row[i]
	}
}

func (p *Parallel) NextTurn() {

	// Send work
	offset := 0
	for _, worker := range p.workers {
		// Send row above
		sendRowToChannel(p.world.matrix[customMod(offset-1, p.world.height)], p.world.width, worker.sendInput)
		// Send rows to compute
		for i := 0; i < worker.rows; i++ {
			sendRowToChannel(p.world.matrix[offset+i], p.world.width, worker.sendInput)
		}
		offset = offset + worker.rows
		// Send row below
		sendRowToChannel(p.world.matrix[customMod(offset, p.world.height)], p.world.width, worker.sendInput)
	}

	// Collect results from workers
	rowCounter := 0
	for _, worker := range p.workers {
		for i := 0; i < worker.rows; i++ {
			for w := 0; w < p.world.width; w++ {
				p.world.matrix[rowCounter][w] = <-worker.receiveResult
			}
			rowCounter++
		}
	}
}

func (p *Parallel) GetWorld() World {
	return p.world.Clone()
}
