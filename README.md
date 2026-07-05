# SuProxy Backend

Production-grade VPN platform backend built with Go, following Clean Architecture and Domain-Driven Design principles.

## 🏗️ Architecture

This project follows **Clean Architecture** and **Domain-Driven Design (DDD)** principles:

```
├── cmd/
│   └── api/              # Application entry point
├── internal/
│   ├── domain/           # Business logic & entities (planned)
│   ├── application/      # Use cases & application logic (planned)
│   ├── infrastructure/   # External dependencies
│   │   ├── config/       # Configuration management
│   │   ├── logger/       # Logging system
│   │   └── server/       # HTTP server setup
│   └── interfaces/       # API handlers & DTOs (planned)
├── configs/              # Configuration files
├── migrations/           # Database migrations (planned)
└── docs/                 # Documentation (planned)
```

## 🚀 Tech Stack

- **Go 1.21** - Programming language
- **Gin** - HTTP web framework
- **PostgreSQL** - Database
- **GORM** - ORM
- **Docker & Docker Compose** - Containerization
- **Viper** - Configuration management
- **Zap** - Structured logging
- **Air** - Hot reload for development

## 📋 Prerequisites

- Go 1.21 or higher
- Docker & Docker Compose
- Make (optional, for using Makefile commands)

## 🛠️ Getting Started

### 1. Clone the repository

```bash
git clone <repository-url>
cd suproxy-backend
```

### 2. Environment Configuration

Copy the example environment file:

```bash
cp .env.example .env
```

### 3. Run with Docker

Start all services:

```bash
make docker-up
```

Or manually:

```bash
docker-compose up -d
```

### 4. Verify Installation

Check health endpoint:

```bash
curl http://localhost:8080/health
```

Expected response:

```json
{
  "status": "healthy",
  "service": "suproxy-backend",
  "version": "1.0.0"
}
```

## 📝 Available Commands

```bash
make help           # Display all available commands
make build          # Build the application
make run            # Run the application
make dev            # Run with hot reload (requires air)
make docker-build   # Build Docker image
make docker-up      # Start Docker containers
make docker-down    # Stop Docker containers
make docker-logs    # Show Docker logs
make test           # Run tests
make lint           # Run linter
make clean          # Clean build artifacts
```

## 🔌 API Endpoints

### Health Check

- **GET** `/health` - Service health check
- **GET** `/api/v1/ping` - Simple ping endpoint

## 🏗️ Project Structure Details

### Clean Architecture Layers

1. **Domain Layer** (Planned)
   - Business entities
   - Business rules
   - Domain services

2. **Application Layer** (Planned)
   - Use cases
   - Application services
   - DTOs

3. **Infrastructure Layer** (Current)
   - Database implementation
   - External services
   - Configuration
   - Logging

4. **Interface Layer** (Planned)
   - HTTP handlers
   - Request/Response models
   - Middleware

## 🔐 Security

- Environment variables for sensitive data
- Structured logging (no sensitive data in logs)
- Graceful shutdown handling
- CORS support (planned)
- JWT authentication (planned)

## 🧪 Development

### Local Development with Hot Reload

Install Air:

```bash
go install github.com/cosmtrek/air@latest
```

Run with hot reload:

```bash
make dev
```

### Database Migrations

Database migration support with golang-migrate (planned).

## 📊 Monitoring & Logging

- Structured JSON logging with Zap
- Request/Response logging middleware
- Log levels: debug, info, warn, error, fatal
- Configurable log format (JSON/Console)

## 🔄 Version

Current version: **1.0.0**

## 🛣️ Roadmap

- [ ] Phase 1: Infrastructure Setup ✅ (Current)
- [ ] Phase 2: User Management System
- [ ] Phase 3: Authentication & Authorization (JWT)
- [ ] Phase 4: Payment Integration
- [ ] Phase 5: VPN Server Management (Xray)
- [ ] Phase 6: Admin Panel APIs
- [ ] Phase 7: Android App Integration

## 📱 Mobile App Integration

This backend is designed to integrate with the SuProxy Android application. The architecture is scalable and extensible to support future mobile requirements.

## 🤝 Contributing

This is a commercial project. Contribution guidelines will be added later.

## 📄 License

Proprietary - All rights reserved

## 👤 Author

SuProxy Team

---

**Status:** 🟢 Phase 1 Complete - Infrastructure Setup

**Last Updated:** 2026-07-05
