# Tournament Management System

A full-stack web application for managing tournaments, built as a learning project to explore **Go** backend development
and **SvelteKit** frontend frameworks.

## 🎯 Purpose

This project serves as a hands-on learning experience for:

- **Go (Golang)** - Backend API development, clean architecture patterns, and database operations
- **SvelteKit** - Modern frontend framework with server-side rendering and reactive components
- **Full-stack integration** - Connecting frontend and backend systems

## 🏗️ Architecture

### Backend (Go)

- **Clean Architecture** with hexagonal design patterns
- **Domain-driven design** with clear separation of concerns
- **RESTful API** endpoints for tournament management
- **PostgreSQL** database with repository pattern
- **Dependency injection** for testability

```
backend/
├── cmd/                    # Application entry points
├── internal/
│   ├── adapters/          # External interfaces
│   │   ├── driven/        # Database repositories
│   │   └── driving/       # HTTP handlers & validation
│   ├── application/       # Business logic services
│   ├── domain/           # Core business entities
│   ├── middleware/       # HTTP middleware
│   └── ports/           # Interface definitions
├── pkg/                 # Shared utilities
└── tests/              # Test files
```

### Frontend (SvelteKit)

- **Server-side rendering** with SvelteKit
- **TypeScript** for type safety
- **Tailwind CSS** + **DaisyUI** for styling
- **Component-based architecture**

```
frontend/src/
├── lib/
│   ├── components/     # Reusable UI components
│   ├── types/         # TypeScript type definitions
│   └── utils/         # Helper functions
└── routes/           # SvelteKit file-based routing
    ├── (app)/        # Authenticated application pages
    └── (public)/     # Public pages (tournament signup)
```

## 🚀 Features

- **Tournament Management**
    - Create and configure tournaments
    - Manage tournament status (Draft, Active, Completed, Cancelled)
    - View tournament details and statistics

- **Player Management**
    - Player registration and signup
    - Qualifying rounds with time tracking
    - Player rankings and positions

- **Dashboard**
    - Overview of all tournaments
    - Filtering and sorting capabilities
    - Responsive design for mobile and desktop

## 🛠️ Tech Stack

### Backend

- **Go 1.23** - Programming language
- **Chi Router** - HTTP router
- **PostgreSQL** - Database
- **Docker** - Containerization

### Frontend

- **SvelteKit 2.22** - Frontend framework
- **TypeScript 5.0** - Type safety
- **Tailwind CSS 4.1** - Styling
- **DaisyUI 5.0** - UI components
- **Vite 7.0** - Build tool

## 📋 Prerequisites

- **Go 1.23+**
- **Node.js 18+**
- **PostgreSQL 13+**
- **Docker** (optional)

Here's the updated README section for Tilt setup:

## 🏃‍♂️ Getting Started

### Quick Start with Tilt (Recommended)

The easiest way to get the entire project running is with **Tilt**, which orchestrates the development environment with
hot reloading for both backend and frontend.

1. **Install Tilt**

```shell script
# macOS
brew install tilt
```

2. **Install Docker** (required for Tilt)
    - [Docker Desktop](https://www.docker.com/products/docker-desktop) for your platform

3. **Start the development environment**

```shell script
# From the project root
   tilt up
```

4. **Access the Tilt UI**
    - Open `http://localhost:10350` to view the Tilt dashboard
    - Monitor logs, rebuilds, and service status in real-time

5. **Access the applications**
    - Frontend: `http://localhost:5173`
    - Backend API: `http://localhost:3000`

### What Tilt Does

- **Automated Setup**: Builds and starts all services with one command
- **Hot Reloading**: Automatically rebuilds and restarts services when code changes
- **Dependency Management**: Ensures database is ready before starting the backend
- **Live Logs**: View logs from all services in one dashboard

### Tilt Commands

```shell script
# Start development environment
tilt up

# Start with specific resources only
tilt up frontend backend

# View logs for a specific service
tilt logs backend

# Stop all services
tilt down

# Force rebuild all images
tilt up --force-rebuild
```

### Manual Setup (Alternative)

If you prefer to run services individually without Tilt:

#### Backend Setup

1. **Navigate to backend directory**

```shell script
cd backend
```

2. **Install Go dependencies**

```shell script
go mod download
```

3. **Set up database**

```shell script
# Create PostgreSQL database
   createdb tournament_db
   
   # Run migrations (if available)
   go run cmd/migrate/main.go
```

4. **Start the server**

```shell script
go run cmd/api/main.go
```

Server will run on `http://localhost:3000`

#### Frontend Setup

1. **Navigate to frontend directory**

```shell script
cd frontend
```

2. **Install dependencies**

```shell script
npm install
```

3. **Start development server**

```shell script
npm run dev
```

Application will run on `http://localhost:5173`

### Docker Compose Setup (Alternative)

```shell script
# Build and run with Docker Compose
docker-compose up --build
```

## 🔧 Development Workflow

### Using Tilt (Recommended)

1. Make code changes in your editor
2. Tilt automatically detects changes and rebuilds affected services
3. View build status and logs in Tilt UI (`http://localhost:10350`)
4. Test changes in the browser (auto-refreshes on frontend changes)

### Tilt Benefits for Learning

- **Microservices Simulation**: Learn how multiple services work together
- **DevOps Practices**: Experience with containerization and orchestration
- **Development Efficiency**: Focus on code instead of environment setup
- **Production-like Environment**: Closer to real deployment scenarios

## 📋 Prerequisites

### For Tilt Setup

- **Docker Desktop** - Container platform
- **Tilt** - Development environment orchestrator

### For Manual Setup

- **Go 1.23+**
- **Node.js 18+**
- **PostgreSQL 13+**
- **Docker** (optional)

## 📚 Learning Outcomes

### Go Backend

- **Clean Architecture** implementation
- **Dependency injection** patterns
- **Repository pattern** for data access
- **HTTP middleware** and request handling
- **Error handling** and validation
- **Testing** with Go's testing package

### SvelteKit Frontend

- **Server-side rendering** concepts
- **Reactive programming** with Svelte stores
- **TypeScript integration** in Svelte
- **Component composition** and props
- **Form handling** and validation
- **API integration** and data fetching

## 🧪 Testing

### Backend Tests

```shell script
cd backend
go test ./...
```

## 📖 API Documentation

The API follows RESTful conventions:

- `GET /api/tournament` - List all tournaments
- `POST /api/tournament` - Create tournament
- `GET /api/tournament/{id}` - Get tournament details
- `PATCH /api/tournament/{id}/status` - Update tournament status
- `GET /api/tournament/{id}/qualifying` - Get qualifying results
