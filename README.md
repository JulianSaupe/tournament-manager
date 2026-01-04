<div align="center">

# 🏆 Tournament

**A modern tournament management system built with microservices architecture**

[![Go](https://img.shields.io/badge/Go-1.23-00ADD8?style=flat&logo=go)](https://go.dev/)
[![Rust](https://img.shields.io/badge/Rust-2024-000000?style=flat&logo=rust)](https://www.rust-lang.org/)
[![SvelteKit](https://img.shields.io/badge/SvelteKit-2-FF3E00?style=flat&logo=svelte)](https://kit.svelte.dev/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-18-316192?style=flat&logo=postgresql)](https://www.postgresql.org/)

</div>

---

## 📋 Overview

Tournament is a full-stack tournament management platform featuring a microservices backend with gRPC communication, a
modern SvelteKit frontend, and robust authentication services.

## 🛠️ Tech Stack

### Backend Services

#### **Engine** (Go)

- **Language**: Go 1.23
- **Router**: Chi v5
- **Database**: PostgreSQL 18

#### **Identity Service** (Rust)

- **Language**: Rust (Edition 2021)
- **Framework**: Tonic (gRPC)
- **Database**: SQLx with PostgreSQL
- **Runtime**: Tokio async
- **Database**: PostgreSQL 18

### Frontend

- **Framework**: SvelteKit 2 + Svelte 5
- **Language**: TypeScript 5
- **Styling**: Tailwind CSS 4 + DaisyUI

### Infrastructure

- **Development**: Tilt (hot reload orchestration)
- **Containers**: Docker Compose
- **Databases**: Dual PostgreSQL 18 instances
- **Communication**: gRPC + Protocol Buffers

## 📁 Project Structure

```
Tournament/
├── backend/
│   ├── engine/              # Go-based tournament engine (REST API)
│   │   ├── cmd/api/        # Main application entry
│   │   └── ...
│   └── identity-service/    # Rust-based auth service (gRPC)
│       ├── proto/          # Protocol Buffer definitions
│       ├── migrations/     # Database migrations
│       └── src/
├── frontend/                # SvelteKit web application
│   ├── src/
│   │   ├── lib/
│   │   └── routes/
│   └── static/
└── docker/                  # Docker Compose setup
    └── docker-compose.yml
```

## 🚀 Getting Started

### Prerequisites

- [Go 1.23+](https://go.dev/dl/)
- [Rust (2024 edition)](https://rustup.rs/)
- [Node.js](https://nodejs.org/)
- [Docker](https://www.docker.com/)
- [Tilt](https://tilt.dev/)

### Development

Start all services with hot reload:

```bash
tilt up
```

**Available Services:**

| Service         | URL                   | Description                  |
|-----------------|-----------------------|------------------------------|
| 🌐 Frontend     | http://localhost:5173 | SvelteKit web application    |
| 🔧 API (Engine) | http://localhost:3000 | Tournament management API    |
| 🔐 Identity     | http://localhost:5000 | gRPC authentication service  |
| 🗄️ Engine DB   | localhost:5432        | PostgreSQL (tournament data) |
| 🗄️ Identity DB | localhost:5433        | PostgreSQL (user data)       |

---

<div align="center">

**Built with ❤️ using Go, Rust, and SvelteKit**

</div>
