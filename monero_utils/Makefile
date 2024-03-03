.PHONY: all clean

# Define variables
GOROOT := $(shell go env GOROOT)
DIST_DIR := ./dist
WASM_EXEC := $(DIST_DIR)/wasm_exec.js
WASM_OUTPUT := $(DIST_DIR)/moneroutils.wasm
MAIN_GO := ./cmd/wasm/main.go

all: $(DIST_DIR) $(WASM_EXEC) $(WASM_OUTPUT)

$(DIST_DIR):
	mkdir -p $(DIST_DIR)

$(WASM_EXEC): $(DIST_DIR)
	cp "$(GOROOT)/misc/wasm/wasm_exec.js" $(DIST_DIR)/

$(WASM_OUTPUT): $(DIST_DIR)
	GOOS=js GOARCH=wasm go build -o $(WASM_OUTPUT) $(MAIN_GO)

clean:
	rm -rf $(DIST_DIR)