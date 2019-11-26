package main

import (
	"sync"
)

// An alternative version of the stage 1b parallel but using shared memory
type ParallelShared struct {
	world             World
	rowsForEachThread []int
}

func (p *ParallelShared) Init(world World, threads int) {
	p.world = world
	p.rowsForEachThread = rowsForEachThread(threads, p.world.height)
}

type parallelCell struct {
	offset int
	rows   int
	result [][]byte
}

// GoL for one cell
func (p *parallelCell) compute(wg *sync.WaitGroup, world *World) {
	defer wg.Done()

	p.result = gameOfLifeTurn(func(i int) []byte {
		return world.matrix[customMod(i+p.offset, world.height)]
	}, p.rows, world.width)
}

func (p *ParallelShared) NextTurn() {

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

func (p *ParallelShared) GetWorld() World {
	return p.world.Clone()
}
