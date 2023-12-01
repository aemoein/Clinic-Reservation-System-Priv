#!/bin/bash

# Stop and remove existing containers
docker-compose down

# Build and run containers
docker-compose up --build
