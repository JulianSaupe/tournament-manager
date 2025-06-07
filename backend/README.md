# Tournament Backend

This project follows the Hexagonal Architecture (Ports and Adapters) pattern.

## Project Structure

```
backend/
├── cmd/                    # Application entry points
│   ├── api/                # HTTP API entry point
│   │   └── main.go         # Main application file
│   └── migrate/            # Database migration script
│       └── main.go         # Migration file generator
├── internal/               # Private application code
│   ├── domain/             # Domain models and business logic
│   ├── application/        # Application services and use cases
│   ├── ports/              # Interfaces defining how to interact with the application
│   │   ├── input/          # Primary/driving ports (interfaces the application exposes)
│   │   └── output/         # Secondary/driven ports (interfaces the application uses)
│   └── adapters/           # Implementations of the ports
│       ├── driving/        # Primary/driving adapters (e.g., HTTP handlers, CLI)
│       └── driven/         # Secondary/driven adapters (e.g., database repositories)
├── migrations/             # Database migration files
│   ├── *.up.sql            # Migration files for upgrading the database
│   └── *.down.sql          # Migration files for downgrading the database
└── pkg/                    # Public libraries that can be used by other applications
```

## Hexagonal Architecture Overview

The Hexagonal Architecture (also known as Ports and Adapters) is a software design pattern that allows an application to be equally driven by users, programs, automated tests, or batch scripts, and to be developed and tested in isolation from its eventual run-time devices and databases.

### Key Components

1. **Domain**: Contains the business logic and domain models.
2. **Application**: Contains the use cases and application services.
3. **Ports**: Defines interfaces for communication with the outside world.
   - **Primary/Driving Ports**: Interfaces that the application exposes.
   - **Secondary/Driven Ports**: Interfaces that the application uses to communicate with external systems.
4. **Adapters**: Implements the ports.
   - **Primary/Driving Adapters**: Implements the primary ports (e.g., HTTP handlers, CLI).
   - **Secondary/Driven Adapters**: Implements the secondary ports (e.g., database repositories, external services).

## Detailed Component Description

### Internal Directory

The internal directory contains the private application code that is not intended to be imported by other applications. It follows the Go convention for organizing private code.

#### Domain Layer

The domain layer is the core of the application and contains:

- Domain entities (e.g., Tournament)
- Value objects
- Domain services
- Domain events
- Business rules and logic

**Principles:**
- The domain layer should be independent of any external concerns
- It should not have dependencies on other layers
- It should contain the business rules and logic of the application
- It should be the most stable part of the application

**Example:** The `tournament.go` file defines the Tournament entity and its associated types, which represent the core domain concepts of the application.

#### Application Layer

The application layer is responsible for:

- Implementing use cases
- Orchestrating the flow of data to and from the domain entities
- Coordinating with the ports to interact with external systems
- Applying application-specific business rules

**Principles:**
- The application layer depends on the domain layer and ports
- It should not have dependencies on the adapters
- It should not contain business rules that belong to the domain layer
- It should be independent of the delivery mechanism (e.g., HTTP, CLI)

**Example:** The `tournament_service_impl.go` file implements the TournamentService interface defined in the input ports. It uses the TournamentRepository interface defined in the output ports to interact with the data storage.

#### Ports Layer

The ports layer defines interfaces for how the application interacts with the outside world:

- **Input Ports (Driving Ports):** Interfaces that the application exposes to be used by external actors
  - Define the API that the application exposes to the outside world
  - Are implemented by the application layer
  - Are used by the primary adapters (e.g., HTTP handlers, CLI commands)

  **Principles:**
  - Input ports should be defined in terms of the domain model
  - They should be independent of the specific technology used to implement them
  - They should represent use cases or user stories
  - They should be stable and change less frequently than the adapters

  **Example:** The `tournament_service.go` file defines the TournamentService interface, which represents the operations that can be performed on tournaments.

- **Output Ports (Driven Ports):** Interfaces that the application uses to interact with external systems
  - Define how the application interacts with external systems (e.g., databases, external services)
  - Are used by the application layer
  - Are implemented by the secondary adapters

  **Principles:**
  - Output ports should be defined in terms of the domain model
  - They should be independent of the specific technology used to implement them
  - They should represent the operations that the application needs to perform on external systems
  - They should be stable and change less frequently than the adapters

  **Example:** The `tournament_repository.go` file defines the TournamentRepository interface, which represents the operations that the application can perform on the tournament data store.

#### Adapters Layer

The adapters layer is responsible for:

- Implementing the interfaces defined in the ports layer
- Connecting the application to external systems
- Translating between the application's domain model and the external world

**Principles:**
- Adapters should be replaceable without changing the application core
- They should be specific to a particular technology or external system
- They should be isolated from each other
- They should be testable independently of the application core

- **Primary/Driving Adapters:** Implement the input ports
  - Handle incoming requests from the outside world
  - Convert external data formats to the domain model
  - Invoke the appropriate use cases in the application layer
  - Convert the results back to the external format

  **Example:** The `http/tournament_handler.go` file implements an HTTP adapter for the tournament operations.

- **Secondary/Driven Adapters:** Implement the output ports
  - Implement the interfaces defined in the output ports
  - Connect the application to external systems (e.g., databases, external services)
  - Convert between the domain model and the external data format
  - Handle the technical details of interacting with external systems

  **Example:** The `memory/tournament_repository.go` file implements an in-memory repository for tournaments.

### Command Directory

The command directory contains the entry points for the application:

- Provides entry points for the application
- Wires together the different components of the application
- Configures the application
- Starts the application

**Principles:**
- Each subdirectory should represent a different entry point or executable
- The code in this directory should be minimal and focused on wiring together the components
- It should not contain business logic
- It should depend on the application, domain, and adapters layers

**Example:** The `api/main.go` file is the entry point for the HTTP API. It creates the necessary components and wires them together to start the HTTP server.

### Package Directory

The package directory contains public libraries that can be used by other applications:

- Provides reusable libraries that can be imported by other applications
- Contains code that is not specific to the Tournament application
- Offers utility functions, helpers, and shared components

**Principles:**
- Code in this directory should be stable and well-tested
- It should have minimal dependencies
- It should be designed for reuse
- It should have clear documentation and examples

**Note:** This directory is currently empty but can be populated as the application grows and common patterns emerge that can be extracted into reusable libraries.

## Database Migrations

This project uses [golang-migrate/migrate](https://github.com/golang-migrate/migrate) for database migrations. Migrations are SQL scripts that are used to create, modify, or drop database objects like tables, indexes, etc.

### Migration Files

Migration files are stored in the `migrations` directory and follow the naming convention:

- `<version>_<description>.up.sql`: Scripts for upgrading the database
- `<version>_<description>.down.sql`: Scripts for downgrading the database

For example:
- `000001_create_tournaments_table.up.sql`: Creates the tournaments table
- `000001_create_tournaments_table.down.sql`: Drops the tournaments table

### Migration Execution

Migrations are automatically executed when the application starts. The migration process is handled by the `MigrateDatabase` function in the `config` package.

### Migration Commands

The project includes a migration script with several commands for managing database migrations.

#### Generate a New Migration

To create a new migration:

1. Navigate to the project root directory
2. Run the migration script with the `generate` command and a descriptive name for your migration:

```bash
# From the repository root
go run backend/cmd/migrate/main.go generate create_users_table
```

This will create two files in the `migrations` directory:
- `VersionYYYYMMDDHHMMSS.up.sql`: For upgrading the database
- `VersionYYYYMMDDHHMMSS.down.sql`: For downgrading the database

Where `YYYYMMDDHHMMSS` is the current timestamp (e.g., `20230607145623`).

3. Edit the generated files to add your SQL statements

#### Execute All Migrations

To execute all pending migrations:

```bash
# From the repository root
go run backend/cmd/migrate/main.go migrate
```

This will execute all migrations that haven't been applied to the database yet.

#### Execute a Specific Migration

To execute a specific migration:

```bash
# Execute in the up direction (apply the migration)
go run backend/cmd/migrate/main.go execute <migration_name> --up

# Execute in the down direction (revert the migration)
go run backend/cmd/migrate/main.go execute <migration_name> --down
```

Replace `<migration_name>` with the name of the migration you want to execute.

Note: Migrations are also automatically executed when the application starts.

#### Manual Creation (Alternative)

Alternatively, you can manually create migration files:

1. Create a new migration file in the `migrations` directory with a higher version number than the existing migrations
2. Write the SQL statements for the migration in the `.up.sql` file
3. Write the SQL statements to revert the migration in the `.down.sql` file
4. The migrations will be executed automatically when the application starts

## Benefits of Hexagonal Architecture

- **Separation of concerns:** Each layer has a specific responsibility
- **Testability:** The application core can be tested in isolation from external dependencies
- **Flexibility:** Implementation details can be changed without affecting the core business logic
- **Clear boundaries:** The architecture defines clear boundaries between different parts of the application
- **Independence from frameworks:** The application core is not tied to any specific framework or technology
