# Production Deployment Checklist

Use this checklist to ensure a smooth production deployment.

## Pre-Deployment

### Environment Setup
- [ ] Copy `.env.production.example` to `.env.production`
- [ ] Generate secure JWT_SECRET: `openssl rand -base64 32`
- [ ] Set strong POSTGRES_PASSWORD
- [ ] Set strong REDIS_PASSWORD
- [ ] Set strong MINIO credentials
- [ ] Configure EXPO_PUBLIC_API_URL (production API domain)
- [ ] Configure EXPO_PUBLIC_AI_SERVICE_URL (production AI service domain)
- [ ] (Optional) Set SENTRY_DSN for error tracking

### Infrastructure
- [ ] Domain names registered and DNS configured
- [ ] SSL/TLS certificates obtained (Let's Encrypt or purchased)
- [ ] Firewall rules configured (allow 80/443, block direct DB access)
- [ ] Server/cloud platform access configured
- [ ] Container registry account created (GHCR, Docker Hub, or ECR)
- [ ] Database backup strategy planned

### Code Preparation
- [ ] All tests passing locally
- [ ] Code reviewed and approved
- [ ] Version tagged (e.g., v1.0.0)
- [ ] Production branch up to date

## Build & Push Images

- [ ] Run `./scripts/push-images.sh --version 1.0.0` or trigger GitHub Actions
- [ ] Verify images pushed to registry successfully
- [ ] Tag images with semantic version
- [ ] Tag images with `latest`

## Deployment

### Initial Setup (First Time)
- [ ] Create databases (PostgreSQL, Redis)
- [ ] Run database migrations
- [ ] Seed initial data (if needed)
- [ ] Configure S3/MinIO buckets
- [ ] Set up monitoring and logging

### Deploy Services
- [ ] Pull production images
- [ ] Start PostgreSQL and verify health
- [ ] Start Redis and verify health
- [ ] Start MinIO/S3 and verify health
- [ ] Start backend service and verify health
- [ ] Start AI service and verify health
- [ ] Start frontend service and verify health

### Platform-Specific
**If using AWS ECS:**
- [ ] Create ECS cluster
- [ ] Register task definitions
- [ ] Create services with desired count
- [ ] Configure Application Load Balancer
- [ ] Set up CloudWatch logging
- [ ] Configure auto-scaling

**If using Fly.io:**
- [ ] Create Fly apps for each service
- [ ] Set secrets with `flyctl secrets set`
- [ ] Deploy with `flyctl deploy`
- [ ] Scale services as needed
- [ ] Configure custom domains

**If using Railway:**
- [ ] Create Railway project
- [ ] Connect GitHub repository
- [ ] Add PostgreSQL and Redis
- [ ] Set environment variables
- [ ] Deploy services
- [ ] Configure custom domains

**If using Docker Compose:**
- [ ] Copy production compose file to server
- [ ] Run `make prod-start` or `docker-compose -f docker-compose.production.yml up -d`
- [ ] Configure reverse proxy (Nginx/Traefik)
- [ ] Set up SSL with Certbot
- [ ] Configure systemd for auto-restart

## Post-Deployment

### Health Checks
- [ ] Backend health endpoint responding: `curl https://api.yourdomain.com/health`
- [ ] AI service health endpoint responding: `curl https://ai.yourdomain.com/health`
- [ ] Frontend health endpoint responding: `curl https://yourdomain.com/health`
- [ ] All services showing "healthy" status in Docker/platform
- [ ] Database connections working
- [ ] Redis connections working
- [ ] S3/MinIO accessible

### E2E Testing
- [ ] User signup works
- [ ] User login works
- [ ] Project creation works
- [ ] Blueprint upload works (test various file sizes)
- [ ] AI analysis completes successfully
- [ ] Bid generation works
- [ ] PDF download works
- [ ] PDF opens and displays correctly

### Cross-Browser Testing
- [ ] Chrome (desktop)
- [ ] Firefox (desktop)
- [ ] Safari (desktop)
- [ ] Edge (desktop)
- [ ] Mobile Safari (iOS)
- [ ] Mobile Chrome (Android)

### Performance Testing
- [ ] API response time < 200ms (p95)
- [ ] Blueprint upload < 30s for 10MB
- [ ] AI analysis < 2 minutes
- [ ] PDF generation < 15 seconds
- [ ] Page load < 3 seconds
- [ ] No memory leaks after 1 hour

### Security Validation
- [ ] HTTPS enforced on all endpoints
- [ ] Security headers present (X-Frame-Options, CSP, etc.)
- [ ] Authentication required for protected endpoints
- [ ] Authorization working correctly
- [ ] No secrets exposed in logs
- [ ] Database not publicly accessible
- [ ] Redis not publicly accessible
- [ ] File upload size limits working
- [ ] File type validation working
- [ ] Rate limiting configured (if applicable)

### Monitoring Setup
- [ ] Error tracking configured (Sentry)
- [ ] Log aggregation configured
- [ ] Uptime monitoring configured
- [ ] Performance monitoring configured
- [ ] Alerts configured for critical errors
- [ ] Alerts configured for service downtime
- [ ] Disk space monitoring configured
- [ ] Database performance monitoring configured

### Backup & Recovery
- [ ] Database backup script tested
- [ ] Database restore process tested
- [ ] Backup schedule configured (daily/hourly)
- [ ] Backup storage configured (S3/cloud storage)
- [ ] Backup retention policy set
- [ ] Disaster recovery plan documented

### Documentation
- [ ] Deployment process documented
- [ ] Environment variables documented
- [ ] Troubleshooting guide updated
- [ ] Runbook created for common issues
- [ ] Contact information for on-call support

## Launch Readiness

### Final Checks
- [ ] All above items completed
- [ ] Load testing completed successfully
- [ ] Security audit passed
- [ ] Stakeholder approval obtained
- [ ] Support team trained
- [ ] Rollback plan tested
- [ ] Communication plan ready

### Go Live
- [ ] Update DNS to point to production (if not already)
- [ ] Monitor logs during first hour
- [ ] Monitor performance metrics
- [ ] Watch for errors in Sentry
- [ ] Be available for immediate issues

### Post-Launch
- [ ] Monitor for 24 hours continuously
- [ ] Collect user feedback
- [ ] Document any issues encountered
- [ ] Create tickets for non-critical issues
- [ ] Plan for next iteration
- [ ] Celebrate! 🎉

## Rollback Procedure

If critical issues are found:

1. [ ] Immediately notify team
2. [ ] Assess severity of issue
3. [ ] If needed, rollback to previous version:
   ```bash
   export VERSION=1.0.0
   docker-compose -f docker-compose.production.yml up -d
   ```
4. [ ] Verify rollback successful
5. [ ] Document issue for post-mortem
6. [ ] Fix issue in development
7. [ ] Re-deploy when ready

## Support Contacts

- **DevOps Lead**: [Name/Email]
- **Backend Lead**: [Name/Email]
- **Frontend Lead**: [Name/Email]
- **On-Call Support**: [Phone/Slack]
- **Escalation**: [Management Contact]

## Resources

- [DEPLOYMENT.md](./DEPLOYMENT.md) - Complete deployment guide
- [E2E_TESTING.md](./E2E_TESTING.md) - Testing procedures
- [QUICKSTART.md](./QUICKSTART.md) - Quick reference
- [M7_IMPLEMENTATION_SUMMARY.md](./M7_IMPLEMENTATION_SUMMARY.md) - Implementation details

---

**Last Updated**: December 2024
**Version**: 1.0.0
