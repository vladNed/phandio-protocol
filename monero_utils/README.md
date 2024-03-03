## Makefile Commands

The following commands are available in the Makefile:

- `make all`: Builds both the WebAssembly code and the CLI executable.
- `make build-cli`: Build only the CLI executable.
- `make clean`: Clean up generated files and directories.

```bash
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