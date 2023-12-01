#!/bin/bash

DOCKERFILES=("Dockerfile.zookeeper" "Dockerfile.mysql" "Dockerfile.react" "Dockerfile.kafka" "Dockerfile.golang")

# Build and run containers
for DOCKERFILE in "${DOCKERFILES[@]}"; do
    IMAGE_NAME=$(basename "${DOCKERFILE}" .dockerfile | tr '[:upper:]' '[:lower:]')
    
    # Build Docker image
    echo "Building ${IMAGE_NAME} image..."
    docker build -t "${IMAGE_NAME}" -f "${DOCKERFILE}" .

    # Run Docker container
    echo "Running ${IMAGE_NAME} container..."
    docker run -d --name "${IMAGE_NAME}_container" "${IMAGE_NAME}"

    echo "========================================"
done

echo "Containers built and running successfully!"
