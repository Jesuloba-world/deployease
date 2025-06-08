# DeployEase Deployment Guide

This guide covers various deployment options for DeployEase, from development to production environments.

## Table of Contents

- [Prerequisites](#prerequisites)
- [Development Deployment](#development-deployment)
- [Production Deployment](#production-deployment)
- [Docker Deployment](#docker-deployment)
- [Cloud Deployment](#cloud-deployment)
- [Environment Configuration](#environment-configuration)
- [SSL/TLS Setup](#ssltls-setup)
- [Database Setup](#database-setup)
- [Monitoring and Logging](#monitoring-and-logging)
- [Backup and Recovery](#backup-and-recovery)
- [Troubleshooting](#troubleshooting)

## Prerequisites

### System Requirements

- **CPU**: 2+ cores recommended
- **RAM**: 4GB minimum, 8GB recommended
- **Storage**: 20GB minimum, SSD recommended
- **Network**: Stable internet connection

### Software Requirements

- **Docker**: 20.10+ and Docker Compose 2.0+
- **Git**: For repository management
- **Domain**: For production deployment (optional for development)

### Supported Operating Systems

- Ubuntu 20.04 LTS or later
- CentOS 8 or later
- Debian 11 or later
- macOS 11+ (development only)
- Windows 10/11 with WSL2 (development only)

## Development Deployment

### Quick Start

1. **Clone the Repository**
   ```bash
   git clone https://github.com/your-org/deployease.git
   cd deployease
   ```

2. **Environment Setup**
   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

3. **Start Development Environment**
   ```bash
   task docker:up
   ```

4. **Access the Application**
   - Frontend: http://localhost:3000
   - Backend API: http://localhost:8080
   - API Docs: http://localhost:8080/docs

### Development Services

The development environment includes:

- **Frontend**: React development server with hot reload
- **Backend**: Go application with air for hot reload
- **PostgreSQL**: Database server
- **Dragonfly**: Redis-compatible cache
- **Nginx**: Reverse proxy and static file serving

## Production Deployment

### Option 1: Docker Compose (Recommended)

1. **Server Preparation**
   ```bash
   # Update system
   sudo apt update && sudo apt upgrade -y
   
   # Install Docker
   curl -fsSL https://get.docker.com -o get-docker.sh
   sudo sh get-docker.sh
   sudo usermod -aG docker $USER
   
   # Install Docker Compose
   sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
   sudo chmod +x /usr/local/bin/docker-compose
   ```

2. **Application Setup**
   ```bash
   # Clone repository
   git clone https://github.com/your-org/deployease.git
   cd deployease
   
   # Create production environment file
   cp .env.example .env.production
   # Edit .env.production with production values
   ```

3. **Production Configuration**
   ```bash
   # Create production docker-compose file
   cp docker/docker-compose.yml docker-compose.prod.yml
   # Edit docker-compose.prod.yml for production
   ```

4. **Deploy**
   ```bash
   # Build and start services
   docker-compose -f docker-compose.prod.yml up -d
   
   # Run database migrations
   docker-compose -f docker-compose.prod.yml exec backend ./deployease migrate
   ```

### Option 2: Manual Installation

1. **Install Dependencies**
   ```bash
   # Install Go 1.24+
   wget https://go.dev/dl/go1.24.linux-amd64.tar.gz
   sudo tar -C /usr/local -xzf go1.24.linux-amd64.tar.gz
   echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
   
   # Install Node.js 18+
   curl -fsSL https://deb.nodesource.com/setup_18.x | sudo -E bash -
   sudo apt-get install -y nodejs
   
   # Install PostgreSQL
   sudo apt-get install -y postgresql postgresql-contrib
   ```

2. **Build Application**
   ```bash
   # Build backend
   cd backend
   go build -o deployease ./cmd/server
   
   # Build frontend
   cd ../frontend
   npm install
   npm run build
   ```

3. **Configure Services**
   ```bash
   # Create systemd service for backend
   sudo tee /etc/systemd/system/deployease.service > /dev/null <<EOF
   [Unit]
   Description=DeployEase Backend
   After=network.target postgresql.service
   
   [Service]
   Type=simple
   User=deployease
   WorkingDirectory=/opt/deployease
   ExecStart=/opt/deployease/deployease
   Restart=always
   RestartSec=5
   Environment=ENV=production
   EnvironmentFile=/opt/deployease/.env
   
   [Install]
   WantedBy=multi-user.target
   EOF
   
   sudo systemctl enable deployease
   sudo systemctl start deployease
   ```

## Docker Deployment

### Production Docker Compose

```yaml
# docker-compose.prod.yml
version: '3.8'

services:
  postgres:
    image: postgres:16-alpine
    environment:
      POSTGRES_DB: deployease_prod
      POSTGRES_USER: deployease
      POSTGRES_PASSWORD: ${POSTGRES_PASSWORD}
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./backups:/backups
    restart: unless-stopped
    networks:
      - deployease

  dragonfly:
    image: docker.dragonflydb.io/dragonflydb/dragonfly:v1.13.0
    command: dragonfly --requirepass=${DRAGONFLY_PASSWORD}
    volumes:
      - dragonfly_data:/data
    restart: unless-stopped
    networks:
      - deployease

  backend:
    build:
      context: .
      dockerfile: docker/Dockerfile.backend
    environment:
      - ENV=production
      - DATABASE_URL=postgres://deployease:${POSTGRES_PASSWORD}@postgres:5432/deployease_prod?sslmode=disable
      - DRAGONFLY_URL=dragonfly://default:${DRAGONFLY_PASSWORD}@dragonfly:6379
      - JWT_SECRET=${JWT_SECRET}
    depends_on:
      - postgres
      - dragonfly
    restart: unless-stopped
    networks:
      - deployease
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock
      - ./data:/app/data

  frontend:
    build:
      context: .
      dockerfile: docker/Dockerfile.frontend
    restart: unless-stopped
    networks:
      - deployease

  nginx:
    image: nginx:alpine
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./nginx/nginx.conf:/etc/nginx/nginx.conf
      - ./nginx/ssl:/etc/nginx/ssl
      - certbot_data:/var/www/certbot
    depends_on:
      - backend
      - frontend
    restart: unless-stopped
    networks:
      - deployease

  certbot:
    image: certbot/certbot
    volumes:
      - certbot_data:/var/www/certbot
      - ./nginx/ssl:/etc/letsencrypt
    command: certonly --webroot --webroot-path=/var/www/certbot --email ${ADMIN_EMAIL} --agree-tos --no-eff-email -d ${DOMAIN}

volumes:
  postgres_data:
  dragonfly_data:
  certbot_data:

networks:
  deployease:
    driver: bridge
```

### Docker Build Optimization

```dockerfile
# Multi-stage build for smaller images
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY backend/go.mod backend/go.sum ./
RUN go mod download
COPY backend/ .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o deployease ./cmd/server

FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/
COPY --from=builder /app/deployease .
EXPOSE 8080
CMD ["./deployease"]
```

## Cloud Deployment

### AWS Deployment

#### Using AWS ECS

1. **Create ECS Cluster**
   ```bash
   aws ecs create-cluster --cluster-name deployease-prod
   ```

2. **Create Task Definition**
   ```json
   {
     "family": "deployease",
     "networkMode": "awsvpc",
     "requiresCompatibilities": ["FARGATE"],
     "cpu": "1024",
     "memory": "2048",
     "executionRoleArn": "arn:aws:iam::account:role/ecsTaskExecutionRole",
     "containerDefinitions": [
       {
         "name": "deployease-backend",
         "image": "your-registry/deployease-backend:latest",
         "portMappings": [
           {
             "containerPort": 8080,
             "protocol": "tcp"
           }
         ],
         "environment": [
           {
             "name": "ENV",
             "value": "production"
           }
         ],
         "logConfiguration": {
           "logDriver": "awslogs",
           "options": {
             "awslogs-group": "/ecs/deployease",
             "awslogs-region": "us-east-1",
             "awslogs-stream-prefix": "ecs"
           }
         }
       }
     ]
   }
   ```

#### Using AWS App Runner

```yaml
# apprunner.yaml
version: 1.0
runtime: docker
build:
  commands:
    build:
      - echo Build started on `date`
      - docker build -t deployease-backend -f docker/Dockerfile.backend .
run:
  runtime-version: latest
  command: ./deployease
  network:
    port: 8080
    env: PORT
  env:
    - name: ENV
      value: production
```

### DigitalOcean Deployment

#### Using DigitalOcean App Platform

```yaml
# .do/app.yaml
name: deployease
services:
- name: backend
  source_dir: /
  github:
    repo: your-org/deployease
    branch: main
  run_command: ./deployease
  environment_slug: go
  instance_count: 1
  instance_size_slug: basic-xxs
  dockerfile_path: docker/Dockerfile.backend
  http_port: 8080
  env:
  - key: ENV
    value: production
  - key: DATABASE_URL
    scope: RUN_AND_BUILD_TIME
    type: SECRET

- name: frontend
  source_dir: /
  github:
    repo: your-org/deployease
    branch: main
  run_command: serve -s build
  environment_slug: node-js
  instance_count: 1
  instance_size_slug: basic-xxs
  dockerfile_path: docker/Dockerfile.frontend
  http_port: 3000

databases:
- name: postgres
  engine: PG
  version: "16"
  size: db-s-1vcpu-1gb
```

### Google Cloud Platform

#### Using Cloud Run

```bash
# Build and push to Container Registry
gcloud builds submit --tag gcr.io/PROJECT_ID/deployease-backend

# Deploy to Cloud Run
gcloud run deploy deployease-backend \
  --image gcr.io/PROJECT_ID/deployease-backend \
  --platform managed \
  --region us-central1 \
  --allow-unauthenticated \
  --set-env-vars ENV=production
```

## Environment Configuration

### Production Environment Variables

```bash
# .env.production

# Application
ENV=production
SERVER_PORT=8080
SERVER_HOST=0.0.0.0
DOMAIN=your-domain.com

# Database
DATABASE_URL=postgres://user:password@host:5432/deployease_prod?sslmode=require
DATABASE_MAX_CONNECTIONS=25
DATABASE_MAX_IDLE_CONNECTIONS=5

# Cache
DRAGONFLY_URL=dragonfly://password@host:6379
DRAGONFLY_MAX_CONNECTIONS=10

# Authentication
JWT_SECRET=your-super-secret-jwt-key-min-32-chars
JWT_EXPIRY=1h
REFRESH_TOKEN_EXPIRY=168h

# Session
SESSION_SECRET=your-session-secret-key
SESSION_SECURE=true
SESSION_DOMAIN=.your-domain.com

# SSL/TLS
SSL_ENABLED=true
SSL_CERT_PATH=/etc/ssl/certs/deployease.crt
SSL_KEY_PATH=/etc/ssl/private/deployease.key

# Git Providers
GITHUB_CLIENT_ID=your-github-client-id
GITHUB_CLIENT_SECRET=your-github-client-secret
GITLAB_CLIENT_ID=your-gitlab-client-id
GITLAB_CLIENT_SECRET=your-gitlab-client-secret

# Cloud Providers
AWS_ACCESS_KEY_ID=your-aws-access-key
AWS_SECRET_ACCESS_KEY=your-aws-secret-key
AWS_REGION=us-east-1

DO_TOKEN=your-digitalocean-token
DO_REGION=nyc3

# Docker
DOCKER_REGISTRY=your-registry.com
DOCKER_USERNAME=your-username
DOCKER_PASSWORD=your-password

# Monitoring
SENTRY_DSN=your-sentry-dsn
PROMETHEUS_ENABLED=true
PROMETHEUS_PORT=9090

# Logging
LOG_LEVEL=info
LOG_FORMAT=json
LOG_FILE=/var/log/deployease/app.log

# Email
SMTP_HOST=smtp.your-provider.com
SMTP_PORT=587
SMTP_USERNAME=your-smtp-username
SMTP_PASSWORD=your-smtp-password
SMTP_FROM=noreply@your-domain.com

# Notifications
SLACK_WEBHOOK_URL=your-slack-webhook-url
DISCORD_WEBHOOK_URL=your-discord-webhook-url

# Frontend
VITE_API_URL=https://api.your-domain.com
VITE_WS_URL=wss://api.your-domain.com
VITE_SENTRY_DSN=your-frontend-sentry-dsn
```

### Security Configuration

```bash
# Generate secure secrets
openssl rand -base64 32  # For JWT_SECRET
openssl rand -base64 32  # For SESSION_SECRET

# Set proper file permissions
chmod 600 .env.production
chown deployease:deployease .env.production
```

## SSL/TLS Setup

### Let's Encrypt with Certbot

1. **Install Certbot**
   ```bash
   sudo apt-get install certbot python3-certbot-nginx
   ```

2. **Obtain Certificate**
   ```bash
   sudo certbot --nginx -d your-domain.com -d www.your-domain.com
   ```

3. **Auto-renewal**
   ```bash
   # Add to crontab
   0 12 * * * /usr/bin/certbot renew --quiet
   ```

### Manual SSL Certificate

```nginx
# nginx/ssl.conf
server {
    listen 443 ssl http2;
    server_name your-domain.com;
    
    ssl_certificate /etc/ssl/certs/your-domain.com.crt;
    ssl_certificate_key /etc/ssl/private/your-domain.com.key;
    
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers ECDHE-RSA-AES256-GCM-SHA512:DHE-RSA-AES256-GCM-SHA512:ECDHE-RSA-AES256-GCM-SHA384:DHE-RSA-AES256-GCM-SHA384;
    ssl_prefer_server_ciphers off;
    
    location / {
        proxy_pass http://backend:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

## Database Setup

### PostgreSQL Production Configuration

```sql
-- Create production database and user
CREATE DATABASE deployease_prod;
CREATE USER deployease WITH ENCRYPTED PASSWORD 'secure_password';
GRANT ALL PRIVILEGES ON DATABASE deployease_prod TO deployease;

-- Performance tuning
ALTER SYSTEM SET shared_buffers = '256MB';
ALTER SYSTEM SET effective_cache_size = '1GB';
ALTER SYSTEM SET maintenance_work_mem = '64MB';
ALTER SYSTEM SET checkpoint_completion_target = 0.9;
ALTER SYSTEM SET wal_buffers = '16MB';
ALTER SYSTEM SET default_statistics_target = 100;
ALTER SYSTEM SET random_page_cost = 1.1;
ALTER SYSTEM SET effective_io_concurrency = 200;

SELECT pg_reload_conf();
```

### Database Migrations

```bash
# Run migrations in production
docker-compose exec backend ./deployease migrate

# Or manually
goose -dir migrations postgres "$DATABASE_URL" up
```

### Database Backup

```bash
#!/bin/bash
# backup.sh

DATE=$(date +%Y%m%d_%H%M%S)
BACKUP_DIR="/backups"
DB_NAME="deployease_prod"

# Create backup
pg_dump -h postgres -U deployease -d $DB_NAME > $BACKUP_DIR/backup_$DATE.sql

# Compress backup
gzip $BACKUP_DIR/backup_$DATE.sql

# Remove backups older than 7 days
find $BACKUP_DIR -name "backup_*.sql.gz" -mtime +7 -delete

echo "Backup completed: backup_$DATE.sql.gz"
```

## Monitoring and Logging

### Prometheus Configuration

```yaml
# prometheus.yml
global:
  scrape_interval: 15s

scrape_configs:
  - job_name: 'deployease'
    static_configs:
      - targets: ['backend:8080']
    metrics_path: '/metrics'
    scrape_interval: 5s

  - job_name: 'postgres'
    static_configs:
      - targets: ['postgres-exporter:9187']

  - job_name: 'nginx'
    static_configs:
      - targets: ['nginx-exporter:9113']
```

### Grafana Dashboard

```json
{
  "dashboard": {
    "title": "DeployEase Monitoring",
    "panels": [
      {
        "title": "HTTP Requests",
        "type": "graph",
        "targets": [
          {
            "expr": "rate(http_requests_total[5m])",
            "legendFormat": "{{method}} {{status}}"
          }
        ]
      },
      {
        "title": "Response Time",
        "type": "graph",
        "targets": [
          {
            "expr": "histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m]))",
            "legendFormat": "95th percentile"
          }
        ]
      }
    ]
  }
}
```

### Log Aggregation

```yaml
# docker-compose.logging.yml
version: '3.8'

services:
  elasticsearch:
    image: docker.elastic.co/elasticsearch/elasticsearch:8.5.0
    environment:
      - discovery.type=single-node
      - "ES_JAVA_OPTS=-Xms512m -Xmx512m"
    volumes:
      - elasticsearch_data:/usr/share/elasticsearch/data

  kibana:
    image: docker.elastic.co/kibana/kibana:8.5.0
    ports:
      - "5601:5601"
    environment:
      - ELASTICSEARCH_HOSTS=http://elasticsearch:9200
    depends_on:
      - elasticsearch

  logstash:
    image: docker.elastic.co/logstash/logstash:8.5.0
    volumes:
      - ./logstash/pipeline:/usr/share/logstash/pipeline
    depends_on:
      - elasticsearch

volumes:
  elasticsearch_data:
```

## Backup and Recovery

### Automated Backup Script

```bash
#!/bin/bash
# full-backup.sh

set -e

BACKUP_DIR="/backups/$(date +%Y%m%d)"
S3_BUCKET="your-backup-bucket"

# Create backup directory
mkdir -p $BACKUP_DIR

# Database backup
echo "Backing up database..."
docker-compose exec -T postgres pg_dump -U deployease deployease_prod > $BACKUP_DIR/database.sql

# Application data backup
echo "Backing up application data..."
tar -czf $BACKUP_DIR/app-data.tar.gz ./data

# Configuration backup
echo "Backing up configuration..."
cp .env.production $BACKUP_DIR/
cp docker-compose.prod.yml $BACKUP_DIR/

# Upload to S3
echo "Uploading to S3..."
aws s3 sync $BACKUP_DIR s3://$S3_BUCKET/deployease/$(date +%Y%m%d)/

# Cleanup local backups older than 3 days
find /backups -type d -mtime +3 -exec rm -rf {} +

echo "Backup completed successfully"
```

### Recovery Procedure

```bash
#!/bin/bash
# restore.sh

BACKUP_DATE=$1
S3_BUCKET="your-backup-bucket"

if [ -z "$BACKUP_DATE" ]; then
    echo "Usage: $0 YYYYMMDD"
    exit 1
fi

# Download backup from S3
aws s3 sync s3://$S3_BUCKET/deployease/$BACKUP_DATE/ ./restore/

# Stop services
docker-compose down

# Restore database
echo "Restoring database..."
docker-compose up -d postgres
sleep 10
docker-compose exec -T postgres psql -U deployease -c "DROP DATABASE IF EXISTS deployease_prod;"
docker-compose exec -T postgres psql -U deployease -c "CREATE DATABASE deployease_prod;"
docker-compose exec -T postgres psql -U deployease deployease_prod < ./restore/database.sql

# Restore application data
echo "Restoring application data..."
rm -rf ./data
tar -xzf ./restore/app-data.tar.gz

# Restore configuration
cp ./restore/.env.production .
cp ./restore/docker-compose.prod.yml .

# Start services
docker-compose up -d

echo "Restore completed successfully"
```

## Troubleshooting

### Common Issues

#### 1. Database Connection Issues

```bash
# Check database connectivity
docker-compose exec backend pg_isready -h postgres -p 5432

# Check database logs
docker-compose logs postgres

# Test connection manually
docker-compose exec postgres psql -U deployease -d deployease_prod -c "SELECT 1;"
```

#### 2. Memory Issues

```bash
# Check memory usage
docker stats

# Increase memory limits in docker-compose.yml
services:
  backend:
    deploy:
      resources:
        limits:
          memory: 1G
        reservations:
          memory: 512M
```

#### 3. SSL Certificate Issues

```bash
# Check certificate validity
openssl x509 -in /etc/ssl/certs/your-domain.com.crt -text -noout

# Test SSL configuration
openssl s_client -connect your-domain.com:443 -servername your-domain.com

# Renew Let's Encrypt certificate
sudo certbot renew --force-renewal
```

#### 4. Performance Issues

```bash
# Check application metrics
curl http://localhost:8080/metrics

# Monitor database performance
docker-compose exec postgres psql -U deployease -d deployease_prod -c "
  SELECT query, calls, total_time, mean_time 
  FROM pg_stat_statements 
  ORDER BY total_time DESC 
  LIMIT 10;
"

# Check disk usage
df -h
docker system df
```

### Log Analysis

```bash
# View application logs
docker-compose logs -f backend

# Search for errors
docker-compose logs backend | grep ERROR

# Monitor real-time logs
tail -f /var/log/deployease/app.log | jq .
```

### Health Checks

```bash
#!/bin/bash
# health-check.sh

# Check backend health
if curl -f http://localhost:8080/health > /dev/null 2>&1; then
    echo "✓ Backend is healthy"
else
    echo "✗ Backend is unhealthy"
fi

# Check database health
if docker-compose exec -T postgres pg_isready > /dev/null 2>&1; then
    echo "✓ Database is healthy"
else
    echo "✗ Database is unhealthy"
fi

# Check disk space
DISK_USAGE=$(df / | awk 'NR==2 {print $5}' | sed 's/%//')
if [ $DISK_USAGE -lt 80 ]; then
    echo "✓ Disk usage is normal ($DISK_USAGE%)"
else
    echo "⚠ Disk usage is high ($DISK_USAGE%)"
fi
```

## Maintenance

### Regular Maintenance Tasks

```bash
#!/bin/bash
# maintenance.sh

# Update Docker images
docker-compose pull
docker-compose up -d

# Clean up unused Docker resources
docker system prune -f

# Vacuum database
docker-compose exec postgres psql -U deployease -d deployease_prod -c "VACUUM ANALYZE;"

# Rotate logs
logrotate /etc/logrotate.d/deployease

# Check for security updates
sudo apt update && sudo apt list --upgradable
```

### Scaling

```yaml
# docker-compose.scale.yml
version: '3.8'

services:
  backend:
    deploy:
      replicas: 3
    
  nginx:
    volumes:
      - ./nginx/nginx-lb.conf:/etc/nginx/nginx.conf
```

```nginx
# nginx-lb.conf
upstream backend {
    server backend_1:8080;
    server backend_2:8080;
    server backend_3:8080;
}

server {
    listen 80;
    location / {
        proxy_pass http://backend;
    }
}
```

This deployment guide provides comprehensive instructions for deploying DeployEase in various environments. Choose the deployment method that best fits your infrastructure and requirements.