.PHONY: run

build:
	go build -o target/server cmd/server.go

build-docker:
	docker build -t signalling-server .

run-local: build
	./target/server

run-docker:
	docker run --name signalling-server -p 8080:8080 signalling-server

clean:
	rm -rf target
	rm -rf .env

generate-local-certs:
	openssl req -x509 -newkey rsa:4096 -keyout resources/certs/key.pem -out resources/certs/cert.pem -days 365 -nodes -addext "subjectAltName = DNS:localhost"

help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  build                Build the signalling server"
	@echo "  build-docker         Build the signalling server docker image"
	@echo "  run-local            Run the signalling server"
	@echo "  run-docker           Run the signalling server in a docker container"
	@echo "  clean                Remove the target directory"
	@echo "  generate-local-certs Generate local certificates"
