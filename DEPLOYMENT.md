# ===========================================
# PRODUCTION DEPLOYMENT GUIDE
# ===========================================

## Prerequisites

### 1. Domain Setup
- Register domain (e.g., yourdomain.com)
- Configure DNS records:
  ```
  A     @              -> Your server IP
  A     api            -> Your server IP
  A     cdn            -> Your CDN IP
  A     admin          -> Your server IP
  CNAME www            -> yourdomain.com
  ```

### 2. Server Requirements
- Ubuntu 22.04 LTS or similar
- 4+ CPU cores
- 8GB+ RAM
- 100GB+ SSD storage
- Docker 24+
- Kubernetes 1.28+ (for K8s deployment)

### 3. External Services
- PostgreSQL (AWS RDS, Google Cloud SQL, or self-hosted)
- Redis (AWS ElastiCache, Redis Cloud, or self-hosted)
- SendGrid account with API key
- Zoom/Agora account for video
- AWS S3 or compatible storage
- Domain with SSL certificate

---

## Deployment Steps

### Step 1: Environment Configuration

```bash
# Copy and edit environment files
cp .env.example .env.production

# Generate secure secrets
export JWT_SECRET=$(openssl rand -hex 64)
export DB_PASSWORD=$(openssl rand -base64 32)
export ENCRYPTION_KEY=$(openssl rand -hex 32)

# Edit .env.production with your values
nano .env.production
```

Required variables:
- `DATABASE_HOST` - Your PostgreSQL host
- `DATABASE_PASSWORD` - Secure password
- `JWT_SECRET` - Generated secret
- `SENDGRID_API_KEY` - From SendGrid
- `ZOOM_API_KEY` - From Zoom
- `AWS_ACCESS_KEY_ID` - From AWS
- `AWS_SECRET_ACCESS_KEY` - From AWS

### Step 2: SendGrid Setup

```bash
# Set your API key
export SENDGRID_API_KEY="SG.your-key-here"

# Run setup script
chmod +x scripts/setup-sendgrid.sh
./scripts/setup-sendgrid.sh
```

This will:
- Validate your API key
- Create email templates
- Configure sender verification

### Step 3: SSL Certificates

```bash
# Install certbot
sudo apt-get update
sudo apt-get install certbot python3-certbot-nginx

# Run SSL setup
chmod +x scripts/setup-ssl.sh
./scripts/setup-ssl.sh
```

Certificates will be installed at:
- `/etc/letsencrypt/live/yourdomain.com/fullchain.pem`
- `/etc/letsencrypt/live/yourdomain.com/privkey.pem`

### Step 4: Database Setup

```bash
# Connect to PostgreSQL
psql -h your-db-host -U postgres

# Create databases
CREATE DATABASE auth_db;
CREATE DATABASE users_db;
CREATE DATABASE courses_db;
CREATE DATABASE tasks_db;
CREATE DATABASE progress_db;
CREATE DATABASE notifications_db;
CREATE DATABASE files_db;
CREATE DATABASE videos_db;

# Create user and grant permissions
CREATE USER langplatform WITH PASSWORD 'your-secure-password';
GRANT ALL PRIVILEGES ON DATABASE auth_db TO langplatform;
-- Repeat for all databases
```

### Step 5: Docker Deployment

```bash
# Build images
docker-compose -f infra/docker/docker-compose.yml build

# Start services
docker-compose -f infra/docker/docker-compose.yml up -d

# Check status
docker-compose ps

# View logs
docker-compose logs -f
```

### Step 6: Kubernetes Deployment

```bash
# Setup secrets
chmod +x scripts/setup-k8s-secrets.sh
./scripts/setup-k8s-secrets.sh

# Deploy services
chmod +x scripts/deploy-production.sh
./scripts/deploy-production.sh production

# Check deployment
kubectl get pods -n language-platform
kubectl get services -n language-platform
```

### Step 7: Configure NGINX

```bash
# Copy production config
sudo cp infra/nginx/nginx-production.conf /etc/nginx/nginx.conf

# Update domains in config
sudo nano /etc/nginx/nginx.conf

# Test configuration
sudo nginx -t

# Reload NGINX
sudo systemctl reload nginx
```

### Step 8: Deploy Frontend

```bash
# Build web frontend
cd frontend/web
npm install
npm run build

# Deploy to server
rsync -avz --delete .next/standalone/ user@server:/var/www/app/
rsync -avz --delete .next/static/ user@server:/var/www/app/.next/static/

# Start with PM2
pm2 start ecosystem.config.js
pm2 save
pm2 startup
```

### Step 9: Verify Deployment

```bash
# Health checks
curl https://api.yourdomain.com/health
curl https://api.yourdomain.com/api/auth/health
curl https://api.yourdomain.com/api/v1/users/health

# Test authentication
curl -X POST https://api.yourdomain.com/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "Test123!",
    "confirmPassword": "Test123!",
    "role": "Student"
  }'
```

---

## Post-Deployment

### Monitoring Setup

```bash
# Access Grafana
https://yourdomain.com:3000
Username: admin
Password: (from .env)

# Add dashboards
- Import dashboard ID 1860 (Node Exporter)
- Import dashboard ID 7362 (PostgreSQL)
- Import custom service dashboards
```

### Backup Configuration

```bash
# Create backup script
cat > /usr/local/bin/backup-db.sh << 'EOF'
#!/bin/bash
DATE=$(date +%Y%m%d_%H%M%S)
BACKUP_DIR="/backups"
pg_dump -h localhost -U langplatform auth_db | gzip > $BACKUP_DIR/auth_db_$DATE.sql.gz
# Repeat for all databases
aws s3 cp $BACKUP_DIR/ s3://your-backup-bucket/ --recursive
EOF

chmod +x /usr/local/bin/backup-db.sh

# Setup cron
crontab -e
# Add: 0 2 * * * /usr/local/bin/backup-db.sh
```

### Security Hardening

```bash
# Setup firewall
sudo ufw allow 22/tcp
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp
sudo ufw enable

# Fail2ban for SSH
sudo apt-get install fail2ban
sudo systemctl enable fail2ban

# Setup monitoring alerts
# Configure Prometheus alertmanager
# Setup PagerDuty/OpsGenie integration
```

---

## Scaling

### Horizontal Scaling (Kubernetes)

```yaml
# Scale deployments
kubectl scale deployment auth-service --replicas=5 -n language-platform
kubectl scale deployment user-service --replicas=5 -n language-platform

# Auto-scaling
kubectl autoscale deployment auth-service \
  --cpu-percent=70 \
  --min=3 \
  --max=10 \
  -n language-platform
```

### Database Scaling

- Use read replicas for read-heavy workloads
- Implement connection pooling (PgBouncer)
- Consider sharding for very large datasets

### Caching Strategy

- Redis for session storage
- Redis for frequently accessed data
- CDN for static assets
- Browser caching headers

---

## Troubleshooting

### Service won't start
```bash
# Check logs
kubectl logs -f deployment/service-name -n language-platform
docker-compose logs service-name

# Check resources
kubectl top pods -n language-platform
```

### Database connection issues
```bash
# Test connection
psql -h db-host -U langplatform -d auth_db

# Check credentials in secrets
kubectl get secret db-secret -n language-platform -o yaml
```

### SSL certificate issues
```bash
# Test certificate
sudo certbot certificates

# Renew manually
sudo certbot renew

# Check NGINX config
sudo nginx -t
```

---

## Rollback Procedure

```bash
# Kubernetes rollback
kubectl rollout undo deployment/service-name -n language-platform

# Docker rollback
docker-compose down
git checkout previous-commit
docker-compose up -d

# Database rollback
psql -h db-host -U langplatform -d auth_db < backup.sql
```

---

## Maintenance

### Regular Tasks

**Daily:**
- Monitor error logs
- Check service health
- Review metrics dashboard

**Weekly:**
- Review security alerts
- Update dependencies
- Test backups

**Monthly:**
- Security patches
- Performance optimization
- Capacity planning

### Update Procedure

```bash
# Pull latest code
git pull origin main

# Backup database
./scripts/backup-db.sh

# Deploy update
./scripts/deploy-production.sh production

# Verify
./scripts/smoke-test.sh
```

---

## Support

For issues:
- Check logs: `kubectl logs` or `docker logs`
- Review documentation in `docs/`
- Contact: devops@yourdomain.com
