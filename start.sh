#!/bin/bash

# Kill any existing processes on ports 8080-8083
lsof -ti:8080,8081,8082,8083 | xargs kill -9 2>/dev/null

echo "Starting Auth Service..."
cd auth-service && go run main.go &
PID1=$!
cd ..

echo "Starting Product Service..."
cd product-service && go run main.go &
PID2=$!
cd ..

echo "Starting Order Service..."
cd order-service && go run main.go &
PID3=$!
cd ..

echo "Starting API Gateway..."
sleep 2 # Wait for services to start
cd api-gateway && go run main.go &
PID4=$!
cd ..

echo "All services started. Press Ctrl+C to stop."
wait $PID1 $PID2 $PID3 $PID4
