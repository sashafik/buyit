# BuyIt E-Commerce Microservices Application

A simple 3-microservice e-commerce application with an API Gateway and Frontend, written in Go.

## Architecture

- **Auth Service** (Port 8081): Handles user registration and login.
- **Product Service** (Port 8082): Manages product inventory and details.
- **Order Service** (Port 8083): Handles order creation and inventory validation.
- **API Gateway** (Port 8080): Routes requests, handles authentication middleware, and manages CORS.
- **Frontend** (Port 3000): Simple HTML/JS UI.

## Project Structure

- `auth-service/`: User management logic (Go).
- `product-service/`: Product catalog logic (Go).
- `order-service/`: Order processing logic (Go).
- `api-gateway/`: Reverse proxy, Auth middleware, CORS (Go).
- `frontend/`: Static HTML/JS files served by Nginx.

## Getting Started

### Prerequisites

- Go 1.21+ (for local run)
- Docker & Docker Compose (for containerized run)

### Running with Docker (Recommended)

To start the entire system using Docker Compose:

```bash
./run_docker_demo.sh
```

Or manually:

```bash
docker-compose up --build
```

Access the application at:
- Frontend: [http://localhost:3000](http://localhost:3000)
- API: [http://localhost:8080](http://localhost:8080)

### Running Locally (Without Docker)

You can start services individually or use the provided script:

```bash
./run_demo.sh
```
(Note: `run_demo.sh` runs services in background and executes tests, then exits. For persistent running, use separate terminals).

To run manually:

1. **Auth Service**: `cd auth-service && go run main.go`
2. **Product Service**: `cd product-service && go run main.go`
3. **Order Service**: `cd order-service && PRODUCT_SERVICE_URL=http://localhost:8082 go run main.go`
4. **API Gateway**: `cd api-gateway && AUTH_SERVICE_URL=http://localhost:8081 PRODUCT_SERVICE_URL=http://localhost:8082 ORDER_SERVICE_URL=http://localhost:8083 go run main.go`
5. **Frontend**: Serve `frontend/` directory (e.g., `cd frontend && python3 -m http.server 3000`).

### Testing

Run the included test script to verify backend services:

```bash
./test_services.sh
```
