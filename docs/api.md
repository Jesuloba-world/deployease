# DeployEase API Documentation

This document provides comprehensive information about the DeployEase REST API.

## Base URL

```
Development: http://localhost:8080
Production: https://your-domain.com
```

## Authentication

DeployEase uses JWT (JSON Web Tokens) for authentication. Include the token in the Authorization header:

```http
Authorization: Bearer <your-jwt-token>
```

### Authentication Endpoints

#### POST /auth/login

Authenticate a user and receive JWT tokens.

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "your-password"
}
```

**Response:**
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIs...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
  "expires_in": 3600,
  "user": {
    "id": "123e4567-e89b-12d3-a456-426614174000",
    "email": "user@example.com",
    "name": "John Doe",
    "created_at": "2024-01-01T00:00:00Z"
  }
}
```

#### POST /auth/refresh

Refresh an expired access token using a refresh token.

**Request Body:**
```json
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIs..."
}
```

**Response:**
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIs...",
  "expires_in": 3600
}
```

#### POST /auth/logout

Invalidate the current session and tokens.

**Headers:**
```http
Authorization: Bearer <access-token>
```

**Response:**
```json
{
  "message": "Successfully logged out"
}
```

## User Management

#### GET /users/profile

Get the current user's profile information.

**Headers:**
```http
Authorization: Bearer <access-token>
```

**Response:**
```json
{
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "email": "user@example.com",
  "name": "John Doe",
  "avatar_url": "https://example.com/avatar.jpg",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

#### PUT /users/profile

Update the current user's profile.

**Request Body:**
```json
{
  "name": "Jane Doe",
  "avatar_url": "https://example.com/new-avatar.jpg"
}
```

## Project Management

#### GET /projects

List all projects for the authenticated user.

**Query Parameters:**
- `page` (optional): Page number (default: 1)
- `limit` (optional): Items per page (default: 20, max: 100)
- `status` (optional): Filter by status (active, inactive, deploying)

**Response:**
```json
{
  "projects": [
    {
      "id": "proj_123456",
      "name": "my-awesome-app",
      "description": "A sample application",
      "status": "active",
      "git_repository": "https://github.com/user/repo.git",
      "domain": "my-app.deployease.com",
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 1,
    "total_pages": 1
  }
}
```

#### POST /projects

Create a new project.

**Request Body:**
```json
{
  "name": "my-new-app",
  "description": "Description of my new app",
  "git_repository": "https://github.com/user/new-repo.git",
  "branch": "main",
  "build_command": "npm run build",
  "start_command": "npm start",
  "environment_variables": {
    "NODE_ENV": "production",
    "API_URL": "https://api.example.com"
  }
}
```

#### GET /projects/{project_id}

Get detailed information about a specific project.

**Response:**
```json
{
  "id": "proj_123456",
  "name": "my-awesome-app",
  "description": "A sample application",
  "status": "active",
  "git_repository": "https://github.com/user/repo.git",
  "branch": "main",
  "domain": "my-app.deployease.com",
  "build_command": "npm run build",
  "start_command": "npm start",
  "environment_variables": {
    "NODE_ENV": "production"
  },
  "ssl_enabled": true,
  "auto_deploy": true,
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

## Deployment Management

#### GET /projects/{project_id}/deployments

List deployments for a project.

**Query Parameters:**
- `page` (optional): Page number
- `limit` (optional): Items per page
- `status` (optional): Filter by status (pending, building, deploying, success, failed)

**Response:**
```json
{
  "deployments": [
    {
      "id": "deploy_789012",
      "project_id": "proj_123456",
      "commit_sha": "abc123def456",
      "commit_message": "Fix user authentication bug",
      "branch": "main",
      "status": "success",
      "started_at": "2024-01-01T10:00:00Z",
      "completed_at": "2024-01-01T10:05:30Z",
      "duration": 330,
      "logs_url": "/projects/proj_123456/deployments/deploy_789012/logs"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 1,
    "total_pages": 1
  }
}
```

#### POST /projects/{project_id}/deployments

Trigger a new deployment.

**Request Body:**
```json
{
  "branch": "main",
  "commit_sha": "abc123def456" // optional, uses latest if not provided
}
```

#### GET /projects/{project_id}/deployments/{deployment_id}/logs

Get real-time deployment logs via WebSocket or HTTP.

**WebSocket URL:**
```
ws://localhost:8080/projects/{project_id}/deployments/{deployment_id}/logs
```

**HTTP Response:**
```json
{
  "logs": [
    {
      "timestamp": "2024-01-01T10:00:00Z",
      "level": "info",
      "message": "Starting deployment..."
    },
    {
      "timestamp": "2024-01-01T10:00:30Z",
      "level": "info",
      "message": "Building Docker image..."
    }
  ]
}
```

## Database Management

#### GET /projects/{project_id}/databases

List databases for a project.

**Response:**
```json
{
  "databases": [
    {
      "id": "db_345678",
      "name": "my-app-postgres",
      "type": "postgresql",
      "version": "16",
      "status": "running",
      "connection_string": "postgresql://user:pass@host:5432/dbname",
      "created_at": "2024-01-01T00:00:00Z"
    }
  ]
}
```

#### POST /projects/{project_id}/databases

Create a new database.

**Request Body:**
```json
{
  "name": "my-app-redis",
  "type": "redis",
  "version": "7",
  "size": "small"
}
```

## Domain Management

#### GET /projects/{project_id}/domains

List custom domains for a project.

**Response:**
```json
{
  "domains": [
    {
      "id": "domain_901234",
      "domain": "myapp.com",
      "status": "active",
      "ssl_status": "active",
      "certificate_expires_at": "2024-12-01T00:00:00Z",
      "created_at": "2024-01-01T00:00:00Z"
    }
  ]
}
```

#### POST /projects/{project_id}/domains

Add a custom domain.

**Request Body:**
```json
{
  "domain": "myapp.com",
  "auto_ssl": true
}
```

## Monitoring

#### GET /projects/{project_id}/metrics

Get project metrics and monitoring data.

**Query Parameters:**
- `from` (required): Start time (ISO 8601)
- `to` (required): End time (ISO 8601)
- `metric` (optional): Specific metric (cpu, memory, requests, response_time)

**Response:**
```json
{
  "metrics": {
    "cpu": {
      "average": 45.2,
      "max": 78.5,
      "data_points": [
        {
          "timestamp": "2024-01-01T10:00:00Z",
          "value": 42.1
        }
      ]
    },
    "memory": {
      "average": 512.8,
      "max": 1024.0,
      "data_points": [
        {
          "timestamp": "2024-01-01T10:00:00Z",
          "value": 498.2
        }
      ]
    }
  }
}
```

## WebSocket Events

DeployEase provides real-time updates via WebSocket connections.

### Connection

```javascript
const ws = new WebSocket('ws://localhost:8080/ws?token=your-jwt-token')
```

### Event Types

#### Deployment Events

```json
{
  "type": "deployment.started",
  "data": {
    "deployment_id": "deploy_789012",
    "project_id": "proj_123456",
    "timestamp": "2024-01-01T10:00:00Z"
  }
}
```

```json
{
  "type": "deployment.completed",
  "data": {
    "deployment_id": "deploy_789012",
    "project_id": "proj_123456",
    "status": "success",
    "timestamp": "2024-01-01T10:05:30Z"
  }
}
```

#### Log Events

```json
{
  "type": "log",
  "data": {
    "deployment_id": "deploy_789012",
    "timestamp": "2024-01-01T10:00:00Z",
    "level": "info",
    "message": "Building Docker image..."
  }
}
```

## Error Handling

All API endpoints return consistent error responses:

```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Invalid request data",
    "details": {
      "field": "email",
      "reason": "Invalid email format"
    }
  }
}
```

### Common Error Codes

- `VALIDATION_ERROR` (400): Invalid request data
- `UNAUTHORIZED` (401): Authentication required
- `FORBIDDEN` (403): Insufficient permissions
- `NOT_FOUND` (404): Resource not found
- `CONFLICT` (409): Resource already exists
- `RATE_LIMITED` (429): Too many requests
- `INTERNAL_ERROR` (500): Server error

## Rate Limiting

API requests are rate limited:

- **Authenticated requests**: 1000 requests per hour
- **Unauthenticated requests**: 100 requests per hour

Rate limit headers are included in responses:

```http
X-RateLimit-Limit: 1000
X-RateLimit-Remaining: 999
X-RateLimit-Reset: 1640995200
```

## SDKs and Libraries

### Official SDKs

- **JavaScript/TypeScript**: `@deployease/sdk-js`
- **Go**: `github.com/deployease/sdk-go`
- **Python**: `deployease-sdk`

### Example Usage (JavaScript)

```javascript
import { DeployEase } from '@deployease/sdk-js'

const client = new DeployEase({
  apiKey: 'your-api-key',
  baseUrl: 'https://api.deployease.com'
})

// List projects
const projects = await client.projects.list()

// Create deployment
const deployment = await client.deployments.create('proj_123456', {
  branch: 'main'
})
```

## Interactive Documentation

For interactive API exploration, visit:

- **Development**: http://localhost:8080/docs
- **Production**: https://your-domain.com/docs

The interactive documentation is automatically generated from the OpenAPI specification and allows you to test API endpoints directly from your browser.