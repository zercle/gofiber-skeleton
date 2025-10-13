# Deployment Guide

This guide covers various deployment options for the Go Fiber Skeleton application.

## 🐳 Docker Deployment

### Development Environment

```bash
# Start all services
docker-compose up -d

# View logs
docker-compose logs -f

# Stop services
docker-compose down
```

### Production Build

```bash
# Build production image
make docker-build

# Or build manually
docker build -t gofiber-skeleton:latest .
```

### Production Run

```bash
# Create network
docker network gofiber-network

# Run PostgreSQL
docker run -d \
  --name gofiber-postgres \
  --network gofiber-network \
  -e POSTGRES_DB=gofiber_skeleton \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=secure_password \
  -p 5432:5432 \
  postgres:16-alpine

# Run application
docker run -d \
  --name gofiber-skeleton \
  --network gofiber-network \
  -p 8080:8080 \
  -e GS_DATABASE_HOST=gofiber-postgres \
  -e GS_DATABASE_PASSWORD=secure_password \
  -e GS_JWT_SECRET=your-super-secret-production-key \
  gofiber-skeleton:latest
```

### Docker Compose Production

Create `docker-compose.prod.yml`:

```yaml
version: '3.8'

services:
  app:
    build: .
    ports:
      - "8080:8080"
    environment:
      - GS_DATABASE_HOST=postgres
      - GS_DATABASE_PASSWORD=secure_password
      - GS_JWT_SECRET=your-super-secret-production-key
      - APP_ENV=production
    depends_on:
      postgres:
        condition: service_healthy
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "wget", "--no-verbose", "--tries=1", "--spider", "http://localhost:8080/health"]
      interval: 30s
      timeout: 10s
      retries: 3

  postgres:
    image: postgres:16-alpine
    environment:
      POSTGRES_DB: gofiber_skeleton
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: secure_password
    volumes:
      - postgres_data:/var/lib/postgresql/data
    restart: unless-stopped
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres"]
      interval: 10s
      timeout: 5s
      retries: 5

volumes:
  postgres_data:
```

```bash
docker-compose -f docker-compose.prod.yml up -d
```

## 🚀 Kubernetes Deployment

### Namespace and ConfigMap

```yaml
# k8s/namespace.yaml
apiVersion: v1
kind: Namespace
metadata:
  name: gofiber-skeleton

---
# k8s/configmap.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: gofiber-skeleton-config
  namespace: gofiber-skeleton
data:
  config.yaml: |
    server:
      host: "0.0.0.0"
      port: 8080
      read_timeout: 15
      write_timeout: 15
      idle_timeout: 60
    logger:
      level: "info"
      format: "json"
      output: "stdout"
```

### Secret

```yaml
# k8s/secret.yaml
apiVersion: v1
kind: Secret
metadata:
  name: gofiber-skeleton-secret
  namespace: gofiber-skeleton
type: Opaque
data:
  jwt-secret: <base64-encoded-jwt-secret>
  db-password: <base64-encoded-db-password>
```

### Deployment

```yaml
# k8s/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: gofiber-skeleton
  namespace: gofiber-skeleton
spec:
  replicas: 3
  selector:
    matchLabels:
      app: gofiber-skeleton
  template:
    metadata:
      labels:
        app: gofiber-skeleton
    spec:
      containers:
      - name: gofiber-skeleton
        image: gofiber-skeleton:latest
        ports:
        - containerPort: 8080
        env:
        - name: GS_DATABASE_HOST
          value: "postgres-service"
        - name: GS_DATABASE_PASSWORD
          valueFrom:
            secretKeyRef:
              name: gofiber-skeleton-secret
              key: db-password
        - name: GS_JWT_SECRET
          valueFrom:
            secretKeyRef:
              name: gofiber-skeleton-secret
              key: jwt-secret
        - name: APP_ENV
          value: "production"
        resources:
          requests:
            memory: "64Mi"
            cpu: "50m"
          limits:
            memory: "256Mi"
            cpu: "200m"
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 10
          periodSeconds: 30
        readinessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 10
```

### Service

```yaml
# k8s/service.yaml
apiVersion: v1
kind: Service
metadata:
  name: gofiber-skeleton-service
  namespace: gofiber-skeleton
spec:
  selector:
    app: gofiber-skeleton
  ports:
  - protocol: TCP
    port: 80
    targetPort: 8080
  type: LoadBalancer
```

### Ingress

```yaml
# k8s/ingress.yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: gofiber-skeleton-ingress
  namespace: gofiber-skeleton
  annotations:
    kubernetes.io/ingress.class: nginx
    cert-manager.io/cluster-issuer: letsencrypt-prod
    nginx.ingress.kubernetes.io/rate-limit: "100"
    nginx.ingress.kubernetes.io/rate-limit-window: "1m"
spec:
  tls:
  - hosts:
    - api.yourdomain.com
    secretName: gofiber-skeleton-tls
  rules:
  - host: api.yourdomain.com
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: gofiber-skeleton-service
            port:
              number: 80
```

### Deploy to Kubernetes

```bash
# Apply all manifests
kubectl apply -f k8s/

# Check status
kubectl get pods -n gofiber-skeleton
kubectl get services -n gofiber-skeleton
kubectl get ingress -n gofiber-skeleton

# View logs
kubectl logs -f deployment/gofiber-skeleton -n gofiber-skeleton
```

## ☁️ Cloud Deployment

### Google Cloud Run

```bash
# Build and push to Google Container Registry
gcloud builds submit --tag gcr.io/PROJECT-ID/gofiber-skeleton

# Deploy to Cloud Run
gcloud run deploy gofiber-skeleton \
  --image gcr.io/PROJECT-ID/gofiber-skeleton \
  --platform managed \
  --region us-central1 \
  --allow-unauthenticated \
  --set-env-vars "GS_DATABASE_HOST=your-db-host,GS_JWT_SECRET=your-secret"
```

### AWS ECS

```json
{
  "family": "gofiber-skeleton",
  "networkMode": "awsvpc",
  "requiresCompatibilities": ["FARGATE"],
  "cpu": "256",
  "memory": "512",
  "executionRoleArn": "arn:aws:iam::account:role/ecsTaskExecutionRole",
  "taskRoleArn": "arn:aws:iam::account:role/ecsTaskRole",
  "containerDefinitions": [
    {
      "name": "gofiber-skeleton",
      "image": "your-account.dkr.ecr.region.amazonaws.com/gofiber-skeleton:latest",
      "portMappings": [
        {
          "containerPort": 8080,
          "protocol": "tcp"
        }
      ],
      "environment": [
        {
          "name": "APP_ENV",
          "value": "production"
        }
      ],
      "secrets": [
        {
          "name": "GS_JWT_SECRET",
          "valueFrom": "arn:aws:secretsmanager:region:account:secret:gofiber-skeleton/jwt-secret"
        }
      ],
      "logConfiguration": {
        "logDriver": "awslogs",
        "options": {
          "awslogs-group": "/ecs/gofiber-skeleton",
          "awslogs-region": "us-west-2",
          "awslogs-stream-prefix": "ecs"
        }
      }
    }
  ]
}
```

### Heroku

```bash
# Install Heroku CLI
# Login to Heroku
heroku login

# Create app
heroku create your-app-name

# Set buildpack
heroku buildpacks:set heroku/go

# Set environment variables
heroku config:set APP_ENV=production
heroku config:set GS_DATABASE_URL=postgres://...
heroku config:set GS_JWT_SECRET=your-super-secret-key

# Deploy
git push heroku main
```

## 📊 Monitoring & Logging

### Prometheus Monitoring

Add to main application:

```go
import "github.com/prometheus/client_golang/prometheus/promhttp"

// Add metrics endpoint
app.Get("/metrics", fiber.WrapHandler(promhttp.Handler()))
```

### Structured Logging

The application already includes structured JSON logging. For production:

```yaml
# config.yaml
logger:
  level: "info"
  format: "json"
  output: "stdout"
```

### Health Checks

Built-in health check at `/health`:

```json
{
  "status": "ok",
  "timestamp": "2024-01-01T00:00:00Z",
  "service": "gofiber-skeleton",
  "version": "1.0.0"
}
```

## 🔧 Environment Configuration

### Production Environment Variables

```bash
# Server
GS_SERVER_HOST=0.0.0.0
GS_SERVER_PORT=8080

# Database
GS_DATABASE_HOST=your-production-db-host
GS_DATABASE_PORT=5432
GS_DATABASE_USER=postgres
GS_DATABASE_PASSWORD=secure_password
GS_DATABASE_DBNAME=gofiber_skeleton
GS_DATABASE_SSLMODE=require

# JWT
GS_JWT_SECRET=your-super-secret-production-key

# Logger
GS_LOGGER_LEVEL=info
GS_LOGGER_FORMAT=json

# Application
APP_ENV=production
```

### Security Considerations

1. **Change default secrets**: Never use default JWT secrets in production
2. **Use HTTPS**: Always use TLS in production
3. **Database SSL**: Enable SSL for database connections
4. **Environment variables**: Store sensitive data in secrets, not code
5. **Rate limiting**: Configure appropriate rate limits
6. **CORS**: Configure CORS for your specific domains

## 🚀 CI/CD Pipeline

### GitHub Actions Example

```yaml
# .github/workflows/deploy.yml
name: Deploy

on:
  push:
    branches: [main]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v3
    - uses: actions/setup-go@v3
      with:
        go-version: 1.25
    - run: make test lint

  build:
    needs: test
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v3
    - uses: docker/build-push-action@v3
      with:
        push: true
        tags: your-registry/gofiber-skeleton:latest

  deploy:
    needs: build
    runs-on: ubuntu-latest
    steps:
    - name: Deploy to production
      run: |
        # Your deployment script
        kubectl apply -f k8s/
```

## 🔍 Troubleshooting

### Common Issues

1. **Database connection failed**
   - Check database credentials
   - Verify network connectivity
   - Check SSL settings

2. **JWT authentication failed**
   - Verify JWT secret matches
   - Check token expiration
   - Ensure proper header format

3. **Container won't start**
   - Check logs: `docker logs container-name`
   - Verify environment variables
   - Check port conflicts

### Debug Commands

```bash
# Docker logs
docker logs gofiber-skeleton

# Kubernetes logs
kubectl logs -f deployment/gofiber-skeleton -n gofiber-skeleton

# Health check
curl http://localhost:8080/health

# Database connection
docker exec -it gofiber-postgres psql -U postgres -d gofiber_skeleton
```

## 📈 Performance Tuning

### Go Fiber Configuration

```go
app := fiber.New(fiber.Config{
    ReadTimeout:           15 * time.Second,
    WriteTimeout:          15 * time.Second,
    IdleTimeout:           30 * time.Second,
    Prefork:               true,  // Enable prefork
    ServerHeader:          "Fiber",
    DisableKeepalive:      false,
})
```

### Database Connection Pool

```yaml
database:
  max_open_conns: 100
  max_idle_conns: 10
```

### Production Optimization

1. **Enable prefork** for better performance
2. **Configure proper timeouts**
3. **Use connection pooling**
4. **Enable Gzip compression**
5. **Configure appropriate cache headers**
6. **Use CDN for static assets**

This deployment guide covers the most common deployment scenarios. Choose the one that best fits your infrastructure and requirements.