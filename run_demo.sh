#!/bin/bash

# Function to cleanup services on exit
cleanup() {
    echo "Stopping services..."
    lsof -ti:8080,8081,8082,8083 | xargs kill -9 2>/dev/null
}
trap cleanup EXIT

# Ensure clean state
cleanup

echo "Starting Auth Service..."
(cd auth-service && go run main.go > ../auth.log 2>&1 &)

echo "Starting Product Service..."
(cd product-service && go run main.go > ../product.log 2>&1 &)

echo "Starting Order Service..."
(cd order-service && go run main.go > ../order.log 2>&1 &)

echo "Starting API Gateway..."
sleep 2
(cd api-gateway && go run main.go > ../gateway.log 2>&1 &)

echo "Waiting for services to initialize..."
sleep 5

echo "Running tests..."
./test_services.sh

echo "Tests completed."
