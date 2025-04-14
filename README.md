# Go Backend Template

A production-ready Go backend template with essential features and best practices. This template provides a solid foundation for building scalable web applications with Go.

## Features

- 🏗️ **Clean Architecture**
  - Layered architecture with clear separation of concerns
  - Domain-driven design principles
  - Dependency injection pattern

- 🔐 **Authentication & Authorization**
  - JWT-based authentication
  - Role-based access control
  - Refresh token mechanism
  - Password hashing with bcrypt

- 👤 **User Management**
  - User registration and login
  - Profile management
  - Password reset functionality
  - Email verification

- 💳 **Payment Integration**
  - Stripe payment processing
  - Multiple payment methods support
  - Payment history and tracking
  - Refund handling

- 📝 **Structured Logging**
  - JSON-formatted logs
  - Log levels and rotation
  - Request/Response logging
  - Error tracking

- 🎯 **Database**
  - PostgreSQL with GORM
  - Database migrations
  - Connection pooling
  - Soft deletes

- 🔄 **Caching**
  - Redis integration
  - Cache middleware
  - Configurable TTL
  - Cache invalidation

- 🐳 **Docker Support**
  - Multi-stage builds
  - Docker Compose setup
  - Production-ready configuration
  - Easy local development

- 🧪 **Testing**
  - Unit tests
  - Integration tests
  - Mock generation
  - Test coverage

- 🔧 **Configuration**
  - Environment-based config
  - Secret management
  - Feature flags
  - Flexible app settings

- 🛠️ **Developer Tools**
  - Hot reload
  - Linting and formatting
  - API documentation
  - Make commands

- 🔍 **Monitoring**
  - Prometheus metrics
  - Health checks
  - Request tracing
  - Performance monitoring

## Project Structure

```
.
├── cmd/                    # Application entry points
│   └── api/               # API server
├── config/                # Configuration files
├── internal/              # Private application code
│   ├── auth/             # Authentication logic
│   ├── database/         # Database connections
│   ├── handlers/         # HTTP handlers
│   ├── middleware/       # HTTP middleware
│   ├── models/           # Data models
│   ├── repository/       # Data access layer
│   └── services/         # Business logic
├── migrations/           # Database migrations
├── pkg/                  # Public libraries
│   ├── logger/          # Logging package
│   ├── utils/           # Common utilities
│   └── validator/       # Validation helpers
├── scripts/             # Build and deployment scripts
├── .env.example         # Example environment variables
├── docker-compose.yml   # Docker compose configuration
├── Dockerfile          # Docker build file
├── go.mod             # Go modules file
└── Makefile          # Build automation
```

## Prerequisites

- Go 1.21 or higher
- Docker and Docker Compose
- PostgreSQL
- Redis
- Make

## Getting Started

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/go_byteeats.git
   cd go_byteeats
   ```

2. Copy the example environment file:
   ```bash
   cp .env.example .env
   ```

3. Update the environment variables in `.env`

4. Start the development environment:
   ```bash
   make dev
   ```

## Development Commands

- `make dev`: Start development server with hot reload
- `make build`: Build the application
- `make test`: Run tests
- `make test-coverage`: Run tests with coverage
- `make migrate`: Run database migrations
- `make migrate-down`: Rollback migrations
- `make lint`: Run linters
- `make fmt`: Format code
- `make docker-build`: Build Docker image
- `make docker-run`: Run Docker container
- `make generate`: Generate mocks and other code
- `make clean`: Clean build artifacts

## API Documentation

API documentation is available at `/swagger/index.html` when running in development mode.

### Main Endpoints

- **Authentication**
  - POST `/api/v1/auth/register`: Register a new user
  - POST `/api/v1/auth/login`: Login user
  - POST `/api/v1/auth/refresh`: Refresh access token
  - POST `/api/v1/auth/logout`: Logout user

- **Users**
  - GET `/api/v1/users/me`: Get current user
  - PUT `/api/v1/users/me`: Update user profile
  - POST `/api/v1/users/password`: Change password
  - POST `/api/v1/users/verify-email`: Verify email

- **Payments**
  - POST `/api/v1/payments`: Create payment
  - GET `/api/v1/payments`: List payments
  - GET `/api/v1/payments/:id`: Get payment
  - POST `/api/v1/payments/:id/refund`: Refund payment

## Configuration

The application uses environment variables for configuration. See `.env.example` for available options.

### Important Configuration Options

- `APP_ENV`: Application environment (development/staging/production)
- `APP_PORT`: HTTP server port
- `DB_*`: Database configuration
- `REDIS_*`: Redis configuration
- `JWT_*`: JWT settings
- `STRIPE_*`: Stripe integration settings

## Deployment

### Docker Deployment

1. Build the Docker image:
   ```bash
   make docker-build
   ```

2. Run the container:
   ```bash
   make docker-run
   ```

### Manual Deployment

1. Build the binary:
   ```bash
   make build
   ```

2. Run migrations:
   ```bash
   make migrate
   ```

3. Start the server:
   ```bash
   ./bin/api
   ```

## Contributing

1. Fork the repository
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Create a new Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgments

- [Gin Web Framework](https://github.com/gin-gonic/gin)
- [GORM](https://gorm.io)
- [golang-migrate](https://github.com/golang-migrate/migrate)
- [zerolog](https://github.com/rs/zerolog)
- [Stripe Go Library](https://github.com/stripe/stripe-go) 