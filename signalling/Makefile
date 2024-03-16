.PHONY: run

build:
	go build -o target/server cmd/server.go

run: build
	./target/server

clean:
	rm -rf target

help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  build       Build the signalling server"
	@echo "  run         Run the signalling server"
	@echo "  clean       Remove the target directory"