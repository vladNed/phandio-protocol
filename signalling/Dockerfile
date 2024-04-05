# Use the official Golang image to create a build artifact.
# This is based on Debian and sets the GOPATH to /go.
# https://hub.docker.com/_/golang
FROM golang:1.22.2 as builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
# This is captured separately from the source code to leverage Docker cache
# and avoid re-downloading dependencies if they haven't changed.
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code into the container
COPY cmd/ cmd/
COPY internal/ internal/
COPY resources/certs/cert.pem /etc/ssl/certs/cert.pem
COPY resources/certs/key.pem /etc/ssl/certs/key.pem


# Build the Go app
# -o specifies the output file, in this case the executable name.
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server cmd/server.go

# Start a new stage from scratch
# This is the final stage where the executable is copied over and run
# Using a scratch (empty) image ensures the image size is kept to a minimum
FROM scratch

# Expose port 8080 to the outside world
EXPOSE 8080

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/server .
COPY --from=builder /etc/ssl/certs/cert.pem /etc/ssl/certs/cert.pem
COPY --from=builder /etc/ssl/certs/key.pem /etc/ssl/certs/key.pem

ENV PORT=8080 \
    HOST=0.0.0.0 \
    LOG_LEVEL=info \
    CERT_FILE=/etc/ssl/certs/cert.pem \
    KEY_FILE=/etc/ssl/certs/key.pem

# Command to run the executable
ENTRYPOINT ["./server"]
