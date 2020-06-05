package main

type parallel struct {
	world   world
	workers []workerExternal
}

// Stage 1b parallel implementation
func initParallel(world world, threads int) implementation {
	s := rowsForEachThread(threads, world.height)
	workers := make([]workerExternal, threads)
	for i, rows := range s {
		input := make(chan byte)
		result := make(chan byte)
		workers[i] = workerExternal{
			rows:          rows,
			sendInput:     input,
			receiveResult: result,
		}
		x := workerInternal{
			rows:         rows,
			width:        world.width,
			receiveInput: input,
			sendResult:   result,
		}
		go parallelWorker(x)
	}

	return &parallel{
		world:   world,
		workers: workers,
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
		worldFragment := make([]byte, (w.rows+2)*w.width)
		for i, _ := range worldFragment {
			worldFragment[i] = <-w.receiveInput
		}
		// Compute GoL
		result := gameOfLifeTurn(worldFragment, w.rows, w.width, 1)

		// Send result
		for i, _ := range result {
			w.sendResult <- result[i]
		}
	}
}

func (p *parallel) nextTurn() {

	// Send work
	offset := 0
	for _, worker := range p.workers {
		// Send row above
		sendRowToChannel(p.world.matrix, customMod(offset-1, p.world.height), p.world.width, worker.sendInput)
		// Send rows to compute
		for i := 0; i < worker.rows; i++ {
			sendRowToChannel(p.world.matrix, offset+i, p.world.width, worker.sendInput)
		}
		offset = offset + worker.rows
		// Send row below
		sendRowToChannel(p.world.matrix, customMod(offset, p.world.height), p.world.width, worker.sendInput)
	}

	// Collect results from workers
	recv := 0
	for _, worker := range p.workers {
		size := worker.rows * p.world.width
		for i := 0; i < worker.rows*p.world.width; i++ {
			p.world.matrix[i+recv] = <-worker.receiveResult
		}
		recv += size
	}
}

func (p *parallel) getWorld() world {
	return p.world.clone()
}

func (p *parallel) close() {}
