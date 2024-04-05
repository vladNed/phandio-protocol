# Signalling Server

This is a simple signalling server that can be used to establish a WebRTC connection between two peers.

## Local Requirements

Before building and running the server there some requirements that need to be configured for local development.

### TLS Certificates

The server requires TLS certificates to be able to run. From the `/signalling` directory run the following command to generate the certificates:

```bash
make generate-local-certs
```

You can skip all the prompts by pressing enter. This will generate a `key.pem` and `cert.pem` file in the `resources/certs` directory.

## Running in a container

### Build the docker image

```bash
make build-docker
```

### Run the docker container

```bash
make run-docker
```

Try ping it in the browser accessing `https://localhost:8080/api/v1/ping`

## Running on localhost

> NOTE: For this step you need to have Go installed on your machine and create an `.env` file in the root of the project and copy the contents of the `.env.example` file into it.

### Build the server

```bash
make build
```

### Run the server

```bash
make run
```
