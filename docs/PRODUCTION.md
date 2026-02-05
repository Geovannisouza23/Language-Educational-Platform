# 🚀 Production Deployment - Quick Reference

## 📋 Overview

Complete Language Learning Platform with 8 microservices, web & mobile frontends, ready for production deployment.

---

## ⚡ Quick Commands

### 1. Initial Setup (One-time)
```bash
# Clone repo
git clone <repo-url> && cd user.api

# Copy environment template
cp config/env/.env.example config/env/.env.production

# Generate secure secrets
export JWT_SECRET=$(openssl rand -hex 64)
export DB_PASSWORD=$(openssl rand -base64 32)

# Edit production config
nano config/env/.env.production
```

### 2. Configure SendGrid
```bash
export SENDGRID_API_KEY="SG.your-key-here"
./scripts/setup-sendgrid.sh
```

### 3. Setup SSL
```bash
./scripts/setup-ssl.sh
```

### 4. Deploy with Docker
```bash
docker-compose -f infra/docker/docker-compose.yml up -d
```

### 5. Deploy with Kubernetes
```bash
./scripts/setup-k8s-secrets.sh
./scripts/deploy-production.sh production
```

### 6. Verify Deployment
```bash
curl https://api.yourdomain.com/health
kubectl get pods -n language-platform
```

---

## 🔑 Required Environment Variables

### Critical (Must Configure)
```bash
# Database
DATABASE_HOST=your-db-host.rds.amazonaws.com
DATABASE_PASSWORD=<generate-secure-password>

# JWT
JWT_SECRET=<generate-with-openssl-rand-hex-64>

# SendGrid
SENDGRID_API_KEY=SG.your-sendgrid-key
SENDGRID_FROM_EMAIL=noreply@yourdomain.com

# Video (Choose ONE)
ZOOM_API_KEY=your-zoom-key              # Recommended
AGORA_APP_ID=your-agora-id              # Alternative
JITSI_DOMAIN=meet.jit.si                # Free option

# Storage
AWS_ACCESS_KEY_ID=your-aws-key
AWS_SECRET_ACCESS_KEY=your-aws-secret
AWS_S3_BUCKET=your-bucket-name

# Domains
PRIMARY_DOMAIN=yourdomain.com
API_DOMAIN=api.yourdomain.com
```

---

## 🎯 Service Ports

| Service | Port | URL |
|---------|------|-----|
| NGINX Gateway | 80/443 | https://api.yourdomain.com |
| Auth Service | 5001 | Internal |
| User Service | 8001 | Internal |
| Course Service | 8002 | Internal |
| Task Service | 8003 | Internal |
| Progress Service | 8004 | Internal |
| Notification Service | 8005 | Internal |
| File Service | 8006 | Internal |
| Video Service | 8007 | Internal |
| Prometheus | 9090 | http://localhost:9090 |
| Grafana | 3000 | http://localhost:3000 |

---

## 🌐 DNS Configuration

```
A     @              -> Your server IP
A     api            -> Your server IP
A     cdn            -> Your CDN IP (optional)
A     admin          -> Your server IP
CNAME www            -> yourdomain.com
```

---

## 📝 Video Provider Setup

### Zoom (Recommended)
1. Go to https://marketplace.zoom.us
2. Create Server-to-Server OAuth app
3. Get Account ID, Client ID, Client Secret
4. Add to `config/env/.env.production`:
   ```
   ZOOM_ACCOUNT_ID=your-account-id
   ZOOM_CLIENT_ID=your-client-id
   ZOOM_CLIENT_SECRET=your-client-secret
   VIDEO_PROVIDER=zoom
   ```

### Agora (Scalable)
1. Go to https://console.agora.io
2. Create project with "Secured mode"
3. Get App ID, App Certificate
4. Add to `config/env/.env.production`:
   ```
   AGORA_APP_ID=your-app-id
   AGORA_APP_CERTIFICATE=your-certificate
   VIDEO_PROVIDER=agora
   ```

### Jitsi (Free)
1. Use public: `JITSI_DOMAIN=meet.jit.si`
2. Or self-host: Follow docs/VIDEO_SETUP.md
3. Add to `config/env/.env.production`:
   ```
   JITSI_DOMAIN=meet.jit.si
   VIDEO_PROVIDER=jitsi
   ```

---

## 🔒 SSL Certificate Setup

```bash
# Auto-setup with Let's Encrypt
./scripts/setup-ssl.sh

# Manual setup
sudo certbot certonly --nginx \
  -d yourdomain.com \
  -d api.yourdomain.com \
  -d cdn.yourdomain.com \
  --email admin@yourdomain.com
```

Certificates installed at:
- `/etc/letsencrypt/live/yourdomain.com/fullchain.pem`
- `/etc/letsencrypt/live/yourdomain.com/privkey.pem`

---

## 📧 SendGrid Configuration

### Quick Setup
```bash
# 1. Get API key from SendGrid dashboard
export SENDGRID_API_KEY="SG.xxx"

# 2. Run setup script
./scripts/setup-sendgrid.sh

# 3. Script will create email templates and return IDs
# 4. Add IDs to config/env/.env.production
```

### Manual Template Creation
See detailed guide in docs/DEPLOYMENT.md

---

## 🐳 Docker Commands

```bash
# Start all services
docker-compose up -d

# View logs
docker-compose logs -f [service-name]

# Restart service
docker-compose restart [service-name]

# Stop all services
docker-compose down

# Rebuild service
docker-compose build [service-name]
docker-compose up -d [service-name]
```

---

## ☸️ Kubernetes Commands

```bash
# Check pods
kubectl get pods -n language-platform

# View logs
kubectl logs -f deployment/auth-service -n language-platform

# Scale service
kubectl scale deployment auth-service --replicas=5 -n language-platform

# Restart deployment
kubectl rollout restart deployment/auth-service -n language-platform

# Rollback
kubectl rollout undo deployment/auth-service -n language-platform
```

---

## 🔍 Health Checks

```bash
# All services
curl https://api.yourdomain.com/health

# Individual services
curl https://api.yourdomain.com/api/auth/health
curl https://api.yourdomain.com/api/v1/users/health
curl https://api.yourdomain.com/api/v1/courses/health
curl https://api.yourdomain.com/api/v1/tasks/health
curl https://api.yourdomain.com/api/v1/progress/health
curl https://api.yourdomain.com/api/v1/notifications/health
curl https://api.yourdomain.com/api/v1/files/health
curl https://api.yourdomain.com/api/v1/sessions/health
```

---

## 🧪 Testing API

### Register User
```bash
curl -X POST https://api.yourdomain.com/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "student@example.com",
    "password": "SecurePass123!",
    "confirmPassword": "SecurePass123!",
    "role": "Student"
  }'
```

### Login
```bash
curl -X POST https://api.yourdomain.com/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "student@example.com",
    "password": "SecurePass123!"
  }'
```

### Get Courses (with token)
```bash
curl https://api.yourdomain.com/api/v1/courses \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

---

## 📊 Monitoring

### Grafana Access
- URL: http://yourdomain.com:3000
- Username: admin
- Password: (from .env)

### Prometheus Access
- URL: http://yourdomain.com:9090

### Key Metrics
- Request rate
- Error rate
- Response time
- CPU/Memory usage
- Database connections
- Cache hit ratio

---

## 🔄 Backup & Restore

### Backup Database
```bash
# Manual backup
pg_dump -h db-host -U langplatform auth_db | gzip > auth_db_backup.sql.gz

# Automated (cron)
0 2 * * * /usr/local/bin/backup-db.sh
```

### Restore Database
```bash
gunzip < auth_db_backup.sql.gz | psql -h db-host -U langplatform auth_db
```

---

## 🚨 Troubleshooting

### Service won't start
```bash
# Check logs
docker-compose logs [service-name]
kubectl logs deployment/[service-name] -n language-platform

# Check resources
kubectl top pods -n language-platform
```

### Database connection failed
```bash
# Test connection
psql -h your-db-host -U langplatform -d auth_db

# Check secrets
kubectl get secret db-secret -n language-platform -o yaml
```

### SSL certificate issues
```bash
# Check certificate
sudo certbot certificates

# Renew manually
sudo certbot renew --force-renewal

# Test SSL
curl -I https://yourdomain.com
```

### Email not sending
```bash
# Test API key
curl -X GET "https://api.sendgrid.com/v3/user/profile" \
  -H "Authorization: Bearer $SENDGRID_API_KEY"

# Check logs
docker-compose logs notification-service
```

---

## 📱 Frontend Deployment

### Web (Next.js)
```bash
cd frontend/web
npm install
npm run build

# Deploy with PM2
pm2 start ../../config/ecosystem.config.js
pm2 save
pm2 startup
```

### Mobile (React Native)
```bash
cd frontend/mobile

# iOS
npm run ios
# Build for production: xcodebuild

# Android
npm run android
# Build for production: ./gradlew assembleRelease
```

---

## 🔐 Security Checklist

- [ ] Change default passwords
- [ ] Generate strong JWT secret
- [ ] Configure firewall (ports 22, 80, 443 only)
- [ ] Enable fail2ban
- [ ] Setup SSL certificates
- [ ] Configure CORS properly
- [ ] Enable rate limiting
- [ ] Setup monitoring alerts
- [ ] Regular security updates
- [ ] Backup encryption

---

## 📚 Documentation

- **Full Guide**: docs/DEPLOYMENT.md
- **Video Setup**: docs/VIDEO_SETUP.md
- **Checklist**: docs/CHECKLIST.md
- **Architecture**: docs/architecture.md
- **API Docs**: docs/api-contracts.md
- **Setup Guide**: docs/SETUP.md

---

## 💰 Cost Estimation

### Minimum (Small Platform)
- Server (2 vCPU, 4GB RAM): $20-40/month
- Managed PostgreSQL: $15-30/month
- Managed Redis: $10-20/month
- SendGrid (40k emails/month): $0-15/month
- Zoom Basic: $0 (40 min limit)
- S3 Storage: $5-10/month
- **Total: ~$50-115/month**

### Production (Medium Platform)
- Server (4 vCPU, 16GB RAM): $80-160/month
- Managed PostgreSQL: $50-100/month
- Managed Redis: $30-50/month
- SendGrid (100k emails): $15-20/month
- Zoom Pro (multi-host): $150-300/month
- S3 + CDN: $20-50/month
- **Total: ~$345-680/month**

### Enterprise (Large Platform)
- Kubernetes Cluster: $200-500/month
- RDS Multi-AZ: $200-400/month
- ElastiCache: $100-200/month
- SendGrid (500k+ emails): $50-100/month
- Agora (usage-based): $100-500/month
- S3 + CloudFront CDN: $100-300/month
- **Total: ~$750-2000/month**

---

## 🎯 Success Metrics

After deployment, verify:
- ✅ All 8 services running
- ✅ Health checks passing
- ✅ User registration working
- ✅ Course creation working
- ✅ File upload working
- ✅ Video sessions creating
- ✅ Emails sending
- ✅ SSL certificates valid
- ✅ Monitoring active
- ✅ No errors in logs

---

## 🆘 Support

- **Documentation**: See docs/ folder
- **DevOps Issues**: devops@yourdomain.com
- **Technical Support**: support@yourdomain.com
- **Emergency**: +1-XXX-XXX-XXXX

---

## 🚀 Next Steps After Deployment

1. **Configure monitoring alerts**
2. **Setup automated backups**
3. **Create admin user**
4. **Import sample courses** (optional)
5. **Test all user flows**
6. **Train support team**
7. **Prepare marketing materials**
8. **Launch! 🎉**

---

**Platform Status: Production Ready** ✅

Last Updated: 2026-02-04
Version: 1.0.0
