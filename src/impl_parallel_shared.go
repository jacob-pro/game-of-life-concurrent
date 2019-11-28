package main

import (
	"sync"
)

// Stage 5 option 1, parallel using shared memory
type parallelShared struct {
	world             world
	rowsForEachThread []int
}

func initParallelShared(world world, threads int) implementation {
	return &parallelShared{
		world:             world,
		rowsForEachThread: rowsForEachThread(threads, world.height),
	}
}

type parallelCell struct {
	offset int
	rows   int
	result [][]byte
}

// GoL for one cell
func (p *parallelCell) compute(wg *sync.WaitGroup, world *world) {
	defer wg.Done()
	p.result = gameOfLifeTurn(func(i int) []byte {
		return world.matrix[customMod(i+p.offset, world.height)]
	}, p.rows, world.width)
}

func (p *parallelShared) nextTurn() {

	// Split the world into cells
	cells := make([]*parallelCell, len(p.rowsForEachThread))
	i := 0
	for thread, rows := range p.rowsForEachThread {
		cells[thread] = &parallelCell{
			offset: i,
			rows:   rows,
			result: nil,
		}
		i = i + rows
	}

	// Do the computations on worker threads
	wg := sync.WaitGroup{}
	for _, cell := range cells {
		wg.Add(1)
		go cell.compute(&wg, &p.world)
	}
	wg.Wait()

	// Reconstruct the world at end of turn
	var result [][]byte
	for _, cell := range cells {
		result = append(result, cell.result...)
	}
	p.world.matrix = result
}

func (p *parallelShared) getWorld() world {
	return p.world.clone()
}

func (p *parallelShared) close() {}
