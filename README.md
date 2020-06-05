# concurrent-coursework

To run in Goland/Jetbrains add a `go build` build configuration in the `src` directory. 
Make sure to enable `run.processes.with.pty` in the registry so that termbox will work

- Install Rust https://www.rust-lang.org/tools/install

To use cgo on Windows requires the Rust library be compiled with the GCC/GNU toolchain. 
- `rustup toolchain install stable-gnu`
- In an MINGW64 terminal `pacman -S mingw-w64-x86_64-gcc`
- Install GNU make `pacman -S make`

Build: `cd src && make`

Running tests:
- `go test -args -i "halo"`
- `go test -bench . -args -i "halo"`
