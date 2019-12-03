# concurrent-coursework

To run in Goland/Jetbrains add a `go build` build configuration in the `src` directory. 
Make sure to enable `run.processes.with.pty` in the registry so that termbox will work

- Install Rust https://www.rust-lang.org/tools/install
- cbindgen is required to auto generate C Headers: `cargo install cbindgen`
- Use the makefile in the root dir `make`

To use cgo on Windows requires gcc. 
- `rustup toolchain install stable-gnu`
- Install `msys2`
- In an MSYS2 terminal `pacman --sync mingw-w64-x86_64-gcc`
- Add `C:\msys64\mingw64\bin` to system PATH.

To run make on windows
- In an MSYS2 terminal `pacman --sync make`
- Then use `C:\msys32\usr\bin\make.exe`

Running tests:
- `go test -args -i "halo"`
- `go test -bench . -args -i "halo"`
