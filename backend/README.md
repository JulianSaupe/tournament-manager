# Tournament Backend

This project follows the Hexagonal Architecture (Ports and Adapters) pattern.

## Project Structure

```
backend/
├── cmd/                    # Application entry points
│   └── api/                # HTTP API entry point
│       └── main.go         # Main application file
├── internal/               # Private application code
│   ├── domain/             # Domain models and business logic
│   ├── application/        # Application services and use cases
│   ├── ports/              # Interfaces defining how to interact with the application
│   │   ├── input/          # Primary/driving ports (interfaces the application exposes)
│   │   └── output/         # Secondary/driven ports (interfaces the application uses)
│   ├── adapters/           # Implementations of the ports
│   │   ├── driving/        # Primary/driving adapters (e.g., HTTP handlers, middleware)
│   │   └── driven/         # Secondary/driven adapters (e.g., database repositories)
│   └── config/             # Application configuration
```

## Hexagonal Architecture Overview

The Hexagonal Architecture (also known as Ports and Adapters) is a software design pattern that allows an application to be equally driven by users, programs, automated tests, or batch scripts, and to be developed and tested in isolation from its eventual run-time devices and databases.

## Project Components

### Domain Layer (`internal/domain`)

Contains the core business entities and logic:
- Tournament entity and related types
- Participant entity
- User entity
- Domain-specific errors

### Application Layer (`internal/application`)

Implements the business use cases:
- Tournament service implementation
- User service implementation
- Participants service implementation

### Ports Layer (`internal/ports`)

Defines interfaces for communication:
- **Input Ports (`internal/ports/input`)**: Interfaces exposed by the application
  - Tournament service interface
  - User service interface
  - Participants service interface
- **Output Ports (`internal/ports/output`)**: Interfaces used by the application
  - Tournament repository interface
  - User repository interface
  - Participants repository interface

### Adapters Layer (`internal/adapters`)

Implements the interfaces defined in the ports layer:
- **Driving Adapters (`internal/adapters/driving`)**: Handle incoming requests
  - HTTP handlers for tournaments and participants
  - Middleware for authentication, response formatting, etc.
  - Request/response models
- **Driven Adapters (`internal/adapters/driven`)**: Connect to external systems
  - PostgreSQL repositories for tournaments, users, and participants

### Configuration (`internal/config`)

Contains application configuration:
- Database configuration and connection setup

### Entry Point (`cmd/api`)

The main application entry point:
- Wires together all components
- Configures the HTTP server
- Sets up middleware
- Starts the application

## Database Configuration

The project uses PostgreSQL as its database. Database configuration is handled in the `config` package, specifically in the `database.go` file.

The database connection is established using the [Bun](https://github.com/uptrace/bun) library, which is a SQL-first database library for Go.

Database initialization happens in the main application entry point (`cmd/api/main.go`), where repositories are created and injected into the application services.

## Benefits of Hexagonal Architecture

- **Separation of concerns:** Each layer has a specific responsibility
- **Testability:** The application core can be tested in isolation from external dependencies
- **Flexibility:** Implementation details can be changed without affecting the core business logic
- **Clear boundaries:** The architecture defines clear boundaries between different parts of the application
- **Independence from frameworks:** The application core is not tied to any specific framework or technology
