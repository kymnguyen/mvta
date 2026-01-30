# MVTA - Management Vehicle Tracking Application

A microservices-based vehicle tracking system built with Go backends and React admin UI.

## Architecture Overview

- **vehicle-svc**: Core vehicle management service (Port 50001)
- **tracking-svc**: Vehicle tracking and history service (Port 50002)
- **auth-svc**: Authentication service (Port 50000)
- **admin-web**: React-based admin interface (Port 3000)

## Prerequisites

- Docker & Docker Compose
- Go 1.22+
- Node.js 18+
- MongoDB (via Docker)
- Kafka (via Docker)

## Quick Start

### 1. Start Infrastructure (MongoDB, Kafka, Zookeeper)

```bash
# From project root
docker-compose up -d

# Verify services are running
docker-compose ps
```

### 2. Initialize Go Backend Workspace

```bash
cd apps/backend

# Initialize workspace (one time only)
go work init ./auth-svc ./vehicle-svc ./tracking-svc

# Download dependencies
go mod download
go mod tidy
```

### 3. Start Backend Services

Each service runs on a different port. Start them in separate terminals:

#### Terminal 1: Vehicle Service (Port 50001)
```bash
cd apps/backend/vehicle-svc

# Copy and configure environment
cp example.env .env
# Edit .env if needed (defaults should work locally)

# Run service
go run ./cmd/main.go
```

#### Terminal 2: Tracking Service (Port 50002)
```bash
cd apps/backend/tracking-svc

# Copy and configure environment
cp example.env .env

# Run service
go run ./cmd/main.go
```

#### Terminal 3: Auth Service (Port 50000)
```bash
cd apps/backend/auth-svc

# Copy and configure environment
cp example.env .env

# Run service
go run ./cmd/main.go
```

### 4. Start Frontend Admin UI

```bash
cd apps/admin-web

# Install dependencies
npm install

# Start development server (Port 3000)
npm run dev
```

Visit http://localhost:3000 in your browser.

## Default Environment Variables

Each service has an `example.env` file with required variables:

**vehicle-svc / tracking-svc:**
```env
APP_ENV=local
APP_PORT=50001  # Change for each service
MONGO_URI=mongodb://localhost:27017
MONGO_DB=vehicle_db
KAFKA_BROKERS=localhost:9092
```

**auth-svc:**
```env
APP_ENV=local
APP_PORT=50000
```

## Service Endpoints

### Vehicle Service (Port 50001)
```
POST   /api/v1/vehicles              # Create vehicle
GET    /api/v1/vehicles              # List all vehicles
GET    /api/v1/vehicles/{id}         # Get vehicle details
PATCH  /api/v1/vehicles/{id}/location  # Update location
PATCH  /api/v1/vehicles/{id}/status    # Update status
PATCH  /api/v1/vehicles/{id}/mileage   # Update mileage
PATCH  /api/v1/vehicles/{id}/fuel      # Update fuel level
GET    /health                        # Health check
```

### Tracking Service (Port 50002)
```
GET    /api/v1/vehicles              # List vehicles with tracking data
GET    /api/v1/vehicles/{id}         # Get vehicle with history
GET    /api/v1/vehicles/{id}/history # Get change history
GET    /health                        # Health check
```

### Auth Service (Port 50000)
```
POST   /api/v1/auth/login            # Login
POST   /api/v1/auth/logout           # Logout
GET    /api/v1/auth/me               # Current user info
GET    /health                        # Health check
```

## Admin Web Features

- **Vehicle List (tracking-svc)**: View all tracked vehicles with real-time updates
- **Vehicle CRUD (vehicle-svc)**: Create, read, update vehicles
  - Create new vehicles with VIN, model, location
  - Update vehicle location
  - Update mileage
  - Update fuel level
  - Change vehicle status
- **Change History**: View all modifications to vehicles with timestamps

Navigate between services via the top menu:
- "Tracking (tracking-svc)" - View tracked vehicles
- "Vehicles (vehicle-svc)" - Manage vehicle master data

## Development Workflows

### Run Tests
```bash
cd apps/backend/vehicle-svc

# Run all tests
go test ./... -v

# Run specific test
go test ./... -run TestName -v
```

### Build Docker Images
```bash
# From each service directory
docker build -t mvta/vehicle-svc:latest -f Dockerfile .
```

### View Logs

Each service logs to stdout with structured logging (zap).

Terminal output shows:
- Startup sequence
- Database operations
- Kafka events
- HTTP requests (via middleware)
- Errors and warnings

### Database

MongoDB default database: `vehicle_db`

Collections:
- `vehicles` - Vehicle master data
- `outbox` - Event sourcing outbox (vehicle-svc)
- `vehicle_change_history` - Change tracking (tracking-svc)

## Architecture Patterns

### Event-Driven Communication
- vehicle-svc publishes domain events to Kafka
- tracking-svc consumes events to build history
- Pattern: vehicle.created, vehicle.location_updated, etc.

### Clean Architecture + DDD
Each service follows:
- `/cmd` - Application entry point
- `/internal/api` - HTTP handlers
- `/internal/application` - Commands/Queries
- `/internal/domain` - Business logic
- `/internal/infrastructure` - External integrations

### CQRS Pattern
- Commands: Modify state (CreateVehicle, UpdateLocation)
- Queries: Read state (GetVehicle, GetVehicleHistory)

## Troubleshooting

### Services won't start
```bash
# Check if ports are in use
lsof -i :50001  # vehicle-svc
lsof -i :50002  # tracking-svc
lsof -i :50000  # auth-svc
lsof -i :3000   # admin-web

# Kill process on port if needed
kill -9 <PID>
```

### MongoDB connection issues
```bash
# Check if MongoDB is running
docker-compose ps

# View MongoDB logs
docker-compose logs mongo

# Restart MongoDB
docker-compose restart mongo
```

### Kafka issues
```bash
# Check if Kafka is running
docker-compose ps

# View Kafka logs
docker-compose logs kafka

# Restart Kafka
docker-compose restart kafka
```

### Admin web won't load
```bash
# Check if development server is running
curl http://localhost:3000

# Clear node_modules and reinstall
rm -rf apps/admin-web/node_modules apps/admin-web/package-lock.json
cd apps/admin-web && npm install
```

## Seed Data

Vehicle-svc automatically seeds 10 sample vehicles on first run:
- Fleet Vehicle 001-010
- Various makes/models (Toyota, Honda, Tesla, etc.)
- Different statuses and locations
- Realistic mileage and fuel levels

Seed runs only if no vehicles exist in the database.

## Next Steps

1. ✅ Start all services (see Quick Start)
2. 📱 Open admin UI: http://localhost:3000
3. 🚗 Create, view, and manage vehicles
4. 📊 Monitor vehicle status and history
5. 🔧 Extend with custom features
