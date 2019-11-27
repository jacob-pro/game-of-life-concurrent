make:
	cbindgen --config cbindgen.toml --output include/gol.h
ifeq ($(OS),Windows_NT)
	rustup run stable-gnu cargo build --release
else
	cargo build --release
endif
