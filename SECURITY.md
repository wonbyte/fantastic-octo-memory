# Security Policy

## Supported Versions

We release patches for security vulnerabilities for the following versions:

| Version | Supported          |
| ------- | ------------------ |
| 1.0.x   | :white_check_mark: |
| < 1.0   | :x:                |

## Reporting a Vulnerability

We take the security of our construction estimation platform seriously. If you believe you have found a security vulnerability, please report it to us as described below.

### How to Report

**Please DO NOT report security vulnerabilities through public GitHub issues.**

Instead, please report them via one of the following methods:

1. **Preferred Method - Private Security Advisory:**
   - Go to the [Security tab](https://github.com/wonbyte/fantastic-octo-memory/security) 
   - Click "Report a vulnerability"
   - Fill in the details using the template

2. **Email:**
   - Send details to: [security@example.com] (replace with actual security contact)
   - Include "[SECURITY]" in the subject line

### What to Include

Please include the following information in your report:

- **Type of vulnerability** (e.g., XSS, SQL injection, authentication bypass)
- **Affected component(s)** (Backend, AI Service, Frontend)
- **Full path to the source file(s)** related to the issue
- **Location of affected source code** (tag/branch/commit or direct URL)
- **Step-by-step instructions** to reproduce the issue
- **Proof-of-concept or exploit code** (if possible)
- **Impact of the issue** and potential attack scenarios
- **Any suggested mitigation or remediation steps**

### What to Expect

When you report a vulnerability, we will:

1. **Acknowledge receipt** within 48 hours
2. **Provide an initial assessment** within 5 business days
3. **Keep you informed** about the progress toward a fix
4. **Notify you** when the vulnerability is fixed
5. **Credit you** in the security advisory (if desired)

### Our Commitment

- We will investigate all legitimate reports
- We will keep you informed of our progress
- We will not take legal action against researchers who:
  - Report in good faith
  - Avoid privacy violations
  - Avoid service disruption
  - Don't access data beyond what's necessary to demonstrate the vulnerability

## Security Best Practices for Users

### For Developers

1. **Keep dependencies up to date:**
   ```bash
   # Backend
   cd backend && go get -u ./...
   
   # AI Service  
   cd ai_service && pip install -U -r requirements.txt
   
   # Frontend
   cd app && npm update
   ```

2. **Use environment variables for secrets:**
   - Never commit `.env` files
   - Use `.env.example` as a template only
   - Rotate secrets regularly

3. **Enable security features:**
   - Rate limiting in production
   - Security headers (HSTS, CSP)
   - CORS with specific origins only

### For Deployment

1. **Production Security Checklist:**
   - [ ] Change all default passwords
   - [ ] Generate secure JWT_SECRET (min 32 characters)
   - [ ] Configure HTTPS/TLS
   - [ ] Set up database backups
   - [ ] Enable error tracking (Sentry)
   - [ ] Configure rate limiting
   - [ ] Set security headers
   - [ ] Restrict CORS to production domains
   - [ ] Review file upload limits
   - [ ] Enable audit logging

2. **Infrastructure Security:**
   - Use managed database services with encryption at rest
   - Enable network isolation (VPC, security groups)
   - Use secrets manager for sensitive data
   - Enable DDoS protection
   - Set up monitoring and alerting
   - Regular security scans

3. **Access Control:**
   - Implement least privilege principle
   - Use strong authentication methods
   - Enable 2FA for admin accounts
   - Regular access audits
   - Separate production and development credentials

## Known Security Features

### Backend (Go)
- JWT-based authentication with configurable expiry
- Bcrypt password hashing (cost factor: 10)
- SQL injection protection via parameterized queries
- Rate limiting per IP and user
- Request body size limits (default: 10MB)
- Security headers middleware (HSTS, CSP, X-Frame-Options)
- Structured logging for audit trails
- Sentry integration for error tracking

### AI Service (Python)
- Input validation on all endpoints
- File type validation for uploads
- Structured logging
- CORS configuration
- Error tracking with Sentry

### Frontend
- Secure token storage (SecureStore on mobile, localStorage with encryption for web)
- HTTPS enforcement in production
- XSS protection via React's auto-escaping
- Input sanitization
- Offline data encryption

## Security Headers

The platform implements the following security headers:

```
Strict-Transport-Security: max-age=31536000; includeSubDomains
Content-Security-Policy: default-src 'self'; script-src 'self' 'unsafe-inline'
X-Frame-Options: DENY
X-Content-Type-Options: nosniff
Referrer-Policy: strict-origin-when-cross-origin
```

## Dependency Security

We use the following tools to monitor dependencies:

- **Go**: `go list -m all` with manual reviews
- **Python**: `pip audit` and safety checks
- **Node.js**: `npm audit` in CI/CD

### Updating Dependencies

We aim to:
- Review and update dependencies monthly
- Apply security patches within 48 hours of disclosure
- Keep all dependencies within one major version of latest

## Security Audit History

| Date       | Type           | Findings | Status   |
|------------|----------------|----------|----------|
| 2025-12-19 | Internal Audit | 0 High   | Complete |

## Vulnerability Disclosure Timeline

We follow this timeline for vulnerability disclosure:

1. **Day 0**: Vulnerability reported
2. **Day 1-2**: Initial triage and acknowledgment
3. **Day 3-7**: Investigation and patch development
4. **Day 8-14**: Testing and validation
5. **Day 15**: Patch release and public disclosure (if applicable)
6. **Day 30**: Detailed write-up (optional, for educational purposes)

## Bug Bounty Program

Currently, we do not have a formal bug bounty program. However, we deeply appreciate security researchers who help us keep our platform secure and will:

- Publicly acknowledge your contribution (if desired)
- Provide detailed feedback on your report
- Consider swag/recognition for significant findings

## Additional Resources

- [OWASP Top 10](https://owasp.org/www-project-top-ten/)
- [OWASP API Security Top 10](https://owasp.org/www-project-api-security/)
- [CWE Top 25](https://cwe.mitre.org/top25/)
- [Security Hardening Guide](./M7_SECURITY_HARDENING.md)

## Security Contacts

- **Security Team**: [security@example.com] (replace with actual contact)
- **Project Maintainer**: [@wonbyte](https://github.com/wonbyte)

## Legal

This security policy is subject to change without notice. Please check back regularly for updates.

Last Updated: December 19, 2025
