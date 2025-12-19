# 🏦 Clean Architecture Banking Service

A complete banking service built with Go following Clean Architecture principles, Domain-Driven Design (DDD), and SOLID principles.

## 📋 Table of Contents

- [Features](#-features)
- [Architecture](#-architecture)
- [Prerequisites](#-prerequisites)
- [Quick Start](#-quick-start)
- [Database Setup](#-database-setup)
- [Project Structure](#-project-structure)
- [API Documentation](#-api-documentation)
- [Testing](#-testing)
- [Configuration](#-configuration)
- [Development](#-development)

## ✨ Features

- **💰 Balance Inquiry** - Check wallet balance by user ID
- **🏧 Money Withdrawal** - Withdraw funds with validation
- **📊 Transaction Recording** - All operations are logged
- **🔍 Input Validation** - Comprehensive request validation
- **🏥 Health Checks** - Service health monitoring
- **📈 Versioned API** - Both legacy and versioned endpoints
- **🧪 Comprehensive Testing** - Unit tests and integration tests
- **🗄️ Database Support** - PostgreSQL with proper schema
- **🏗️ Clean Architecture** - Proper separation of concerns

## 🏛️ Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                    Infrastructure Layer                      │
│  ┌─────────────────┐  ┌─────────────────┐  ┌──────────────┐ │
│  │   HTTP Handlers │  │  HTTP Routing   │  │ Persistence │ │
│  │                 │  │                 │  │   (Postgres) │ │
│  └─────────────────┘  └─────────────────┘  └──────────────┘ │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                    Application Layer                         │
│  ┌─────────────────┐  ┌─────────────────┐                  │
│  │    Services     │  │    Use Cases     │                  │
│  │                 │  │                 │                  │
│  └─────────────────┘  └─────────────────┘                  │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                      Domain Layer                             │
│  ┌─────────────────┐  ┌─────────────────┐  ┌──────────────┐ │
│  │    Entities     │  │  Value Objects  │  │ Repositories │ │
│  │                 │  │                 │  │ (Interfaces) │ │
│  └─────────────────┘  └─────────────────┘  └──────────────┘ │
└─────────────────────────────────────────────────────────────┘
```

## 🚀 Prerequisites

- **Go 1.21+** - Go programming language
- **PostgreSQL 13+** - Database server
- **Docker & Docker Compose** - Containerization (optional)
- **Git** - Version control

## ⚡ Quick Start

### 1. Clone the Repository
```bash
git clone <repository-url>
cd bank
```

### 2. Install Dependencies
```bash
go mod download
```

### 3. Setup Database
See [Database Setup](#database-setup) section below

### 4. Build and Run
```bash
# Build the application
go build -o bank-service ./cmd/service

# Run with default settings (in-memory mode)
./bank-service

# Or run with specific configuration
./bank-service --port=8080 --debug=true
```

### 5. Test the API
```bash
# Health check
curl http://localhost:8080/health

# Check balance
curl "http://localhost:8080/balance?user_id=550e8400-e29b-41d4-a716-446655440000"

# Withdraw money
curl -X POST -H "Content-Type: application/json" \
  -d '{"user_id":"550e8400-e29b-41d4-a716-446655440000","amount":20000}' \
  http://localhost:8080/withdraw
```

## 🗄️ Database Setup

### Using PostgreSQL (Recommended)

#### 1. Start PostgreSQL
```bash
# Using Docker (recommended)
docker run --name postgres-bank \
  -e POSTGRES_USER=rio \
  -e POSTGRES_PASSWORD=rio \
  -e POSTGRES_DB=postgres \
  -p 5433:5432 \
  -d postgres:13

# Or start your existing PostgreSQL instance
```

#### 2. Create Database Schema
Execute the SQL file in your DataGrip or any PostgreSQL client:

```sql
-- File: database/schema.sql
-- Or run directly:
psql -h localhost -p 5433 -U rio -d postgres -f database/schema.sql
```

#### 3. Database Configuration
The application will automatically connect to PostgreSQL using:
- **Host**: localhost
- **Port**: 5433
- **User**: rio
- **Password**: rio
- **Database**: postgres

#### 4. Connection String
```
postgresql://rio:rio@localhost:5433/postgres
```

### Using In-Memory (Testing)
By default, the application uses in-memory repositories for quick testing without database setup.

## 📁 Project Structure

```
bank/
├── cmd/service/                    # Application entry point
│   ├── main.go                     # Main application
│   └── config.go                   # Configuration helpers
├── database/
│   └── schema.sql                  # Database schema
├── internal/
│   ├── domain/
│   │   ├── entity/                 # Business entities
│   │   │   ├── wallet.go
│   │   │   └── transaction.go
│   │   ├── valueobject/            # Value objects
│   │   │   ├── money.go
│   │   │   └── userid.go
│   │   ├── repository/             # Repository interfaces
│   │   │   ├── wallet_repository.go
│   │   │   └── transaction_repository.go
│   │   ├── service/                # Service interfaces
│   │   │   └── wallet_service.go
│   │   ├── usecase/                # Use case interfaces
│   │   │   └── withdraw_usecase.go
│   │   └── mocks/                  # Reusable test mocks
│   ├── application/
│   │   ├── service/                # Service implementations
│   │   │   └── wallet_service.go
│   │   ├── usecase/                # Use case implementations
│   │   │   └── withdraw_usecase.go
│   │   └── dto/                    # Data transfer objects
│   │       ├── balance_response.go
│   │       └── withdraw_response.go
│   └── infrastructure/
│       ├── http/                   # HTTP layer
│       │   ├── handlers/
│       │   │   ├── balance_handler.go
│       │   │   └── withdraw_handler.go
│       │   ├── server.go
│       │   └── responses.go
│       └── persistence/            # Data persistence
│           ├── postgres_wallet_repository.go
│           ├── postgres_transaction_repository.go
│           └── memory_*.go          # In-memory implementations
├── README.md                       # This file
└── go.mod                          # Go modules
```

## 📚 API Documentation

### Base URLs
- **Legacy API**: `http://localhost:8080`
- **Versioned API**: `http://localhost:8080/api/v1`

### Endpoints

#### Health Check
```http
GET /health
GET /api/v1/health
```

**Response:**
```json
{
  "status": "ok",
  "message": "Wallet service is running"
}
```

#### Get Balance
```http
GET /balance?user_id={uuid}
GET /api/v1/balance?user_id={uuid}
```

**Query Parameters:**
- `user_id` (required): UUID of the user

**Response:**
```json
{
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "balance": 100000
}
```

#### Withdraw Money
```http
POST /withdraw
POST /api/v1/withdraw
Content-Type: application/json
```

**Request Body:**
```json
{
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "amount": 20000
}
```

**Response (Success):**
```json
{
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "amount_withdrawn": 20000,
  "new_balance": 80000,
  "success": true,
  "message": "withdrawal successful"
}
```

**Response (Error):**
```json
{
  "error": "insufficient_funds",
  "message": "Insufficient funds for withdrawal"
}
```

### Error Responses

All errors return consistent format:
```json
{
  "error": "error_code",
  "message": "Human readable error message"
}
```

**Common Error Codes:**
- `invalid_request` - Invalid JSON format
- `validation_error` - Input validation failed
- `wallet_not_found` - Wallet doesn't exist
- `insufficient_funds` - Not enough balance

## 🧪 Testing

### Run All Tests
```bash
go test ./...
```

### Run Tests with Coverage
```bash
go test -v -cover ./...
```

### Run Specific Tests
```bash
# Domain layer tests
go test ./internal/domain/...

# Integration tests
go test ./internal/infrastructure/http/...
```

### Test Results
```
✅ All domain entities
✅ All value objects
✅ All use cases
✅ All services
✅ All HTTP handlers
✅ All repositories
```

## ⚙️ Configuration

### Command Line Flags
```bash
./bank-service --help

# Available flags:
--host     Server host (default: 0.0.0.0)
--port     Server port (default: 8080)
--debug    Enable debug logging (default: false)
```

### Environment Variables
```bash
export SERVER_HOST=0.0.0.0
export SERVER_PORT=8080
export DEBUG=true
```

### Database Configuration
Currently hardcoded to:
```go
host: localhost
port: 5433
user: rio
password: rio
database: postgres
```

## 🔧 Development

### Build
```bash
# Development build
go build -o bank-service ./cmd/service

# Production build
go build -ldflags="-s -w" -o bank-service ./cmd/service
```

### Run in Development Mode
```bash
# With debug logging
go run ./cmd/service --debug=true

# Or with live reload using air
air
```

### Linting and Formatting
```bash
# Format code
go fmt ./...

# Run linter
golangci-lint run

# Run staticcheck
staticcheck ./...
```

### Database Migrations
For production use, consider adding migration support:
```bash
# Install migration tool
go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Run migrations
migrate -path database/migrations -database "postgresql://rio:rio@localhost:5433/postgres" up
```

## 🐳 Docker Support

### Dockerfile
```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o bank-service ./cmd/service

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/bank-service .
EXPOSE 8080
CMD ["./bank-service"]
```

### Docker Compose
```yaml
version: '3.8'
services:
  postgres:
    image: postgres:13
    environment:
      POSTGRES_USER: rio
      POSTGRES_PASSWORD: rio
      POSTGRES_DB: postgres
    ports:
      - "5433:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data

  bank-service:
    build: .
    ports:
      - "8080:8080"
    depends_on:
      - postgres
    environment:
      - DATABASE_URL=postgresql://rio:rio@postgres:5432/postgres

volumes:
  postgres_data:
```

## 📊 Monitoring and Observability

### Health Check
```bash
curl http://localhost:8080/health
```

### Logging
The application provides structured logging with:
- Request/Response logging
- Error logging with stack traces
- Debug mode for development

### Metrics
Consider adding Prometheus metrics for production:
- Request count and duration
- Transaction counts
- Error rates

## 🔒 Security Considerations

- **Input Validation**: All inputs are validated using struct tags
- **SQL Injection Protection**: Uses parameterized queries
- **Error Handling**: Sensitive information not exposed in error messages
- **CORS**: Configure CORS headers for production
- **Rate Limiting**: Consider adding rate limiting for production

## 🚀 Production Deployment

### Environment Variables
```bash
export SERVER_HOST=0.0.0.0
export SERVER_PORT=8080
export DATABASE_URL=postgresql://rio:rio@localhost:5433/postgres
export LOG_LEVEL=info
export GIN_MODE=release
```

### Building for Production
```bash
# Build optimized binary
go build -ldflags="-s -w" -o bank-service ./cmd/service

# Or use Makefile
make build
```

### Running in Production
```bash
# Create user
useradd -r -s /bin/false bankuser

# Run with systemd
sudo systemctl start bank-service
sudo systemctl enable bank-service
```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🆘 Troubleshooting

### Common Issues

**1. Port already in use**
```bash
# Find what's using the port
lsof -i :8080

# Kill the process
kill -9 <PID>
```

**2. Database connection failed**
```bash
# Check PostgreSQL is running
docker ps | grep postgres

# Test connection
psql -h localhost -p 5433 -U rio -d postgres
```

**3. Build errors**
```bash
# Clean dependencies
go mod tidy
go mod download

# Rebuild
go clean -cache
go build -o bank-service ./cmd/service
```

### Getting Help

- Check the logs: `./bank-service --debug=true`
- Review the test files for usage examples
- Check the GitHub Issues for known problems

---

**Built with ❤️ using Go and Clean Architecture principles**