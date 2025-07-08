# Build stage
FROM --platform=$BUILDPLATFORM golang:1.21-alpine AS builder

WORKDIR /app

# Set Go environment variables
ENV GOPROXY=https://goproxy.cn,direct \
    CGO_ENABLED=0 \
    GO111MODULE=on

# Copy only go.mod first to leverage Docker cache
COPY go.mod .

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN go build -ldflags='-w -s' -o tcp-echo .


FROM busybox

# Copy the binary from builder
COPY --from=builder /app/tcp-echo /tcp-echo

# Copy SSL certificates (needed for HTTPS if your app makes external requests)
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Set environment variables with defaults
ENV ECHO_SERVER_HOST=0.0.0.0 \
    ECHO_SERVER_PORT=9002 \
    ECHO_SERVER_PREFIX="ECHO: "

# Expose the default port
EXPOSE 9002

# Run the application
ENTRYPOINT ["/tcp-echo"]
