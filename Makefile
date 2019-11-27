make:
	cd rust && $(MAKE)
	go env -w CGO_LDFLAGS_ALLOW=".*" CGO_CFLAGS_ALLOW=".*"
	cd src && $(MAKE)