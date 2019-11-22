package main

func customMod(index int, max int) int {
	if index >= max {
		return index - max
	} else if index < 0 {
		return index + max
	} else {
		return index
	}
}

func updateWorldSerially(w *World) {
	clone := w.Clone()
	for y := 0; y < clone.height; y++ {
		for x := 0; x < clone.width; x++ {
			var neighboursAlive = 0
			//Count alive cells in 3x3 grid
			for i := y - 1; i <= y+1; i++ {
				for j := x - 1; j <= x+1; j++ {
					if clone.matrix[customMod(i, clone.height)][customMod(j, clone.width)] == ALIVE {
						neighboursAlive++
					}
				}
			}
			if clone.matrix[y][x] == ALIVE {
				neighboursAlive--
				if neighboursAlive == 2 || neighboursAlive == 3 {
					w.matrix[y][x] = ALIVE
				} else {
					w.matrix[y][x] = DEAD
				}
			} else if neighboursAlive == 3 {
				w.matrix[y][x] = ALIVE
			}
		}
	}
}
