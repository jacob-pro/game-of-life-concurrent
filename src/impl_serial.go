package main

type serial struct {
	world world
}

// Stage 1a single threaded implementation
func initSerial(world world, _ int) implementation {
	return &serial{world}
}

func (s *serial) nextTurn() {
	s.world.matrix = gameOfLifeTurn(s.world.matrix, s.world.height, s.world.width, 0)
}

func (s *serial) getWorld() world {
	return s.world.clone()
}

func (p *serial) close() {}
