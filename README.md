# DDD Bottom-Up Implementation

A sample Go project implementing Domain-Driven Design (DDD) using a bottom-up approach, following Clean Architecture principles with clear layer separation.

## 📋 Table of Contents

- [Architecture](#architecture)
- [Project Structure](#project-structure)
- [Development Setup](#development-setup)
- [API Endpoints](#api-endpoints)
- [Testing](#testing)
- [Tech Stack](#tech-stack)
- [Design Patterns](#design-patterns)

## 🏗️ Architecture

This project adopts Clean Architecture + DDD layer structure:

```
┌─────────────────────────────┐
│     Presentation Layer      │  ← HTTP handlers, routing
├─────────────────────────────┤
│      Use Case Layer         │  ← Application services
├─────────────────────────────┤
│       Domain Layer          │  ← Entities, Value Objects, Services
├─────────────────────────────┤
│    Infrastructure Layer     │  ← Repository implementations
└─────────────────────────────┘
```

### Domain Layer
- **Entities**: `User`, `Circle`, `CircleMembers`, `Shipment`
- **Value Objects**: `Email`, `FullName`, `CircleName`, `Money`
- **Specifications**: `CircleMemberLimitSpecification`, `RecommendedCircleSpecification`
- **Repository Interfaces**: Data access contracts
- **Domain Services**: Business logic that doesn't belong to entities

### Use Case Layer
Application services orchestrating domain operations:
- `CreateUserUseCase` - User creation
- `GetUserUseCase` - User retrieval
- `UpdateUserUseCase` - User updates
- `DeleteUserUseCase` - User deletion

### Infrastructure Layer
- **Repository Implementations**: Memory-based and database implementations
- **Database**: SQL schema definitions

### Presentation Layer
- **Handlers**: HTTP request/response handling
- **Router**: HTTP route configuration using go-chi

## 📁 Project Structure

```
ddd-bottomup/
├── domain/                 # Domain layer
│   ├── user.go            # User entity & services
│   ├── values.go          # Value objects
│   ├── *_test.go          # Domain tests
│   └── ...
├── usecase/               # Use case layer
│   ├── create_user.go     # Create user use case
│   ├── get_user.go        # Get user use case
│   ├── update_user.go     # Update user use case
│   ├── delete_user.go     # Delete user use case
│   └── *_test.go          # Use case tests
├── infrastructure/        # Infrastructure layer
│   ├── memory_user_repository.go  # In-memory repository
│   └── mysql_user_repository.go   # MySQL repository
├── presentation/          # Presentation layer
│   ├── user_handler.go    # User handlers
│   └── router.go          # Router configuration
├── migrations/            # Database migrations
│   └── 000001_initial_schema.sql
├── main.go               # Application entry point
├── go.mod                # Go module configuration
├── CLAUDE.md             # Development instructions
└── README.md             # This file
```

## 🚀 Development Setup

### Prerequisites
- Go 1.24 or later
- MySQL (for production)

### Installation

```bash
# Clone the repository
git clone <repository-url>
cd ddd-bottomup

# Download dependencies
go mod tidy

# Build the application
go build -o bin/app .

# Run the application
go run main.go
```

The server will start at `http://localhost:8080`.

## 📡 API Endpoints

| Method | Endpoint     | Description |
|--------|-------------|-------------|
| POST   | `/users`     | Create user |
| GET    | `/users/{id}` | Get user |
| PUT    | `/users/{id}` | Update user |
| DELETE | `/users/{id}` | Delete user |
| GET    | `/health`    | Health check |

### Request Examples

#### Create User
```bash
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{
    "firstName": "John",
    "lastName": "Doe", 
    "email": "john@example.com",
    "isPremium": false
  }'
```

#### Get User
```bash
curl http://localhost:8080/users/{user-id}
```

#### Update User
```bash
curl -X PUT http://localhost:8080/users/{user-id} \
  -H "Content-Type: application/json" \
  -d '{
    "firstName": "Jane",
    "lastName": "Smith"
  }'
```

## 🧪 Testing

### Run All Tests
```bash
go test ./...
```

### Run Tests by Layer
```bash
# Domain layer tests
go test ./domain -v

# Use case layer tests
go test ./usecase -v
```

### Test Strategy
- **Domain Tests**: Unit tests for entities, value objects, and domain services
- **Use Case Tests**: Integration tests for application services (using memory repositories)
- **Error Case Tests**: Comprehensive error handling validation

### Test Coverage
- Value object validation (Email, FullName, CircleName, Money)
- Entity behavior (User creation, updates, equality)
- Domain service logic (UserExistenceService)
- Use case orchestration with error scenarios
- Typed error handling throughout all layers

## 🔧 Tech Stack

- **Language**: Go 1.24
- **Web Framework**: go-chi/chi v5 (lightweight HTTP router)
- **Database**: MySQL (production) / In-memory (development/test)
- **UUID Generation**: google/uuid
- **Architecture**: Clean Architecture + DDD

### Dependencies
```go
require (
    github.com/go-chi/chi/v5 v5.2.3
    github.com/go-sql-driver/mysql v1.9.3
    github.com/google/uuid v1.6.0
)
```

## 🏛️ Design Patterns

### Specification Pattern
For expressing complex business rules:
```go
type CircleSpecification interface {
    IsSatisfiedBy(circle *Circle) bool
}
```

### Repository Pattern
Data access abstraction:
```go
type UserRepository interface {
    FindByID(id *UserID) (*User, error)
    FindByName(name *FullName) (*User, error)
    Save(user *User) error
    Delete(id *UserID) error
}
```

### Entity Design
Strong typing and reconstruction patterns:
- Dedicated ID types (`UserID`, `CircleID`)
- Reconstruction patterns for rebuilding from storage
- Value objects for data integrity

### Error Handling
Typed errors with automatic HTTP status mapping:
```go
type DomainError interface {
    error
    HTTPStatus() int
}

// Example implementations
type UserNotFoundError struct { ID string }
type DuplicateUserNameError struct { Name string }
type InvalidEmailError struct { Value string }
```

### Middleware
go-chi middleware for cross-cutting concerns:
- Request logging
- Panic recovery
- Request timeout (60s)
- Content-Type headers

## 📊 Key Features

### Domain-Driven Design
- Rich domain models with behavior
- Value objects ensuring invariants
- Domain services for complex operations
- Specifications for business rules

### Clean Architecture
- Dependency inversion (interfaces pointing inward)
- Independent testability of each layer
- Framework-independent business logic

### Error Handling
- Type-safe error propagation
- HTTP status code mapping
- Detailed error context

### Testing
- Comprehensive unit and integration tests
- Mock implementations for testing
- Table-driven test patterns

## 🚀 Production Deployment

1. Set up database
2. Configure environment variables
3. Run migrations
4. Start application

```bash
# Run migrations
mysql -u user -p database < migrations/000001_initial_schema.sql

# Production start
./bin/app
```

## 🤝 Contributing

This project serves as a learning sample for DDD and Clean Architecture best practices in Go. Feel free to explore the code structure and patterns implemented.

---

This project demonstrates DDD and Clean Architecture implementation patterns for educational purposes.