.PHONY: all, help

build-wasm:
	@echo "Building wasm"
	@GOOS=js GOARCH=wasm go build -o test-app/p2p.wasm ./cmd/wasm/main.go

run-test-server:
	@echo "Running test server"
	@cd ./test-app && http-server -p 3000

help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  build-wasm  - Build wasm"
	@echo "  run-test-server  - Run test server"
