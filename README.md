# DeployEase

> Transform complex deployment workflows into simple, one-click operations

DeployEase is an open-source self-hosting platform that provides Heroku-like simplicity on your own infrastructure. Built for developers who want the power of their own infrastructure without the operational overhead.

## 🚀 Features

### Git-to-Production Pipeline

- Automatic deployment triggered by Git pushes
- Support for multiple Git providers (GitHub, GitLab, Bitbucket)
- Branch-based deployments with environment isolation
- Rollback capabilities for failed deployments

### Infrastructure as Code

- Automated server provisioning and configuration
- Support for multiple cloud providers (AWS, DigitalOcean, Linode)
- Container orchestration with Docker
- Load balancing and auto-scaling capabilities

### Database Management

- One-click database provisioning (PostgreSQL, MySQL, Redis, MongoDB)
- Automated backups and restore functionality
- Database migration management with Goose
- Type-safe database operations using sqlc

### SSL & Domain Management

- Automatic HTTPS certificate generation via Let's Encrypt
- Custom domain routing and DNS management
- Subdomain provisioning for applications
- SSL certificate renewal automation

### Real-time Monitoring

- Live application logs with filtering and search
- Performance metrics and resource usage tracking
- Health checking and uptime monitoring
- Alert system for critical issues

### Team Collaboration

- Multi-user support with role-based permissions
- Project sharing and access control
- Audit logs for deployment activities
- Team dashboard for project overview

## 🏗️ Architecture

### Backend (Golang)

- **Framework**: Go 1.24 with bunrouter
- **API**: Huma v2 for automatic OpenAPI documentation
- **Database**: PostgreSQL 16 with Goose migrations and sqlc
- **Cache**: Dragonfly for high-performance caching
- **Config**: Viper for configuration management
- **WebSocket**: Real-time updates using coder/websocket
- **Auth**: JWT-based authentication with refresh tokens

### Frontend (Inertia.js + React)

- **Framework**: React 19 with TypeScript
- **Routing**: Inertia.js for server-side routing
- **Build Tool**: Vite for fast development and building
- **UI**: shadcn/ui component library with Tailwind CSS
- **WebSocket**: react-use-websocket for real-time updates
- **Forms**: react-hook-form for form handling
- **State**: Zustand for state management
- **Testing**: Jest for unit tests, Storybook for component development

## 📋 Prerequisites

- **Go**: 1.24 or later
- **Node.js**: 18.0 or later
- **PostgreSQL**: 16 or later
- **Docker**: For containerization and deployment
- **Task**: For running build tasks ([Installation guide](https://taskfile.dev/installation/))

## 🚀 Quick Start

### 1. Clone the Repository

```bash
git clone <repository-url>
cd deployease
```

### 2. Environment Setup

```bash
# Copy environment template
cp .env.example .env

# Edit .env with your configuration
# At minimum, configure database and JWT secrets
```

### 3. Install Dependencies

```bash
# Install all dependencies
task install

# Or install individually
task install:backend
task install:frontend
```

### 4. Database Setup

```bash
# Start PostgreSQL (using Docker)
docker run -d \
  --name deployease-postgres \
  -e POSTGRES_DB=deployease_dev \
  -e POSTGRES_USER=deployease \
  -e POSTGRES_PASSWORD=deployease_password \
  -p 5432:5432 \
  postgres:16-alpine

# Run migrations
task db:migrate

# Seed initial data (optional)
task db:seed
```

### 5. Start Development

```bash
# Start both backend and frontend
task dev

# Or start individually
task dev:backend  # Backend on :8080
task dev:frontend # Frontend on :3000
```

### 6. Access the Application

- **Frontend**: http://localhost:3000
- **Backend API**: http://localhost:8080
- **API Documentation**: http://localhost:8080/docs

## 🛠️ Development

### Available Tasks

```bash
# Development
task dev              # Start development environment
task dev:backend      # Start backend only
task dev:frontend     # Start frontend only

# Building
task build            # Build all components
task build:backend    # Build backend binary
task build:frontend   # Build frontend assets

# Testing
task test             # Run all tests
task test:backend     # Run backend tests
task test:frontend    # Run frontend tests

# Database
task db:migrate       # Run database migrations
task db:seed          # Seed database
task db:reset         # Reset database

# Docker
task docker:build     # Build Docker images
task docker:up        # Start Docker containers
task docker:down      # Stop Docker containers

# Linting
task lint             # Run all linters
task lint:backend     # Lint backend code
task lint:frontend    # Lint frontend code
```

### Project Structure

```
├── backend/                 # Go backend application
│   ├── cmd/                # Application entry points
│   ├── internal/           # Private application code
│   ├── migrations/         # Database migrations
│   ├── queries/            # SQL queries for sqlc
│   └── pkg/                # Public packages
├── frontend/               # React frontend application
│   ├── src/               # Source code
│   ├── public/            # Static assets
│   └── stories/           # Storybook stories
├── docker/                # Docker configurations
├── docs/                  # Project documentation
└── scripts/               # Build and deployment scripts
```

### Code Generation

```bash
# Generate SQL code using sqlc
cd backend && task sqlc:generate

# Generate Go code
cd backend && task generate
```

### Database Migrations

```bash
# Create a new migration
cd backend && task db:migrate:create -- create_users_table

# Run migrations
task db:migrate
```

## 🐳 Docker Development

### Using Docker Compose

```bash
# Start all services
task docker:up

# View logs
docker-compose -f docker/docker-compose.yml logs -f

# Stop services
task docker:down
```

### Services

- **Backend**: http://localhost:8080
- **Frontend**: http://localhost:3000
- **PostgreSQL**: localhost:5432
- **Dragonfly**: localhost:6379
- **Nginx**: http://localhost:80

## 📚 API Documentation

The API is automatically documented using Huma v2. Once the backend is running, visit:

- **Interactive Docs**: http://localhost:8080/docs
- **OpenAPI Spec**: http://localhost:8080/openapi.json

## 🧪 Testing

### Backend Testing

```bash
# Run all tests
cd backend && task test

# Run with coverage
cd backend && task test:coverage

# Run specific package
cd backend && go test ./internal/auth
```

### Frontend Testing

```bash
# Run all tests
cd frontend && task test

# Run in watch mode
cd frontend && task test:watch

# Run with coverage
cd frontend && task test:coverage
```

### Component Development

```bash
# Start Storybook
cd frontend && task storybook

# Build Storybook
cd frontend && task storybook:build
```

## 🚀 Deployment

### Production Build

```bash
# Build for production
task build

# Build Docker images
task docker:build

# Deploy
task deploy
```

### Environment Variables

See `.env.example` for all available configuration options. Key variables:

- `DATABASE_URL`: PostgreSQL connection string
- `JWT_SECRET`: Secret for JWT token signing
- `DRAGONFLY_URL`: Dragonfly/Redis connection string
- `SERVER_PORT`: Backend server port
- `VITE_API_URL`: Frontend API URL

## 🤝 Contributing

We welcome contributions! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

### Development Workflow

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Run the test suite
6. Submit a pull request

### Code Style

- **Go**: Follow standard Go conventions, use `gofmt`
- **TypeScript/React**: Use ESLint and Prettier configurations
- **Commits**: Use conventional commit messages

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🆘 Support

- **Documentation**: [docs/](docs/)
- **Issues**: [GitHub Issues](https://github.com/Jesuloba-world/deployease/issues)
- **Discussions**: [GitHub Discussions](https://github.com/Jesuloba-world/deployease/discussions)

## 🗺️ Roadmap

### Phase 1: MVP Foundation ✅

- [x] User authentication and basic dashboard
- [x] Simple Git repository connection
- [x] Basic Docker container deployment
- [x] Manual SSL certificate upload
- [x] Single-user project management

### Phase 2: Core Automation 🚧

- [ ] Automatic deployments from Git pushes
- [ ] Let's Encrypt SSL automation
- [ ] Basic database provisioning (PostgreSQL)
- [ ] Real-time deployment logs
- [ ] Simple monitoring dashboard

### Phase 3: Infrastructure Management 📋

- [ ] Multi-cloud server provisioning
- [ ] Advanced database options (MySQL, Redis, MongoDB)
- [ ] Custom domain management
- [ ] Team collaboration features
- [ ] Role-based access control

### Phase 4: Advanced Features 🔮

- [ ] Auto-scaling and load balancing
- [ ] Advanced monitoring and alerting
- [ ] Backup and disaster recovery
- [ ] API access and CLI tools
- [ ] Marketplace for deployment templates

---

**DeployEase** - Making self-hosting simple, one deployment at a time. 🚀