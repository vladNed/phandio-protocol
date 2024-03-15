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

## CLI

``bash
$ ./bin/cli
A tool to generate and sum monero keys

Usage:
  cli [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  new         Generates a new wallet
  pub         Pass a private key to derive the associated public key.
  sumsk       Pass two private spend keys to generate their sum.
  view        Pass a private spend key to generate a private view key.

Flags:
  -h, --help   help for cli

Use "cli [command] --help" for more information about a command.
```

```bash
$ ./bin/cli sumsk -h
Pass two private spend keys to generate their sum.

Usage:
  cli sumsk [flags]

Flags:
  -h, --help      help for sumsk
  -v, --verbose   provides all keys when possible
```