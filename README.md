# 🏦 Clean Architecture Banking Service

A production-ready banking service built with Go following Clean Architecture principles, Domain-Driven Design (DDD), and Test-Driven Development (TDD). The service provides wallet management with PostgreSQL persistence, transaction safety, and comprehensive testing.

## 📋 Table of Contents

- [Features](#-features)
- [Business Requirements](#-business-requirements)
- [Architecture](#-architecture)
- [Database Schema](#-database-schema)
- [Prerequisites](#-prerequisites)
- [Quick Start](#-quick-start)
- [Project Structure](#-project-structure)
- [API Documentation](#api-documentation)
- [Testing](#-testing)
- [Configuration](#-configuration)
- [Development](#-development)

## ✨ Features

- **💰 Wallet Management** - Each user has one wallet with balance tracking
- **🏧 Safe Withdrawals** - Transactional withdrawals with row-level locking
- **📊 Transaction History** - Complete audit trail of all operations
- **🔍 Input Validation** - Comprehensive UUID and amount validation
- **🏥 Health Checks** - Database connectivity monitoring
- **📈 RESTful API** - Clean JSON API with proper HTTP status codes
- **🧪 Comprehensive Testing** - Unit, integration, and table-driven tests
- **🗄️ PostgreSQL Persistence** - Production-ready database with migrations
- **⚡ Concurrent Safety** - Race condition prevention with row locking
- **🏗️ Clean Architecture** - Proper separation of concerns
- **⚙️ Environment Configuration** - Database connection management
- **🔄 Graceful Shutdown** - Clean server termination

## 🎯 Business Requirements

### Core Business Rules
1. **One Wallet Per User**: Each user has exactly one wallet
2. **Withdrawal Validation**: Cannot withdraw more than available balance
3. **Atomic Operations**: All withdrawals are transactional
4. **Audit Trail**: Every operation is recorded with full details
5. **Integer Currency**: All monetary values use smallest currency unit (no floating point)
6. **Concurrency Safety**: Multiple withdrawals cannot corrupt balance

### Supported Operations
- **Balance Inquiry**: Query current wallet balance
- **Fund Withdrawal**: Withdraw funds with sufficient balance check
- **Transaction Recording**: Automatic audit trail for all operations

### Flow Overview
```text
Withdrawal Flow:
1. Client sends POST /withdraw request
2. API validates UUID and amount
3. Use case begins database transaction
4. Repository locks wallet row (FOR UPDATE)
5. Business logic checks balance sufficiency
6. Repository updates wallet balance
7. Repository records transaction
8. Transaction is committed
9. Success response returned

Balance Inquiry Flow:
1. Client sends GET /balance request
2. API validates UUID parameter
3. Repository queries wallet by user ID
4. Balance response returned
```

## 🏛️ Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                    Infrastructure Layer                      │
│  ┌─────────────────┐  ┌─────────────────┐  ┌──────────────┐ │
│  │   HTTP Handlers │  │  HTTP Routing   │  │ Persistence │ │
│  │                 │  │                 │  │(PostgreSQL)│ │
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

### Layer Responsibilities

**Domain Layer** (Core Business Logic)
- **Entities**: `Wallet`, `Transaction` - Rich domain models with business logic
- **Value Objects**: `Money`, `UserID` - Immutable values with validation
- **Repository Interfaces**: Abstract data access contracts

**Application Layer** (Use Cases)
- **Use Cases**: `WithdrawUseCase`, `BalanceService` - Business process coordination
- **DTOs**: Request/response objects for external communication
- **Services**: Application service orchestration

**Infrastructure Layer** (External Interfaces)
- **HTTP Handlers**: REST API endpoints with validation
- **Persistence**: PostgreSQL repository implementations with SQL queries
- **Database**: Connection management and migrations

## 🗄️ Database Schema

### Tables
```sql
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE users (
                       id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                       name varchar(50),
                       created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                       updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE wallets (
                         id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                         user_id UUID NOT NULL UNIQUE,
                         balance BIGINT NOT NULL DEFAULT 0,
                         created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                         updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),


                         CONSTRAINT wallets_balance_non_negative CHECK (balance >= 0),
                         CONSTRAINT wallets_user_fk FOREIGN KEY (user_id)
                             REFERENCES users(id)
                             ON DELETE CASCADE
);

CREATE TABLE transactions (
                              id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                              wallet_id UUID NOT NULL,
                              transaction_type VARCHAR(20) NOT NULL,
                              amount BIGINT NOT NULL,
                              created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),


                              CONSTRAINT transactions_type_valid CHECK (
                                  transaction_type IN ('WITHDRAWAL', 'DEPOSIT')
                                  ),
                              CONSTRAINT transactions_amount_positive CHECK (amount > 0),
                              CONSTRAINT transactions_wallet_fk FOREIGN KEY (wallet_id)
                                  REFERENCES wallets(id)
                                  ON DELETE CASCADE
);

CREATE INDEX idx_wallets_user_id ON wallets(user_id);
CREATE INDEX idx_transactions_wallet_id ON transactions(wallet_id);

-- 1. Create user
INSERT INTO users (id, name) VALUES ('cfa3b5c8-258a-4d9a-9258-d0ab849ef82d', 'rio');
INSERT INTO users (id, name) VALUES ('cfa3b5c8-258a-4d9a-9258-d0ab849ef82f', 'raihan');


-- 2. Create wallet
INSERT INTO wallets (id, user_id, balance) VALUES (
                                                      '11111111-1111-1111-1111-111111111111',
                                                      'cfa3b5c8-258a-4d9a-9258-d0ab849ef82d',
                                                      100000
                                                  );

INSERT INTO wallets (id, user_id, balance) VALUES (
                                                      '11111111-1111-1111-1111-111111111112',
                                                      'cfa3b5c8-258a-4d9a-9258-d0ab849ef82f',
                                                      500000
                                                  );

```

### Key Features
- **Row-Level Locking**: `SELECT ... FOR UPDATE` prevents concurrent modification
- **Referential Integrity**: Foreign keys ensure data consistency
- **Automatic Timestamps**: Trigger updates `updated_at` automatically
- **Audit Trail**: All operations recorded with complete details

## 🚀 Prerequisites

- **Go 1.21+** - Go programming language
- **PostgreSQL** - Database server (development: any version, production: 12+)
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
```bash
# Start PostgreSQL (Docker example)
docker run --name postgres-bank -e POSTGRES_PASSWORD=postgres -p 5432:5432 -d postgres:latest

# Or use your local PostgreSQL instance
```

### 4. Setup Environment
```bash
# Copy the example environment file
cp .env.example .env

# Edit .env file with your database settings
```

### 5. Run Database Migrations
```bash
# The application will run migrations automatically on first start
# Or run manually:
go run ./cmd/service
```

### 6. Build and Run
```bash
# Build the application
go build -o bank-service ./cmd/service

# Run the application
./bank-service
```

### 7. Test the API
```bash
# Health check
curl http://localhost:8080/health

# Create a wallet (first query will auto-create)
curl "http://localhost:8080/balance?user_id=123e4567-e89b-12d3-a456-426614174000"

# Withdraw funds
curl -X POST http://localhost:8080/withdraw \
  -H "Content-Type: application/json" \
  -d '{"user_id":"123e4567-e89b-12d3-a456-426614174000","amount":5000}'
```

## 📁 Project Structure

```
bank/
├── cmd/service/                    # Application entry point
│   └── main.go                     # Main application
├── .env.example                    # Environment template
├── .env                            # Your environment settings (ignored by Git)
├── go.mod                          # Go modules
├── go.sum                          # Go dependencies lock
├── Makefile                        # Build automation
├── internal/
│   ├── domain/                     # Domain layer (core business logic)
│   │   ├── entity/                 # Business entities
│   │   │   ├── wallet.go
│   │   │   ├── wallet_test.go
│   │   │   ├── transaction.go
│   │   │   └── transaction_test.go
│   │   ├── valueobject/            # Value objects
│   │   │   ├── money.go
│   │   │   ├── money_test.go
│   │   │   ├── userid.go
│   │   │   └── userid_test.go
│   │   ├── repository/             # Repository interfaces
│   │   │   ├── wallet_repository.go
│   │   │   ├── wallet_repository_test.go
│   │   │   ├── wallet_repository_mock.go
│   │   │   ├── transaction_repository.go
│   │   │   └── transaction_repository_mock.go
│   │   ├── service/                # Domain services
│   │   │   └── balance_service.go
│   │   └── usecase/                # Use case interfaces
│   │       ├── withdraw_usecase.go
│   │       └── balance_usecase.go
│   ├── application/                # Application layer (use cases, DTOs)
│   │   ├── service/                # Service implementations
│   │   │   ├── balance_service.go
│   │   │   └── balance_service_test.go
│   │   ├── usecase/                # Use case implementations
│   │   │   ├── withdraw_usecase.go
│   │   │   ├── withdraw_usecase_test.go
│   │   │   ├── balance_usecase.go
│   │   │   ├── balance_usecase_test.go
│   │   │   └── usecase_mock.go
│   │   └── dto/                    # Data transfer objects
│   │       └── wallet_dto.go
│   └── infrastructure/             # Infrastructure layer
│       ├── http/                   # HTTP layer
│       │   ├── balance_handler.go
│       │   ├── balance_handler_test.go
│       │   ├── withdraw_handler.go
│       │   ├── http_mock.go
│       │   ├── responses.go
│       │   └── server.go
│       ├── persistence/            # Database implementations
│       │   ├── wallet_repository.go
│       │   ├── wallet_repository_test.go
│       │   ├── transaction_repository.go
│       │   ├── transaction_repository_test.go
│       │   └── persistence_mock.go
│       └── database/               # Database configuration
│           ├── database.go
│           └── connection_manager.go
├── ARCHITECTURE.md                # Architecture guide
├── DATABASE_IMPLEMENTATION.md    # Database details
└── README.md                       # This file
```

## 📚 API Documentation

### Base URL
- **API**: `http://localhost:8080`

### Endpoints

#### Health Check
```http
GET /health
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
```

**Query Parameters:**
- `user_id` (required): UUID of the user

**Response (Success):**
```json
{
  "user_id": "123e4567-e89b-12d3-a456-426614174000",
  "balance": 100000,
  "message": "balance retrieved successfully"
}
```

**Response (Error - Wallet Not Found):**
```json
{
  "error": "wallet_not_found",
  "message": "wallet not found"
}
```

#### Withdraw Money
```http
POST /withdraw
Content-Type: application/json
```

**Request Body:**
```json
{
  "user_id": "123e4567-e89b-12d3-a456-426614174000",
  "amount": 20000
}
```

**Response (Success):**
```json
{
  "user_id": "123e4567-e89b-12d3-a456-426614174000",
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
  "message": "insufficient funds"
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
- `missing_parameter` - Required query parameter missing
- `validation_error` - Input validation failed (UUID format, amount validation)
- `wallet_not_found` - Wallet doesn't exist
- `insufficient_funds` - Not enough balance for withdrawal
- `failed_to_begin_transaction` - Database transaction error

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

# Infrastructure layer tests
go test ./internal/infrastructure/...

# With coverage report
go test -v -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Test Categories
- **Unit Tests**: Domain entities and value objects
- **Integration Tests**: Repository implementations with SQL mocks
- **Use Case Tests**: Business logic orchestration
- **Service Tests**: Application service coordination
- **HTTP Tests**: API endpoints and validation
- **Table-Driven Tests**: Comprehensive scenario testing

### Test Structure
All tests follow table-driven patterns with clear test cases:
- Success scenarios
- Error conditions
- Edge cases (zero values, maximum values)
- Concurrency scenarios
- Database transaction flows

## ⚙️ Configuration

### Environment Variables

Configure using `.env` file:
```bash
# Database Configuration
DB_HOST=localhost              # Database host
DB_PORT=5432                  # Database port
DB_USER=postgres              # Database user
DB_PASSWORD=postgres          # Database password
DB_NAME=wallet_db             # Database name
DB_SSLMODE=disable             # SSL mode (development: disable)

# Server Configuration
SERVER_HOST=localhost          # Server host (default: 0.0.0.0)
SERVER_PORT=8080              # Server port (default: 8080)

# Logging Configuration
DEBUG=false                   # Enable debug logging (default: false)
```

### Database Setup

**Development:**
```sql
-- Create database
CREATE DATABASE wallet_db;

-- Connect to wallet_db and run migrations automatically
```

**Production:**
```sql
-- Create database and user
CREATE DATABASE wallet_db;
CREATE USER wallet_user WITH PASSWORD 'secure_password';
GRANT ALL PRIVILEGES ON DATABASE wallet_db TO wallet_user;
```

## 🔧 Development

### Build
```bash
# Development build
go build -o bank-service ./cmd/service

# Production build
go build -ldflags="-s -w" -o bank-service ./cmd/service
```

### Using Makefile
```bash
# Build application
make build

# Run tests
make test

# Run with coverage
make test-coverage

# Clean build artifacts
make clean

# Development run
make run

# Production build
make build-prod
```

### Linting and Formatting
```bash
# Format code
go fmt ./...

# Run linter
golangci-lint run

# Run staticcheck
staticcheck ./...

# Vet code
go vet ./...
```

## 🔄 Graceful Shutdown

The application supports graceful shutdown for production use:

- **Signal Handling**: Responds to `Ctrl+C` and `kill` commands
- **Transaction Completion**: Completes in-progress database transactions
- **Resource Cleanup**: Properly closes database connections
- **Zero Data Loss**: Ensures all operations complete safely

## 🐳 Docker Support

### Dockerfile
```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o bank-service ./cmd/service

FROM postgres:15-alpine
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /app/bank-service .
EXPOSE 8080
CMD ["./bank-service"]
```

### Docker Compose
```yaml
version: '3.8'
services:
  postgres:
    image: postgres:15-alpine
    environment:
      POSTGRES_DB: wallet_db
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data

  bank-service:
    build: .
    ports:
      - "8080:8080"
    environment:
      - DB_HOST=postgres
      - DB_PORT=5432
      - DB_USER=postgres
      - DB_PASSWORD=postgres
      - DB_NAME=wallet_db
    depends_on:
      - postgres

volumes:
  postgres_data:
```

## 📊 Monitoring and Observability

### Health Check
```bash
curl http://localhost:8080/health
```

### Logging
- Structured logging with request/response correlation
- Database query logging in debug mode
- Error logging with stack traces
- Transaction flow tracing

### Database Monitoring
- Connection pool health
- Query performance metrics
- Transaction lock monitoring

## 🔒 Security Considerations

- **Input Validation**: All inputs validated with comprehensive rules
- **SQL Injection Prevention**: Parameterized queries throughout
- **Transaction Safety**: Row-level locking prevents race conditions
- **Error Information**: No sensitive data exposed in error messages
- **Environment Security**: Database credentials in environment files

## 🚀 Production Deployment

### Environment Setup
```bash
# Production build
make build-prod

# Configure production environment
cat > .env << EOF
DB_HOST=your-db-host
DB_PORT=5432
DB_USER=your-db-user
DB_PASSWORD=your-secure-password
DB_NAME=wallet_db
DB_SSLMODE=require
SERVER_HOST=0.0.0.0
SERVER_PORT=8080
DEBUG=false
EOF

# Run with process manager (systemd, supervisor, etc.)
./bank-service
```

### Database Setup
```bash
# Create database and run migrations
psql -h your-db-host -U postgres -c "CREATE DATABASE wallet_db;"
psql -h your-db-host -U postgres -d wallet_db -f schema.sql
```

### Performance Considerations
- **Connection Pooling**: PostgreSQL connection pool configuration
- **Index Optimization**: Optimized indexes for common queries
- **Transaction Timeouts**: Reasonable timeout settings
- **Resource Limits**: Appropriate memory and CPU limits

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Write tests first (TDD approach)
4. Implement your feature
5. Ensure all tests pass (`go test ./...`)
6. Run linting (`golangci-lint run`)
7. Commit your changes (`git commit -m 'Add amazing feature'`)
8. Push to the branch (`git push origin feature/amazing-feature`)
9. Open a Pull Request

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🆘 Troubleshooting

### Common Issues

**Database Connection Issues**
```bash
# Check PostgreSQL status
docker ps | grep postgres

# Test connection
psql -h localhost -U postgres -d wallet_db

# Check environment variables
cat .env
```

**Build Issues**
```bash
# Clean dependencies
go clean -modcache
go mod tidy
go mod download

# Rebuild
go build -o bank-service ./cmd/service
```

**Runtime Issues**
```bash
# Enable debug logging
DEBUG=true ./bank-service

# Check database connectivity
curl http://localhost:8080/health

# Check database tables
psql -h localhost -U postgres -d wallet_db -c "\dt"
```

### Getting Help

- Check logs with debug mode enabled
- Review test files for usage examples
- Check [ARCHITECTURE.md](ARCHITECTURE.md) for design details
- Check [DATABASE_IMPLEMENTATION.md](DATABASE_IMPLEMENTATION.md) for database specifics

---

**Built with ❤️ using Go, Clean Architecture, and Domain-Driven Design**