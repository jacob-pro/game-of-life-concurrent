package main

// #cgo CFLAGS: -I../rust/include
// #cgo windows LDFLAGS: ../rust/target/release/gol_rust.lib -lws2_32 -luserenv
// #cgo linux darwin LDFLAGS: ../rust/target/release/libgol_rust.a -ldl
// #include <gol.h>
import "C"

// Stage 5 some crazy shit
type rust struct {
	gol    *C.GameOfLife
	height int
	width  int
}

func initRust(world world, threads int) implementation {
	// Flatten the world
	var k []byte
	for _, v := range world.matrix {
		k = append(k, v...)
	}
	//noinspection ALL
	gol := C.gol_init((*C.uchar)(&k[0]), C.int32_t(world.height), C.int32_t(world.width), C.int32_t(threads))

	return &rust{
		gol:    gol,
		height: world.height,
		width:  world.width,
	}
}

func (r *rust) nextTurn() {
	C.gol_next_turn(r.gol)
}

func (r *rust) getWorld() world {
	// Load the world into a slice
	b := make([]byte, r.width*r.height)
	C.gol_get_world(r.gol, (*C.uchar)(&b[0]))

	// Unflatten
	unflat := make([][]byte, r.height)
	for i := 0; i < r.height; i++ {
		start := i * r.width
		end := start + r.width
		unflat[i] = b[start:end]
	}

	return world{
		width:  r.width,
		height: r.height,
		matrix: unflat,
	}
}

func (r *rust) close() {
	C.gol_destroy(r.gol)
}
