# Contributing to DeployEase

Thank you for your interest in contributing to DeployEase! This document provides guidelines and information for contributors.

## 🚀 Getting Started

### Prerequisites

- Go 1.24 or later
- Node.js 18.0 or later
- PostgreSQL 16 or later
- Docker and Docker Compose
- Task runner ([Installation guide](https://taskfile.dev/installation/))

### Development Setup

1. **Fork and Clone**
   ```bash
   git clone https://github.com/your-username/deployease.git
   cd deployease
   ```

2. **Environment Setup**
   ```bash
   cp .env.example .env
   # Edit .env with your local configuration
   ```

3. **Install Dependencies**
   ```bash
   task install
   ```

4. **Start Development Environment**
   ```bash
   task dev
   ```

## 📋 Development Guidelines

### Code Style

#### Backend (Go)
- Follow standard Go conventions and idioms
- Use `gofmt` for formatting
- Run `golangci-lint` for linting
- Write meaningful variable and function names
- Add comments for exported functions and complex logic
- Use dependency injection for better testability

#### Frontend (TypeScript/React)
- Use TypeScript for all new code
- Follow ESLint and Prettier configurations
- Use functional components with hooks
- Implement proper error boundaries
- Write accessible components (ARIA labels, semantic HTML)
- Use shadcn/ui components when possible

#### Database
- Use Goose for migrations
- Write reversible migrations when possible
- Use sqlc for type-safe database operations
- Follow PostgreSQL naming conventions (snake_case)
- Add proper indexes for performance

### Commit Messages

We use [Conventional Commits](https://www.conventionalcommits.org/) format:

```
type(scope): description

[optional body]

[optional footer]
```

**Types:**
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting, etc.)
- `refactor`: Code refactoring
- `test`: Adding or updating tests
- `chore`: Maintenance tasks

**Examples:**
```
feat(auth): add JWT refresh token support
fix(deploy): resolve container startup race condition
docs(api): update deployment endpoint documentation
```

### Branch Naming

- `feature/description` - New features
- `fix/description` - Bug fixes
- `docs/description` - Documentation updates
- `refactor/description` - Code refactoring

## 🧪 Testing

### Backend Testing

```bash
# Run all tests
cd backend && task test

# Run with coverage
cd backend && task test:coverage

# Run specific package
cd backend && go test ./internal/auth -v
```

**Testing Guidelines:**
- Write unit tests for all business logic
- Use table-driven tests for multiple scenarios
- Mock external dependencies
- Aim for >80% code coverage
- Write integration tests for API endpoints

**Test Structure:**
```go
func TestUserService_CreateUser(t *testing.T) {
    tests := []struct {
        name    string
        input   CreateUserRequest
        want    *User
        wantErr bool
    }{
        {
            name: "valid user creation",
            input: CreateUserRequest{
                Email:    "test@example.com",
                Password: "password123",
            },
            want: &User{
                Email: "test@example.com",
            },
            wantErr: false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation
        })
    }
}
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

**Testing Guidelines:**
- Write unit tests for utility functions
- Use React Testing Library for component tests
- Test user interactions, not implementation details
- Mock API calls and external dependencies
- Write accessibility tests

**Component Test Example:**
```typescript
import { render, screen, fireEvent } from '@testing-library/react'
import { LoginForm } from './LoginForm'

test('submits form with valid credentials', async () => {
  const mockOnSubmit = jest.fn()
  render(<LoginForm onSubmit={mockOnSubmit} />)
  
  fireEvent.change(screen.getByLabelText(/email/i), {
    target: { value: 'test@example.com' }
  })
  
  fireEvent.change(screen.getByLabelText(/password/i), {
    target: { value: 'password123' }
  })
  
  fireEvent.click(screen.getByRole('button', { name: /sign in/i }))
  
  expect(mockOnSubmit).toHaveBeenCalledWith({
    email: 'test@example.com',
    password: 'password123'
  })
})
```

## 🏗️ Architecture Guidelines

### Backend Architecture

We follow Clean Architecture principles:

```
backend/
├── cmd/                    # Application entry points
├── internal/
│   ├── domain/            # Business entities and interfaces
│   ├── usecase/           # Business logic
│   ├── repository/        # Data access layer
│   ├── handler/           # HTTP handlers
│   ├── middleware/        # HTTP middleware
│   └── config/            # Configuration
├── pkg/                   # Public packages
├── migrations/            # Database migrations
└── queries/               # SQL queries for sqlc
```

**Key Principles:**
- Dependencies point inward (domain has no external dependencies)
- Use interfaces for dependency injection
- Separate business logic from infrastructure concerns
- Keep handlers thin (delegate to use cases)

### Frontend Architecture

```
frontend/src/
├── components/            # Reusable UI components
├── pages/                 # Page components (Inertia.js)
├── hooks/                 # Custom React hooks
├── stores/                # Zustand stores
├── utils/                 # Utility functions
├── types/                 # TypeScript type definitions
└── styles/                # Global styles
```

**Key Principles:**
- Use composition over inheritance
- Keep components small and focused
- Use custom hooks for stateful logic
- Implement proper error boundaries
- Follow accessibility best practices

## 🔄 Pull Request Process

1. **Create Feature Branch**
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. **Make Changes**
   - Write code following our guidelines
   - Add tests for new functionality
   - Update documentation if needed

3. **Test Your Changes**
   ```bash
   task test
   task lint
   ```

4. **Commit Changes**
   ```bash
   git add .
   git commit -m "feat(scope): description of changes"
   ```

5. **Push and Create PR**
   ```bash
   git push origin feature/your-feature-name
   ```
   Then create a pull request on GitHub.

### PR Requirements

- [ ] All tests pass
- [ ] Code follows style guidelines
- [ ] New features have tests
- [ ] Documentation is updated
- [ ] Commit messages follow conventional format
- [ ] PR description explains the changes
- [ ] Breaking changes are documented

### PR Template

```markdown
## Description
Brief description of changes

## Type of Change
- [ ] Bug fix
- [ ] New feature
- [ ] Breaking change
- [ ] Documentation update

## Testing
- [ ] Unit tests added/updated
- [ ] Integration tests added/updated
- [ ] Manual testing completed

## Checklist
- [ ] Code follows style guidelines
- [ ] Self-review completed
- [ ] Documentation updated
- [ ] Tests pass locally
```

## 🐛 Bug Reports

When reporting bugs, please include:

1. **Environment Information**
   - OS and version
   - Go version
   - Node.js version
   - Browser (for frontend issues)

2. **Steps to Reproduce**
   - Clear, numbered steps
   - Expected vs actual behavior
   - Screenshots if applicable

3. **Additional Context**
   - Error messages
   - Log output
   - Configuration details

## 💡 Feature Requests

For feature requests, please provide:

1. **Problem Description**
   - What problem does this solve?
   - Who would benefit from this feature?

2. **Proposed Solution**
   - Detailed description of the feature
   - How should it work?
   - Any alternatives considered?

3. **Additional Context**
   - Mockups or examples
   - Related issues or discussions

## 📚 Documentation

### Writing Documentation

- Use clear, concise language
- Include code examples
- Add screenshots for UI features
- Keep documentation up to date with code changes
- Use proper markdown formatting

### Documentation Structure

```
docs/
├── api/                   # API documentation
├── deployment/            # Deployment guides
├── development/           # Development setup
├── user-guide/            # User documentation
└── architecture/          # Technical architecture
```

## 🤝 Community Guidelines

### Code of Conduct

- Be respectful and inclusive
- Welcome newcomers and help them learn
- Focus on constructive feedback
- Respect different opinions and approaches
- Follow the [Contributor Covenant](https://www.contributor-covenant.org/)

### Getting Help

- **GitHub Discussions**: For questions and general discussion
- **GitHub Issues**: For bug reports and feature requests
- **Discord**: For real-time chat (link in README)

### Recognition

We recognize contributors in several ways:

- Contributors list in README
- Release notes mention significant contributions
- Special recognition for first-time contributors
- Maintainer status for consistent contributors

## 🏷️ Release Process

### Versioning

We follow [Semantic Versioning](https://semver.org/):

- `MAJOR.MINOR.PATCH`
- Major: Breaking changes
- Minor: New features (backward compatible)
- Patch: Bug fixes (backward compatible)

### Release Checklist

- [ ] All tests pass
- [ ] Documentation updated
- [ ] CHANGELOG.md updated
- [ ] Version bumped in relevant files
- [ ] Release notes prepared
- [ ] Docker images built and tested

## 📞 Contact

For questions about contributing:

- Open a GitHub Discussion
- Create an issue with the "question" label
- Reach out to maintainers directly

Thank you for contributing to DeployEase! 🚀