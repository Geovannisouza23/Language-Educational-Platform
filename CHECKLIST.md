# Quick Start Checklist

## Pre-Deployment Checklist

### 1. Domain & DNS ✓
- [ ] Domain registered (yourdomain.com)
- [ ] DNS A records configured
- [ ] Nameservers updated (allow 24-48h)

### 2. Server Setup ✓
- [ ] Ubuntu 22.04 server provisioned
- [ ] SSH key authentication configured
- [ ] Firewall configured (ports 22, 80, 443)
- [ ] Docker installed
- [ ] kubectl installed (for K8s)

### 3. Database Setup ✓
- [ ] PostgreSQL instance created (AWS RDS/self-hosted)
- [ ] 8 databases created (auth_db, users_db, etc.)
- [ ] User and permissions configured
- [ ] Connection tested
- [ ] Backups configured

### 4. Redis Setup ✓
- [ ] Redis instance created (ElastiCache/self-hosted)
- [ ] Password configured (if applicable)
- [ ] Connection tested

### 5. SendGrid Configuration ✓
- [ ] SendGrid account created
- [ ] API key generated
- [ ] Sender email verified
- [ ] Email templates created
- [ ] Test email sent successfully

### 6. Video Provider Setup ✓
Choose ONE:
- [ ] Zoom: App created, OAuth configured, credentials obtained
- [ ] Agora: Project created, App ID obtained, token generation tested
- [ ] Jitsi: Self-hosted installed OR using meet.jit.si

### 7. File Storage Setup ✓
- [ ] AWS S3 bucket created OR MinIO installed
- [ ] Access keys generated
- [ ] Bucket permissions configured
- [ ] CDN configured (optional but recommended)

### 8. SSL Certificates ✓
- [ ] Let's Encrypt certbot installed
- [ ] SSL certificates obtained
- [ ] Auto-renewal tested
- [ ] HTTPS redirect configured

---

## Deployment Checklist

### Step 1: Clone Repository
```bash
git clone <your-repo-url>
cd user.api
```
- [ ] Repository cloned
- [ ] On main/master branch

### Step 2: Environment Configuration
```bash
cp .env.example .env.production
nano .env.production
```
Required variables:
- [ ] DATABASE_HOST, DATABASE_PASSWORD
- [ ] REDIS_HOST, REDIS_PASSWORD
- [ ] JWT_SECRET (generated)
- [ ] SENDGRID_API_KEY
- [ ] Video provider credentials (ZOOM/AGORA/JITSI)
- [ ] AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY
- [ ] PRIMARY_DOMAIN, API_DOMAIN
- [ ] All template IDs configured

### Step 3: Generate Secrets
```bash
export JWT_SECRET=$(openssl rand -hex 64)
export DB_PASSWORD=$(openssl rand -base64 32)
export ENCRYPTION_KEY=$(openssl rand -hex 32)
```
- [ ] JWT_SECRET generated
- [ ] DB_PASSWORD generated
- [ ] ENCRYPTION_KEY generated
- [ ] Secrets added to .env.production

### Step 4: Setup SendGrid
```bash
export SENDGRID_API_KEY="your-key"
./scripts/setup-sendgrid.sh
```
- [ ] Script executed successfully
- [ ] Template IDs saved
- [ ] Test email received

### Step 5: Setup SSL
```bash
./scripts/setup-ssl.sh
```
- [ ] Certificates obtained
- [ ] Auto-renewal configured
- [ ] HTTPS working

### Step 6: Configure NGINX
```bash
sudo cp infra/nginx/nginx-production.conf /etc/nginx/nginx.conf
sudo nano /etc/nginx/nginx.conf  # Update domains
sudo nginx -t
sudo systemctl reload nginx
```
- [ ] Config copied
- [ ] Domains updated
- [ ] Config tested
- [ ] NGINX reloaded

### Step 7: Deploy Services

**Option A: Docker Compose**
```bash
docker-compose -f infra/docker/docker-compose.yml up -d
```
- [ ] All services started
- [ ] Health checks passing

**Option B: Kubernetes**
```bash
./scripts/setup-k8s-secrets.sh
./scripts/deploy-production.sh production
```
- [ ] Secrets created
- [ ] Deployments successful
- [ ] Pods running

### Step 8: Run Migrations
```bash
# Auth service
docker exec auth-service dotnet ef database update

# Other services auto-migrate on startup
```
- [ ] Auth DB migrated
- [ ] All tables created

### Step 9: Deploy Frontend
```bash
cd frontend/web
npm install
npm run build
# Deploy to server
pm2 start ecosystem.config.js
```
- [ ] Dependencies installed
- [ ] Build successful
- [ ] Frontend running

### Step 10: Smoke Tests
```bash
# Health checks
curl https://api.yourdomain.com/health
curl https://api.yourdomain.com/api/auth/health

# Test registration
curl -X POST https://api.yourdomain.com/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"Test123!","role":"Student"}'
```
- [ ] All health checks passing
- [ ] Registration working
- [ ] Login working
- [ ] Course creation working

---

## Post-Deployment Checklist

### Monitoring ✓
- [ ] Grafana accessible
- [ ] Dashboards imported
- [ ] Alerts configured
- [ ] Sentry error tracking (optional)

### Backups ✓
- [ ] Database backup script configured
- [ ] Backup schedule set (cron)
- [ ] Backup restoration tested
- [ ] S3 backup bucket created

### Security ✓
- [ ] Firewall rules applied
- [ ] Fail2ban configured
- [ ] SSH key-only authentication
- [ ] Rate limiting tested
- [ ] CORS configured correctly

### Performance ✓
- [ ] CDN configured
- [ ] Gzip compression enabled
- [ ] Database indexes created
- [ ] Redis caching working

### Documentation ✓
- [ ] API documentation accessible
- [ ] Runbook created
- [ ] Support team trained
- [ ] Incident response plan

---

## Verification Commands

### Service Health
```bash
# Check all services
kubectl get pods -n language-platform
docker-compose ps

# Test each endpoint
curl https://api.yourdomain.com/health
curl https://api.yourdomain.com/api/auth/health
curl https://api.yourdomain.com/api/v1/users/health
curl https://api.yourdomain.com/api/v1/courses/health
curl https://api.yourdomain.com/api/v1/tasks/health
curl https://api.yourdomain.com/api/v1/progress/health
curl https://api.yourdomain.com/api/v1/notifications/health
curl https://api.yourdomain.com/api/v1/files/health
curl https://api.yourdomain.com/api/v1/sessions/health
```

### Database Connectivity
```bash
# Test PostgreSQL
psql -h your-db-host -U langplatform -d auth_db -c "SELECT 1;"

# Test Redis
redis-cli -h your-redis-host ping
```

### SSL Certificate
```bash
# Check certificate
sudo certbot certificates

# Test SSL
curl -I https://yourdomain.com
curl -I https://api.yourdomain.com
```

### Load Testing (Optional)
```bash
# Install Apache Bench
sudo apt install apache2-utils

# Test API
ab -n 1000 -c 10 https://api.yourdomain.com/health
```

---

## Rollback Plan

If deployment fails:

1. **Check logs first**
```bash
kubectl logs -f deployment/service-name -n language-platform
docker-compose logs service-name
```

2. **Rollback Kubernetes**
```bash
kubectl rollout undo deployment/service-name -n language-platform
```

3. **Rollback Docker**
```bash
docker-compose down
git checkout previous-commit
docker-compose up -d
```

4. **Restore Database**
```bash
psql -h db-host -U langplatform -d auth_db < backup.sql
```

---

## Success Criteria

✅ All services responding to health checks
✅ User registration and login working
✅ Course creation and enrollment working
✅ File upload working
✅ Video session creation working
✅ Email notifications being sent
✅ Frontend accessible and functional
✅ SSL certificates valid
✅ Monitoring dashboards showing data
✅ No errors in logs

---

## Support Contacts

- **DevOps**: devops@yourdomain.com
- **Backend**: backend@yourdomain.com
- **Frontend**: frontend@yourdomain.com
- **Emergency**: +1-XXX-XXX-XXXX

---

## Estimated Time

- Initial setup: 4-6 hours
- Configuration: 2-3 hours
- Deployment: 1-2 hours
- Testing & verification: 1-2 hours
- **Total: 8-13 hours**

---

**Ready for production!** 🚀
