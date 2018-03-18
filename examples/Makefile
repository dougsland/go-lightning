# Build all go files
#
gofiles = $(shell ls *.go)
all:
	@echo "Building all examples..."
	$(foreach gfile,$(gofiles),go build $(gfile);)

clean:
	@echo "Claning binaries..."
	$(shell find . -type f ! -name '*.go' ! -name 'Makefile' -delete)
