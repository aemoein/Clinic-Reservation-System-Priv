#!/bin/bash

# Wait for Kafka to start
sleep 10

# Create Kafka topic
/opt/kafka_2.13-2.8.1/bin/kafka-topics.sh --create --zookeeper zookeeper:2181 --replication-factor 1 --partitions 1 --topic clinic_reservation

# Run the original Kafka entrypoint
/docker-entrypoint.sh "$@"