.PHONY: all clean build-cli

# Define variables
GOROOT := $(shell go env GOROOT)
DIST_DIR := ./dist
WASM_EXEC := $(DIST_DIR)/wasm_exec.js
WASM_OUTPUT := $(DIST_DIR)/moneroutils.wasm
MAIN_GO := ./cmd/wasm/main.go

all: $(DIST_DIR) $(WASM_EXEC) $(WASM_OUTPUT) build-cli

$(DIST_DIR):
	mkdir -p $(DIST_DIR)

$(WASM_EXEC): $(DIST_DIR)
	cp "$(GOROOT)/misc/wasm/wasm_exec.js" $(DIST_DIR)/

$(WASM_OUTPUT): $(DIST_DIR)
	GOOS=js GOARCH=wasm go build -o $(WASM_OUTPUT) $(MAIN_GO)

clean:
	rm -rf $(DIST_DIR) ./bin

build-cli:
	go build -o ./bin/cli ./cmd/cli/