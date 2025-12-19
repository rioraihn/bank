# 🚀 Quick Start Guide

Get your Clean Architecture Banking Service running in minutes!

## ⚡ Quick Setup (5 minutes)

### 1. Database Setup
**Using Docker (Easiest):**
```bash
docker run --name postgres-bank \
  -e POSTGRES_USER=rio \
  -e POSTGRES_PASSWORD=rio \
  -e POSTGRES_DB=postgres \
  -p 5433:5432 \
  -d postgres:13
```

**Using DataGrip:**
1. Connect to `localhost:5433`
2. User: `rio`, Password: `rio`
3. Execute the SQL from `database/schema.sql`

### 2. Build & Run
```bash
# Build the service
go build -o bank-service ./cmd/service

# Run it!
./bank-service --port=8080 --debug=true
```

### 3. Test It Works
```bash
# Health check
curl http://localhost:8080/health

# Check balance (test user created by schema)
curl "http://localhost:8080/balance?user_id=550e8400-e29b-41d4-a716-446655440000"

# Withdraw $20
curl -X POST -H "Content-Type: application/json" \
  -d '{"user_id":"550e8400-e29b-41d4-a716-446655440000","amount":20000}' \
  http://localhost:8080/withdraw
```

## 🎯 You're Done!

Your banking service is now running with:
- ✅ Clean Architecture
- ✅ PostgreSQL Database
- ✅ REST API
- ✅ Transaction Logging
- ✅ Input Validation
- ✅ Error Handling

## 📚 Need More?
- Full documentation: `README.md`
- Database schema: `database/schema.sql`
- API docs: See `README.md#api-documentation`

**Happy coding! 🎉**