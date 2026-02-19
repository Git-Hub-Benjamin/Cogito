GO ?= $(shell which go 2>/dev/null || echo $(HOME)/.local/go/bin/go)

build:
	$(GO) build -o cogito .

install: build
	cp cogito ~/.local/bin/cogito
	cp launch-cogito.sh ~/.local/bin/launch-cogito.sh

uninstall:
	rm -f ~/.local/bin/cogito ~/.local/bin/launch-cogito.sh
