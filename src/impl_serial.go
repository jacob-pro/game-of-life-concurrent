package main

// Stage 1a single threaded implementation
type Serial struct {
	world World
}

func (s *Serial) Init(world World, _ int) {
	s.world = world
}

func (s *Serial) NextTurn() {
	s.world.matrix = gameOfLifeTurn(func(i int) []byte {
		return s.world.matrix[customMod(i, s.world.height)]
	}, s.world.height, s.world.width)
}

func (s *Serial) GetWorld() World {
	return s.world.Clone()
}
