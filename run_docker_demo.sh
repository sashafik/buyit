#!/bin/bash

echo "Building and starting Docker containers..."
docker-compose up --build -d

echo "Waiting for services to be ready..."
sleep 10

echo "Services are running!"
echo "Frontend: http://localhost:3000"
echo "API Gateway: http://localhost:8080"
echo "Auth Service: http://localhost:8081"
echo "Product Service: http://localhost:8082"
echo "Order Service: http://localhost:8083"

echo ""
echo "To stop services, run: docker-compose down"
