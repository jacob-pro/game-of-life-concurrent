mod logic;

use std::slice;
use rayon::{ThreadPoolBuilder, ThreadPool};
use rayon::prelude::*;
use logic::calculate_row;

pub struct GameOfLife {
    pool: ThreadPool,
    world: World,
}

pub struct World {
    cells: Vec<u8>,
    height: i32,
    width: i32,
}

/// World is a pointer to a byte array of the flattened world
/// It should therefore have length = height * width
/// The array does not need to live any longer than this function call because it will be cloned
#[no_mangle]
pub extern fn gol_init(world: *mut u8, height: i32, width: i32, threads: i32) -> *mut GameOfLife {

    let slice = unsafe { slice::from_raw_parts(world, (height * width) as usize) };

    let world = World { cells: slice.to_vec(), height, width};

    let pool = ThreadPoolBuilder::new()
        .num_threads(threads as usize)
        .build()
        .unwrap();

    let state = GameOfLife { pool, world };

    Box::into_raw(Box::new(state))
}

#[no_mangle]
pub extern fn gol_next_turn(gol: *mut GameOfLife) {
    let mut gol = unsafe { &mut *gol };

    gol.world.cells = gol.pool.install(|| {

        // Compute each row
        let k: Vec<Vec<u8>> = (0..gol.world.height).into_par_iter().map(|row| {
            calculate_row(&gol.world, row)
        }).collect();

        // Join the rows back together
        return k.into_iter().flat_map(|z| z.into_iter()).collect()
    });

}

/// World is a pointer to an allocated slice of length = width * height
/// Which will be populated with the result
#[no_mangle]
pub extern fn gol_get_world(gol: *mut GameOfLife, world: *mut u8) {
    let gol = unsafe { &*gol };
    let slice = unsafe { slice::from_raw_parts_mut(world, gol.world.cells.len()) };
    slice.copy_from_slice(gol.world.cells.as_slice());
}

#[no_mangle]
pub extern fn gol_destroy(gol: *mut GameOfLife) {
    // When the box goes out of scope it will be deallocated
    unsafe { Box::from_raw(gol); }
}
