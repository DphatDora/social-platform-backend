# Social Platform Backend

[![CI/CD Pipeline](https://github.com/DphatDora/social-platform-backend/actions/workflows/ci-cd.yml/badge.svg)](https://github.com/DphatDora/social-platform-backend/actions/workflows/ci-cd.yml)
[![Go Version](https://img.shields.io/badge/Go-1.24.4-00ADD8?logo=go)](https://go.dev/)
[![Test Coverage](https://img.shields.io/badge/coverage-20.2%25-yellow)](./docs/TEST_SUMMARY_PHASE3.md)

A RESTful API backend for a social media platform built with Go, featuring user authentication, communities, posts, comments, real-time notifications, and direct messaging.

## Features

- **User Management**: Registration, authentication with JWT, profile management, password reset
- **Community System**: Create and manage communities with moderator roles and permissions
- **Content Management**: Posts (text, link, media), comments with threading, voting system
- **Real-time Features**: Server-Sent Events (SSE) for live notifications and messages
- **Messaging**: Direct messaging between users with read receipts and conversation management
- **Notifications**: Customizable notification settings with push and email options
- **Recommendation Engine**: Personalized post recommendations based on user interests and tags
- **Search**: Full-text search for posts, comments, communities, and users

## Tech Stack

- **Language**: Go 1.24.4
- **Web Framework**: Gin
- **Database**: PostgreSQL with GORM ORM
- **Cache**: Redis for session management and rate limiting
- **Authentication**: JWT with Google OAuth2 integration
- **Testing**: testify/mock for unit tests with 20.2% service coverage
- **Deployment**: Render with automated CI/CD via GitHub Actions

## Project Structure

```text
social-platform-backend/
├── src/
│   ├── cmd/
│   │   └── server/          # Application entry point
│   ├── config/              # Configuration management
│   ├── internal/
│   │   ├── domain/
│   │   │   ├── model/       # Database models
│   │   │   └── repository/  # Repository interfaces
│   │   ├── infrastructure/
│   │   │   ├── cache/       # Redis implementation
│   │   │   └── db/          # PostgreSQL implementation
│   │   ├── interface/
│   │   │   ├── dto/         # Request/Response DTOs
│   │   │   ├── handler/     # HTTP handlers
│   │   │   ├── middleware/  # Authentication, CORS, etc.
│   │   │   └── router/      # Route definitions
│   │   └── service/         # Business logic layer
│   ├── package/
│   │   ├── constant/        # Application constants
│   │   ├── err/             # Custom error types
│   │   ├── template/        # Email and notification templates
│   │   └── util/            # Helper utilities
│   └── go.mod
├── docs/                    # Documentation and diagrams
├── .github/workflows/       # CI/CD workflows
└── scripts/                 # Utility scripts
```

## Getting Started

### Prerequisites

- Go 1.24.4 or higher
- PostgreSQL 14+
- Redis 7+
- Git

### Installation

Clone the repository:

```bash
git clone https://github.com/DphatDora/social-platform-backend.git
cd social-platform-backend
```

Set up environment configuration:

```bash
cd src
cp .env.example .env
cp config.yaml.example config.yaml
```

Configure database and Redis in `.env`:

```env
DATABASE_URL=postgresql://user:password@localhost:5432/social_platform
REDIS_URL=redis://localhost:6379
JWT_SECRET=your-secret-key
GOOGLE_CLIENT_ID=your-google-client-id
GOOGLE_CLIENT_SECRET=your-google-client-secret
```

Install dependencies:

```bash
go mod download
```

### Running the Application

Development mode:

```bash
cd src
go run ./cmd/server
```

Production build:

```bash
cd src
go build -o ../server ./cmd/server
cd ..
./server
```

The API will be available at `http://localhost:8080` (or configured port).
