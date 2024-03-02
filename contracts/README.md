# XMR-MVX Atomic Swap Contracts

This repository contains the smart contracts for the XMR-MVX atomic swap protocol.

## Development

### Requirements

- Multiversx Python CLI (mxpy) 9.5.1

### Build

To build the contracts you can either run the `build` command from the root of the repository:
```bash
make build
```

Or go into each individual contract directory and run the `build` command:
```bash
mxpy contract build
```

### Clean

To clean the contracts you can run the `clean` command from the root of the repository:
```bash
make clean
```