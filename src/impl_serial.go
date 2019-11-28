package main

// Stage 1a single threaded implementation
type serial struct {
	world world
}

func initSerial(world world, _ int) implementation {
	return &serial{world}
}

func (s *serial) nextTurn() {
	s.world.matrix = gameOfLifeTurn(func(i int) []byte {
		return s.world.matrix[customMod(i, s.world.height)]
	}, s.world.height, s.world.width)
}

func (s *serial) getWorld() world {
	return s.world.clone()
}

func (p *serial) close() {}
