# MVX-MRN Atomic Swap P2P Node

This is a simple implementation of a P2P node for the MVX-MRN Atomic Swap protocol. It is written in Go and uses pion's WebRTC library for the P2P communication.

## Test Run

To run the P2P node, you need to have installed the following:
- Go
- http-server

Run build wasm first to generate the wasm file:
```bash
make build-wasm
```

Then, run the P2P node in browser:
```bash
make run-test-server
```

Open your browser and go to `http://localhost:3000` to see the P2P node in action.
