package main

// Stage 1a single threaded implementation
type Serial struct {
	world World
}

func InitSerial(world World, _ int) Implementation {
	return &Serial{world}
}

func (s *Serial) NextTurn() {
	s.world.matrix = gameOfLifeTurn(func(i int) []byte {
		return s.world.matrix[customMod(i, s.world.height)]
	}, s.world.height, s.world.width)
}

func (s *Serial) GetWorld() World {
	return s.world.Clone()
}

func (p *Serial) Close() {}
