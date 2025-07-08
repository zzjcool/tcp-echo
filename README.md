# TCP Echo Server

A simple TCP echo server written in Go that echoes back any message it receives, with an optional prefix.

## Features

- Configurable host and port
- Customizable echo prefix
- Displays server information on connection
- Containerized with Docker
- Simple Makefile for common tasks

## Prerequisites

- Go 1.21 or higher
- Docker (optional, for containerized deployment)
- Make (optional, for using the Makefile)

## Getting Started

### Building and Running Locally

1. Clone the repository:
   ```bash
   git clone <repository-url>
   cd tcp-echo
   ```

2. Build the application:
   ```bash
   make build
   ```

3. Run the server:
   ```bash
   ./tcp-echo
   ```

   Or use the Makefile:
   ```bash
   make run
   ```

### Using Docker

1. Build the Docker image:
   ```bash
   make docker-build
   ```

2. Run the container:
   ```bash
   make docker-run
   ```

### Configuration

The server can be configured using the following environment variables:

| Variable | Description | Default |
|----------|-------------|---------|
| `ECHO_SERVER_HOST` | Host to bind the server to | `0.0.0.0` |
| `ECHO_SERVER_PORT` | Port to listen on | `9002` |
| `ECHO_SERVER_PREFIX` | Prefix to add to echoed messages | `ECHO: ` |

Example with custom configuration:
```bash
ECHO_SERVER_HOST=127.0.0.1 ECHO_SERVER_PORT=8080 ECHO_SERVER_PREFIX="SERVER: " ./tcp-echo
```

## Testing

To run the tests:

```bash
make test
```

## Connecting to the Server

Use any TCP client to connect to the server. For example, using `netcat`:

```bash
nc localhost 9002
```

Once connected, type any message and press Enter to see it echoed back with the configured prefix.

## Building for Production

### Docker

To build and push the Docker image to a registry:

```bash
export DOCKER_REPO=your-docker-repo
make docker-push
```

### Binary Release

To build for a specific platform:

```bash
GOOS=linux GOARCH=amd64 go build -o tcp-echo-linux-amd64
```

## License

[MIT](LICENSE)

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
