.PHONY: build run test clean docker-build docker-run docker-push release help

# Variables
APP_NAME := tcp-echo
DOCKER_IMAGE ?= $(APP_NAME)
DOCKER_TAG ?= latest
DOCKER_REPO ?= 
DOCKER_IMAGE_FULL = $(if $(DOCKER_REPO),$(DOCKER_REPO)/$(DOCKER_IMAGE):$(DOCKER_TAG),$(DOCKER_IMAGE):$(DOCKER_TAG))
PLATFORMS ?= linux/amd64,linux/arm64
BUILDX_NAME := $(APP_NAME)-builder

# Default target
all: build

# Build the application
build:
	@echo "Building $(APP_NAME)..."
	go build -o $(APP_NAME)

# Run the application
run: build
	@echo "Running $(APP_NAME)..."
	./$(APP_NAME)

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -f $(APP_NAME)

# Ensure buildx builder is set up
setup-buildx:
	@if ! docker buildx inspect $(BUILDX_NAME) >/dev/null 2>&1; then \
		docker buildx create --name $(BUILDX_NAME) --use; \
	fi

# Build Docker image
docker-build: setup-buildx
	@echo "Building Docker image $(DOCKER_IMAGE_FULL)"
	docker buildx build \
		--platform $(PLATFORMS) \
		-t $(DOCKER_IMAGE_FULL) \
		--load \
		.

# Run Docker container
docker-run: docker-build
	@echo "Running Docker container..."
	docker run -it --rm -p 9002:9002 \
		--name $(APP_NAME) \
		$(DOCKER_IMAGE_FULL)

# Push Docker image
docker-push: docker-build
	@if [ -z "$(DOCKER_REPO)" ]; then \
		echo "Error: DOCKER_REPO is not set"; \
		exit 1; \
	fi
	@echo "Pushing $(DOCKER_IMAGE_FULL)..."
	docker push $(DOCKER_IMAGE_FULL)

# Get git commit hash, fallback to 'dev' if not in git repo
GIT_HASH := $(shell git rev-parse --short HEAD 2>/dev/null || echo "dev")

# Release: build and push to zzjcool/tcp-echo
release: setup-buildx
	@echo "Building and pushing release image to zzjcool/tcp-echo"
	@read -p "This will build and push to Docker Hub. Continue? [y/N] " confirm && [ $$confirm = y ] || exit 1
	
	# Build and push with commit hash tag
	docker buildx build \
		--platform $(PLATFORMS) \
		-t zzjcool/tcp-echo:$(GIT_HASH) \
		-t zzjcool/tcp-echo:latest \
		--push \
		.
	
	@echo "✅ Successfully released zzjcool/tcp-echo:$(GIT_HASH) and zzjcool/tcp-echo:latest"

# Show help
help:
	@echo "Available targets:"
	@echo "  build         - Build the application"
	@echo "  run           - Run the application locally"
	@echo "  test          - Run tests"
	@echo "  clean         - Clean build artifacts"
	@echo "  docker-build  - Build Docker image"
	@echo "  docker-run    - Run Docker container"
	@echo "  docker-push   - Push Docker image to registry (requires DOCKER_REPO)"
	@echo "  release      - Build and push release image to zzjcool/tcp-echo"
	@echo "  help          - Show this help message"

# Default target
.DEFAULT_GOAL := help
