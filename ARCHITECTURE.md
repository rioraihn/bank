# Wallet Service - Clean Architecture Structure

## Overview
This project follows Clean Architecture principles with Domain-Driven Design (DDD) and Test-Driven Development (TDD) practices.

## Architecture Layers

### Domain Layer (`internal/domain/`)
- **Entities**: Core business objects with identity (Wallet, Transaction)
- **Value Objects**: Immutable concepts without identity (Money, UserID)
- **Repository Interfaces**: Contracts for data persistence
- **Domain Services**: Business logic that doesn't naturally fit in entities
- **Domain Errors**: Business-specific error types

### Application Layer (`internal/application/`)
- **Use Cases**: Application-specific business rules (WithdrawFunds, GetBalance)
- **DTOs**: Data transfer objects for requests/responses
- **Application Services**: Orchestration of domain objects

### Infrastructure Layer (`internal/infrastructure/`)
- **Persistence**: Database implementations of repository interfaces
- **HTTP**: API controllers, handlers, and middleware
- **Config**: Configuration management

## Key Principles
1. **Dependency Inversion**: Domain layer doesn't depend on infrastructure
2. **Single Responsibility**: Each component has one reason to change
3. **Testability**: All components can be unit tested in isolation
4. **Domain Focus**: Business rules are centralized in the domain layer

## Development Workflow
1. Write failing tests (TDD)
2. Implement domain entities and value objects
3. Create repository interfaces in domain
4. Implement use cases in application layer
5. Implement infrastructure (database, HTTP)
6. Integration tests to verify end-to-end functionality