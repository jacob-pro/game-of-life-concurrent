use crate::World;

const ALIVE: u8 = 0xFF;
const DEAD: u8 = 0x00;

#[inline(always)]
fn custom_mod(index: i32, max: i32) -> i32 {
    if index >= max {
        return index - max
    } else if index < 0 {
        return index + max
    } else {
        return index
    }
}

#[inline(always)]
fn index(row: i32, column: i32, width: i32) -> usize {
    ((row * width) + column) as usize
}

pub fn calculate_row(world: &World, row: i32) -> Vec<u8> {

    let start = (world.width * row) as usize;
    let end = start + world.width as usize;
    let mut result: Vec<u8> = world.cells[start..end].iter().cloned().collect();

    for col  in 0..world.width  {
        let mut neighbours_alive = 0;

        for i in (row-1)..(row+2) {
            for j in (col-1)..(col+2) {
                let idx = index(custom_mod(i, world.height),custom_mod(j, world.width), world.width);
                if world.cells[idx] == ALIVE {
                    neighbours_alive = neighbours_alive + 1;
                }
            }
        }

        if world.cells[index(row, col, world.width)] == ALIVE {
            neighbours_alive = neighbours_alive - 1;
            if neighbours_alive == 2 || neighbours_alive == 3 {
                result[col as usize] = ALIVE
            } else {
                result[col as usize] = DEAD
            }
        } else if neighbours_alive == 3 {
            result[col as usize] = ALIVE
        }
    }

    return result
}
