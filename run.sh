#!/bin/bash

# Create network if it doesn't exist
docker network create mynetwork 2>/dev/null || true

# Build and run RabbitMQ container
docker build -t rabbitmqimage -f Dockerfile.rabbitmq .
docker run -d --net mynetwork --name rabbitmqcontainer -p 5672:5672 -p 15672:15672 rabbitmqimage

# Build and run MySQL container
docker build -t dbimage -f Dockerfile.mysql .
docker run -d --net mynetwork --name dbcontainer -p 3307:3306 dbimage

echo "Waiting for MySQL container to start..."
sleep 10

# Build and run React container
docker build -t frontimage -f Dockerfile.react .
docker run -d --name frontcontainer --net mynetwork -p 3000:3000 frontimage

# Build and run Golang container
docker build -t goimage -f Dockerfile.golang .
docker run -d --net mynetwork --name gocontainer -p 8081:8081 goimage