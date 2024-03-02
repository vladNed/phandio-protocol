all: help

.PHONY: clean
clean:
	cd swap && mxpy contract clean
	cd swap-router && mxpy contract clean
	cargo clean

.PHONY: build
build:
	cd swap && mxpy contract build
	cd swap-router && mxpy contract build

help:
	@echo "> clean"
	@echo "> build"
