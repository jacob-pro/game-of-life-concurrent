package main

// Calculate how many rows each thread will be responsible for
func rowsForEachThread(threads int, rows int) []int {
	s := make([]int, threads)
	for row := 0; row < rows; row++ {
		s[row%threads]++
	}
	return s
}

func customMod(index int, max int) int {
	if index >= max {
		return index - max
	} else if index < 0 {
		return index + max
	} else {
		return index
	}
}

// The getRow closure should return values for (-1) to (height + 1)
// Each row should have length equal to width
func gameOfLifeTurn(getRow func(int) []byte, height int, width int) [][]byte {
	result := make([][]byte, height)
	for y := 0; y < height; y++ {
		result[y] = make([]byte, width)
		for x := 0; x < width; x++ {

			//Clone current world into cell result
			result[y][x] = getRow(y)[x]

			var neighboursAlive = 0
			//Count alive cells in 3x3 grid
			for i := y - 1; i <= y+1; i++ {
				for j := x - 1; j <= x+1; j++ {
					if getRow(i)[customMod(j, width)] == ALIVE {
						neighboursAlive++
					}
				}
			}
			if getRow(y)[x] == ALIVE {
				neighboursAlive--
				if neighboursAlive == 2 || neighboursAlive == 3 {
					result[y][x] = ALIVE
				} else {
					result[y][x] = DEAD
				}
			} else if neighboursAlive == 3 {
				result[y][x] = ALIVE
			}
		}
	}
	return result
}
