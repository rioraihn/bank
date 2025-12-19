# 🏦 Clean Architecture Banking Service

A complete banking service built with Go following Clean Architecture principles, Domain-Driven Design (DDD), and SOLID principles.

## 📋 Table of Contents

- [Features](#-features)
- [Architecture](#-architecture)
- [Prerequisites](#-prerequisites)
- [Quick Start](#-quick-start)
- [Project Structure](#project-structure)
- [API Documentation](#api-documentation)
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
- **🗄️ In-Memory Storage** - Fast, zero-setup persistence
- **🏗️ Clean Architecture** - Proper separation of concerns
- **⚙️ Environment Configuration** - `.env` file support
- **🔄 Graceful Shutdown** - Clean server termination

## 🏛️ Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                    Infrastructure Layer                      │
│  ┌─────────────────┐  ┌─────────────────┐  ┌──────────────┐ │
│  │   HTTP Handlers │  │  HTTP Routing   │  │ Persistence │ │
│  │                 │  │                 │  │   (Memory)   │ │
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

### 3. Setup Environment
```bash
# Copy the example environment file
cp .env.example .env

# Edit .env file with your preferred settings
# The application will automatically load .env file
```

### 4. Build and Run
```bash
# Build the application
go build -o bank-service ./cmd/service

# Run the application
./bank-service
```

### 5. Test the API
The application starts with empty memory repositories. You need to create wallets first:

```bash
# Health check
curl http://localhost:8080/health

# Note: The following will return "wallet not found" until you create wallets
curl "http://localhost:8080/balance?user_id=your-user-id"
```

**Important**: The application no longer includes pre-populated test data. You'll need to create your own wallets and perform transactions via the API.

## 📁 Project Structure

```
bank/
├── cmd/service/                    # Application entry point
│   ├── main.go                     # Main application
│   └── mocks_test.go               # Test mocks
├── .env.example                    # Environment template
├── .env                            # Your environment settings (ignored by Git)
├── .gitignore                      # Git ignore rules
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
│   │   │   ├── wallet_service.go
│   │   │   └── mocks_test.go       # Service test mocks
│   │   ├── usecase/                # Use case implementations
│   │   │   ├── withdraw_usecase.go
│   │   │   └── mocks_test.go       # Use case test mocks
│   │   └── dto/                    # Data transfer objects
│   │       └── wallet_dto.go
│   └── infrastructure/
│       ├── http/                   # HTTP layer
│       │   ├── handlers/
│       │   │   ├── balance_handler.go
│       │   │   └── withdraw_handler.go
│       │   ├── server.go
│       │   └── responses.go
│       └── persistence/            # Data persistence (in-memory)
│           ├── memory_wallet_repository.go
│           └── memory_transaction_repository.go
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

**Response (Success):**
```json
{
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "balance": 100000
}
```

**Response (Error - Wallet Not Found):**
```json
{
  "error": "wallet_not_found",
  "message": "Wallet not found"
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

# Application layer tests
go test ./internal/application/...

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

### Environment Variables

The application uses a `.env` file for configuration. Copy `.env.example` to `.env` and modify as needed:

**Priority Order:**
1. **Command line flags** (highest priority)
2. **Environment variables** from `.env` file
3. **System environment variables**
4. **Default values** (lowest priority)

#### Configuration Options
```bash
# .env file

# Server Configuration
SERVER_HOST=localhost          # Server host (default: 0.0.0.0)
SERVER_PORT=8080              # Server port (default: 8080)

# Logging Configuration
DEBUG=false                   # Enable debug logging (default: false)
```

### Command Line Flags
```bash
./bank-service --help

# Available flags:
--host     Server host (overrides .env file)
--port     Server port (overrides .env file)
--debug    Enable debug logging (overrides .env file)
```

### Examples

**Development (local only):**
```bash
# .env
SERVER_HOST=localhost
SERVER_PORT=8080
DEBUG=true
```

**Production (all interfaces):**
```bash
# .env
SERVER_HOST=0.0.0.0
SERVER_PORT=8080
DEBUG=false
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

# Or using .env file
echo "DEBUG=true" > .env
go run ./cmd/service
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

## 🔄 Graceful Shutdown

The application supports graceful shutdown for production use:

- **Signal Handling**: Responds to `Ctrl+C` and `kill` commands
- **Request Completion**: Finishes processing existing requests (30-second timeout)
- **Resource Cleanup**: Properly closes connections and releases resources
- **Zero Data Loss**: Ensures database operations complete

**Usage:**
```bash
# Start server
./bank-service

# Graceful shutdown (Ctrl+C)
# OR kill <pid>
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
COPY --from=builder /app/.env.example .env
EXPOSE 8080
CMD ["./bank-service"]
```

### Docker Compose
```yaml
version: '3.8'
services:
  bank-service:
    build: .
    ports:
      - "8080:8080"
    environment:
      - SERVER_HOST=0.0.0.0
      - SERVER_PORT=8080
      - DEBUG=false
    volumes:
      - ./your-.env-file:/app/.env:ro
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
- Debug mode for development (shows file and line numbers)

## 🔒 Security Considerations

- **Input Validation**: All inputs are validated using struct tags
- **Error Handling**: Sensitive information not exposed in error messages
- **Environment Variables**: Secrets stored in `.env` file (ignored by Git)
- **Graceful Shutdown**: Prevents data corruption during termination

## 🚀 Production Deployment

### Environment Setup
```bash
# Copy and configure environment
cp .env.example .env
# Edit .env with production settings

# Build optimized binary
go build -ldflags="-s -w" -o bank-service ./cmd/service

# Run with systemd or process manager
./bank-service
```

### Production Configuration
```bash
# .env for production
SERVER_HOST=0.0.0.0
SERVER_PORT=8080
DEBUG=false
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

# Or change port in .env
echo "SERVER_PORT=8081" >> .env
```

**2. Environment variables not working**
```bash
# Ensure .env file exists
cp .env.example .env

# Check .env file format
cat .env

# Verify application loads .env
./bank-service  # Should show "No .env file found" message if missing
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

**4. Wallet not found errors**
The application starts with empty repositories. You need to create wallets via API before checking balances or making withdrawals.

### Getting Help

- Check the logs: `DEBUG=true ./bank-service`
- Review the test files for usage examples
- Check the GitHub Issues for known problems

---

**Built with ❤️ using Go and Clean Architecture principles**