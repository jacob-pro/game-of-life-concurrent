package main

func sendRowToChannel(world []byte, row int, width int, dest chan<- byte) {
	for i := 0; i < width; i++ {
		dest <- world[(row*width)+i]
	}
}

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

// Rows is the number of rows to process
// Width is how wide each row is
// The offset is where the starting row is in the matrix
// The result will have size rows * width
func gameOfLifeTurn(matrix []byte, rows int, width int, offset int) []byte {
	mHeight := len(matrix) / width
	result := make([]byte, rows*width)
	for y := 0; y < rows; y++ {
		realY := y + offset
		for x := 0; x < width; x++ {

			result[(y*width)+x] = matrix[(realY*width)+x]

			var neighboursAlive = 0
			//Count alive cells in 3x3 grid
			for i := y - 1; i <= y+1; i++ {
				realI := i + offset
				for j := x - 1; j <= x+1; j++ {
					if matrix[(customMod(realI, mHeight)*width)+customMod(j, width)] == ALIVE {
						neighboursAlive++
					}
				}
			}
			if result[(y*width)+x] == ALIVE {
				neighboursAlive--
				if neighboursAlive == 2 || neighboursAlive == 3 {
					result[(y*width)+x] = ALIVE
				} else {
					result[(y*width)+x] = DEAD
				}
			} else if neighboursAlive == 3 {
				result[(y*width)+x] = ALIVE
			}
		}
	}
	return result
}
