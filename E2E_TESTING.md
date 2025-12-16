# End-to-End Testing & Validation Guide

Complete guide for validating the full user journey from signup to PDF download in production.

## Table of Contents

- [Overview](#overview)
- [Testing Environments](#testing-environments)
- [Pre-Deployment Checklist](#pre-deployment-checklist)
- [E2E Test Scenarios](#e2e-test-scenarios)
- [User Testing Checklist](#user-testing-checklist)
- [Automated E2E Tests](#automated-e2e-tests)
- [Performance Testing](#performance-testing)
- [Security Testing](#security-testing)
- [Troubleshooting](#troubleshooting)

---

## Overview

### Complete User Flow

The platform supports the following end-to-end user journey:

```
1. User Signup/Login
   ↓
2. Project Creation
   ↓
3. Blueprint Upload
   ↓
4. AI Analysis & Processing
   ↓
5. Bid Generation
   ↓
6. PDF Download
```

Each step must be validated in production to ensure the system works correctly.

---

## Testing Environments

### Development
- **URL**: http://localhost:3000
- **Purpose**: Local development and initial testing
- **Data**: Mock/test data

### Staging (Recommended)
- **URL**: https://staging.yourdomain.com
- **Purpose**: Pre-production testing with production-like setup
- **Data**: Anonymized production-like data

### Production
- **URL**: https://yourdomain.com
- **Purpose**: Live user environment
- **Data**: Real user data

---

## Pre-Deployment Checklist

Before running E2E tests, verify:

### Infrastructure
- [ ] All services are running and healthy
  ```bash
  curl https://api.yourdomain.com/health
  curl https://ai.yourdomain.com/health
  curl https://yourdomain.com/health
  ```
- [ ] Database migrations completed successfully
- [ ] Redis is accessible
- [ ] S3/MinIO storage is configured and accessible
- [ ] SSL/TLS certificates are valid

### Configuration
- [ ] Environment variables set correctly
- [ ] CORS configured for frontend domain
- [ ] API URLs configured in frontend
- [ ] File upload size limits appropriate
- [ ] JWT token expiry set correctly

### Security
- [ ] Secrets rotated from defaults
- [ ] Authentication working
- [ ] Authorization enforced on protected endpoints
- [ ] Rate limiting enabled
- [ ] Security headers configured

---

## E2E Test Scenarios

### Scenario 1: New User Registration & First Project

**Objective**: Verify complete new user onboarding flow

#### Steps:

**1. User Signup** ✅
```bash
# API Test
curl -X POST https://api.yourdomain.com/auth/signup \
  -H "Content-Type: application/json" \
  -d '{
    "email": "testuser@example.com",
    "password": "SecureP@ssw0rd",
    "name": "Test User",
    "company_name": "Test Construction Co"
  }'

# Expected Response: 201 Created
# {
#   "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
#   "user": {
#     "id": "uuid",
#     "email": "testuser@example.com",
#     "name": "Test User",
#     "company_name": "Test Construction Co"
#   }
# }
```

**Manual Test:**
1. Open https://yourdomain.com
2. Click "Sign Up" button
3. Enter email, password, name, company name
4. Click "Create Account"
5. Verify redirect to dashboard
6. Verify user info displayed correctly

**Expected Results:**
- ✅ User account created
- ✅ JWT token received and stored
- ✅ User redirected to dashboard
- ✅ No console errors

---

**2. Project Creation** ✅
```bash
# Save token from signup
TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."

# Create project
curl -X POST https://api.yourdomain.com/api/projects \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Downtown Office Complex",
    "description": "5-story office building renovation",
    "location": "123 Main St, City, State",
    "client_name": "Acme Properties"
  }'

# Expected Response: 201 Created
# {
#   "id": "project-uuid",
#   "name": "Downtown Office Complex",
#   "description": "5-story office building renovation",
#   "status": "active",
#   "created_at": "2024-01-01T00:00:00Z"
# }
```

**Manual Test:**
1. Click "New Project" button
2. Fill in project details:
   - Name: "Downtown Office Complex"
   - Description: "5-story office building renovation"
   - Location: "123 Main St, City, State"
   - Client: "Acme Properties"
3. Click "Create Project"
4. Verify project appears in project list
5. Click project to open details

**Expected Results:**
- ✅ Project created successfully
- ✅ Project visible in list
- ✅ Project details page loads
- ✅ Status shows "Active"

---

**3. Blueprint Upload** ✅
```bash
# Upload blueprint
curl -X POST https://api.yourdomain.com/api/blueprints/upload \
  -H "Authorization: Bearer $TOKEN" \
  -F "project_id=project-uuid" \
  -F "file=@test-blueprint.pdf" \
  -F "name=Floor Plan - Level 1"

# Expected Response: 201 Created
# {
#   "id": "blueprint-uuid",
#   "project_id": "project-uuid",
#   "name": "Floor Plan - Level 1",
#   "filename": "floor-plan-level-1.pdf",
#   "file_size": 2458624,
#   "status": "uploaded",
#   "upload_url": "https://...",
#   "created_at": "2024-01-01T00:00:00Z"
# }
```

**Manual Test:**
1. Open project details page
2. Click "Upload Blueprint" button
3. Select PDF file (recommended: valid architectural blueprint, 1-10 MB)
4. Add name: "Floor Plan - Level 1"
5. Click "Upload"
6. Verify upload progress indicator
7. Wait for upload completion
8. Verify blueprint appears in list

**Expected Results:**
- ✅ File upload starts immediately
- ✅ Progress indicator shows percentage
- ✅ Upload completes within 30 seconds (for <10MB files)
- ✅ Blueprint appears in project blueprints list
- ✅ Status shows "Uploaded"
- ✅ Thumbnail/preview available

**Test Files:**
- Small PDF: 1-2 MB architectural drawing
- Medium PDF: 5-7 MB multi-page blueprint
- Large PDF: 9-10 MB (test size limits)

---

**4. AI Analysis Trigger** ✅
```bash
# Trigger analysis
curl -X POST https://api.yourdomain.com/api/blueprints/blueprint-uuid/analyze \
  -H "Authorization: Bearer $TOKEN"

# Expected Response: 202 Accepted
# {
#   "job_id": "job-uuid",
#   "blueprint_id": "blueprint-uuid",
#   "status": "queued",
#   "message": "Analysis job queued successfully"
# }
```

**Manual Test:**
1. Click "Analyze Blueprint" button
2. Verify confirmation dialog
3. Click "Confirm"
4. Verify status changes to "Analyzing"
5. Wait for analysis to complete (typically 30-120 seconds)

**Expected Results:**
- ✅ Analysis job created
- ✅ Status indicator shows "Analyzing"
- ✅ Real-time status updates (if WebSocket/polling enabled)
- ✅ Progress messages appear

---

**5. Check Analysis Status** ✅
```bash
# Poll for status
curl -X GET https://api.yourdomain.com/api/jobs/job-uuid \
  -H "Authorization: Bearer $TOKEN"

# Response (in progress):
# {
#   "id": "job-uuid",
#   "blueprint_id": "blueprint-uuid",
#   "status": "processing",
#   "progress": 45,
#   "message": "Extracting text from blueprint..."
# }

# Response (completed):
# {
#   "id": "job-uuid",
#   "blueprint_id": "blueprint-uuid",
#   "status": "completed",
#   "progress": 100,
#   "result": {
#     "total_area_sqft": 5000,
#     "materials_detected": [...],
#     "cost_estimate": 250000
#   },
#   "completed_at": "2024-01-01T00:05:00Z"
# }
```

**Manual Test:**
1. Monitor status updates on blueprint page
2. Verify progress messages
3. Wait for completion
4. Verify analysis results display

**Expected Results:**
- ✅ Status updates in real-time or on refresh
- ✅ Analysis completes within 2 minutes (for typical blueprints)
- ✅ Results display correctly:
  - Total area detected
  - Materials list
  - Cost estimates
  - Confidence scores
- ✅ No errors in console
- ✅ Status changes to "Completed"

---

**6. Bid Generation** ✅
```bash
# Generate bid from analysis
curl -X POST https://api.yourdomain.com/api/bids/generate \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "project_id": "project-uuid",
    "blueprint_id": "blueprint-uuid",
    "bid_name": "Initial Estimate - Downtown Office",
    "markup_percentage": 15
  }'

# Expected Response: 201 Created
# {
#   "id": "bid-uuid",
#   "project_id": "project-uuid",
#   "blueprint_id": "blueprint-uuid",
#   "bid_name": "Initial Estimate - Downtown Office",
#   "total_cost": 287500,
#   "status": "draft",
#   "pdf_url": null,
#   "created_at": "2024-01-01T00:06:00Z"
# }
```

**Manual Test:**
1. Click "Generate Bid" button
2. Enter bid details:
   - Bid name: "Initial Estimate - Downtown Office"
   - Markup: 15%
   - Contingency: 10%
3. Review summary
4. Click "Generate"
5. Wait for bid generation (typically 5-15 seconds)

**Expected Results:**
- ✅ Bid form displays with pre-filled data from analysis
- ✅ Cost calculations update dynamically
- ✅ Bid created successfully
- ✅ Bid appears in project bids list
- ✅ Status shows "Draft"

---

**7. PDF Generation & Download** ✅
```bash
# Request PDF generation
curl -X POST https://api.yourdomain.com/api/bids/bid-uuid/generate-pdf \
  -H "Authorization: Bearer $TOKEN"

# Expected Response: 202 Accepted
# {
#   "bid_id": "bid-uuid",
#   "status": "generating",
#   "message": "PDF generation in progress"
# }

# Check status
curl -X GET https://api.yourdomain.com/api/bids/bid-uuid \
  -H "Authorization: Bearer $TOKEN"

# Response (completed):
# {
#   "id": "bid-uuid",
#   "pdf_url": "https://...presigned-url...",
#   "pdf_generated_at": "2024-01-01T00:06:30Z",
#   "status": "completed"
# }

# Download PDF
curl -o bid.pdf "https://...presigned-url..."
```

**Manual Test:**
1. Click "Generate PDF" button
2. Wait for PDF generation (typically 5-10 seconds)
3. Click "Download PDF" button
4. Verify PDF downloads successfully
5. Open PDF and verify:
   - Company branding/logo
   - Project details
   - Blueprint reference
   - Cost breakdown
   - Line items
   - Total cost
   - Terms and conditions
   - Professional formatting

**Expected Results:**
- ✅ PDF generates within 15 seconds
- ✅ Download link available
- ✅ PDF downloads successfully
- ✅ PDF opens without errors
- ✅ All data present and correctly formatted
- ✅ File size reasonable (typically 50-500 KB)
- ✅ PDF is printable
- ✅ Professional appearance

---

### Scenario 2: Existing User Login & Multiple Projects

**Objective**: Verify login and multi-project management

#### Steps:

**1. User Login**
```bash
curl -X POST https://api.yourdomain.com/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "testuser@example.com",
    "password": "SecureP@ssw0rd"
  }'
```

**Manual Test:**
1. Navigate to login page
2. Enter credentials
3. Click "Login"
4. Verify redirect to dashboard
5. Verify existing projects displayed

**Expected Results:**
- ✅ Login successful
- ✅ Token received
- ✅ Dashboard loads
- ✅ Projects list populated

---

**2. Multiple Blueprint Management**

**Manual Test:**
1. Create 3 different projects
2. Upload 2-3 blueprints to each project
3. Trigger analysis on multiple blueprints
4. Monitor concurrent processing
5. Generate bids for completed analyses

**Expected Results:**
- ✅ Multiple uploads work simultaneously
- ✅ Analyses queue and process correctly
- ✅ No race conditions or conflicts
- ✅ Each project isolated correctly
- ✅ Performance remains acceptable

---

### Scenario 3: Error Handling & Edge Cases

**Objective**: Verify system handles errors gracefully

#### Test Cases:

**1. Invalid File Upload**
- Upload non-PDF file → Should reject
- Upload 0-byte file → Should reject
- Upload 50MB file → Should reject (size limit)
- Upload corrupted PDF → Should detect and error

**2. Authentication Errors**
- Expired token → Should redirect to login
- Invalid token → Should return 401
- Missing token → Should reject

**3. Network Issues**
- Slow connection during upload → Should show progress
- Connection loss mid-upload → Should fail gracefully
- API timeout → Should show error message

**4. Concurrent Users**
- Multiple users uploading simultaneously
- Race conditions in job processing
- Database connection pooling

---

## User Testing Checklist

### Device Testing

Test on the following devices/browsers:

#### Desktop Browsers
- [ ] Chrome (latest)
- [ ] Firefox (latest)
- [ ] Safari (latest)
- [ ] Edge (latest)

#### Mobile Devices
- [ ] iOS Safari (iPhone)
- [ ] Android Chrome
- [ ] Tablet (iPad)

### Feature Checklist

For each device/browser, verify:

#### Authentication
- [ ] Sign up works
- [ ] Login works
- [ ] Logout works
- [ ] Session persists on refresh
- [ ] Expired token handled correctly
- [ ] Password validation works

#### Project Management
- [ ] Create project
- [ ] Edit project
- [ ] Delete project
- [ ] Project list loads
- [ ] Search/filter projects
- [ ] Sort projects

#### Blueprint Upload
- [ ] File picker opens
- [ ] Upload progress shows
- [ ] Multiple files supported
- [ ] Large files handle correctly
- [ ] Error messages clear
- [ ] Cancel upload works

#### Analysis
- [ ] Trigger analysis
- [ ] Status updates
- [ ] Results display correctly
- [ ] Error handling
- [ ] Retry on failure

#### Bid Generation
- [ ] Form validation
- [ ] Cost calculations
- [ ] Preview available
- [ ] Edit bid
- [ ] Delete bid

#### PDF Download
- [ ] Generate PDF
- [ ] Download link works
- [ ] PDF opens correctly
- [ ] PDF content accurate
- [ ] Re-download works

### UI/UX Checklist

- [ ] All buttons clickable
- [ ] Forms submit correctly
- [ ] Validation messages clear
- [ ] Loading states shown
- [ ] Error messages helpful
- [ ] Success messages displayed
- [ ] Navigation intuitive
- [ ] Responsive on mobile
- [ ] No layout breaks
- [ ] No console errors
- [ ] Fast page loads (<3s)
- [ ] Smooth interactions

### Accessibility
- [ ] Keyboard navigation works
- [ ] Screen reader compatible
- [ ] Sufficient color contrast
- [ ] Focus indicators visible
- [ ] Error messages announced
- [ ] Form labels present

---

## Automated E2E Tests

### Using Playwright (Recommended)

The project includes comprehensive Playwright E2E tests in the `e2e/` directory:

**Test Files:**
- `e2e/basic.spec.ts` - Basic functionality and accessibility tests
- `e2e/user-journey.spec.ts` - Complete user flow tests
- `e2e/revision-comparison.spec.ts` - Revision comparison feature tests
- `e2e/complete-workflow.spec.ts` - Comprehensive workflow from signup to PDF download

**Test Coverage:**
- ✅ Authentication (signup, login, logout)
- ✅ Project creation and management
- ✅ Blueprint upload and analysis
- ✅ Bid generation and export
- ✅ PDF download
- ✅ Error handling and edge cases
- ✅ Mobile responsiveness
- ✅ Accessibility features
- ✅ Dark mode
- ✅ Offline mode

### Running E2E Tests

```bash
# Install dependencies (first time only)
npm install
npx playwright install --with-deps

# Run all E2E tests
npm run test:e2e

# Run with UI mode (interactive)
npm run test:e2e:ui

# Run in headed mode (see browser)
npm run test:e2e:headed

# View test report
npm run test:e2e:report

# Run specific test file
npx playwright test e2e/complete-workflow.spec.ts

# Run tests in specific browser
npx playwright test --project=chromium
npx playwright test --project=firefox
npx playwright test --project=webkit

# Debug mode
npx playwright test --debug

# Run on specific URL
BASE_URL=https://staging.yourdomain.com npm run test:e2e
```

### Configuration

The Playwright configuration is in `playwright.config.ts`:

```typescript
export default defineConfig({
  testDir: './e2e',
  fullyParallel: true,
  retries: process.env.CI ? 2 : 0,
  workers: process.env.CI ? 1 : undefined,
  reporter: 'html',
  use: {
    baseURL: process.env.BASE_URL || 'http://localhost:3000',
    trace: 'on-first-retry',
    screenshot: 'only-on-failure',
  },
  projects: [
    { name: 'chromium', use: { ...devices['Desktop Chrome'] } },
    { name: 'firefox', use: { ...devices['Desktop Firefox'] } },
    { name: 'webkit', use: { ...devices['Desktop Safari'] } },
    { name: 'Mobile Chrome', use: { ...devices['Pixel 5'] } },
    { name: 'Mobile Safari', use: { ...devices['iPhone 12'] } },
  ],
});
```

### CI/CD Integration

E2E tests run automatically in CI/CD:

```yaml
# .github/workflows/ci.yml
e2e-tests:
  name: E2E Tests (Playwright)
  runs-on: ubuntu-latest
  steps:
    - uses: actions/checkout@v4
    - name: Set up Node.js
      uses: actions/setup-node@v4
      with:
        node-version: '22'
    - name: Install dependencies
      run: npm install
    - name: Install Playwright Browsers
      run: npx playwright install --with-deps
    - name: Run E2E tests
      run: npm run test:e2e
```

### Test Results and Artifacts

After running tests:
- **HTML Report**: `playwright-report/index.html`
- **Screenshots**: `test-results/` (on failure)
- **Traces**: Available in report for debugging
- **Videos**: Captured for failed tests (if configured)

To view the report:
```bash
npm run test:e2e:report
```

---

## Performance Testing

### Load Testing with Artillery

The project includes comprehensive load testing configurations in the `load-tests/` directory.

**Configuration Files:**
- `load-tests/artillery-backend.yml` - Backend API load tests
- `load-tests/artillery-ai-service.yml` - AI service load tests
- `load-tests/README.md` - Detailed load testing guide

### Running Load Tests

```bash
# Install Artillery
npm install

# Test backend API
npm run test:load:backend

# Test AI service
npm run test:load:ai

# Run all load tests
npm run test:load:all

# Run against custom target
artillery run --target https://api.yourdomain.com load-tests/artillery-backend.yml

# Run with output report
artillery run --output report.json load-tests/artillery-backend.yml
artillery report report.json
```

### Load Test Scenarios

**Backend API Tests:**
- Health check monitoring
- User authentication flow (signup/login)
- Project CRUD operations
- Blueprint workflow
- Bid generation and export

**AI Service Tests:**
- Health check
- Blueprint analysis (with S3)
- Bid generation with takeoff data
- Mixed operations
- Concurrent processing

### Performance Targets

**Backend API:**
- P95 response time: < 500ms
- P99 response time: < 1000ms
- Error rate: < 1%
- Throughput: > 100 RPS

**AI Service:**
- P95 response time: < 2s
- P99 response time: < 5s
- Error rate: < 2%
- Throughput: > 20 RPS

For detailed load testing documentation, see [load-tests/README.md](./load-tests/README.md)
      name: Sustained load
  variables:
    email: "loadtest-{{ $randomString() }}@example.com"
    password: "LoadTest123!"

scenarios:
  - name: "Complete User Flow"
    flow:
      - post:
          url: "/auth/signup"
          json:
            email: "{{ email }}"
            password: "{{ password }}"
            name: "Load Test User"
          capture:
            - json: "$.token"
              as: "token"
      
      - post:
          url: "/api/projects"
          headers:
            Authorization: "Bearer {{ token }}"
          json:
            name: "Load Test Project"
            description: "Performance testing"
      
      - think: 5
      
      - get:
          url: "/api/projects"
          headers:
            Authorization: "Bearer {{ token }}"
```

Run load test:
```bash
npm install -g artillery
artillery run artillery-config.yml
```

### Performance Targets

- **API Response Time**: < 200ms (p95)
- **Blueprint Upload**: < 30s for 10MB file
- **Analysis Completion**: < 2 minutes
- **PDF Generation**: < 15 seconds
- **Page Load**: < 3 seconds
- **Time to Interactive**: < 5 seconds

---

## Security Testing

### Penetration Testing Checklist

- [ ] SQL Injection attempts blocked
- [ ] XSS attempts sanitized
- [ ] CSRF protection working
- [ ] Authentication bypasses prevented
- [ ] Authorization properly enforced
- [ ] File upload restrictions work
- [ ] Rate limiting functional
- [ ] Sensitive data not exposed in errors
- [ ] HTTPS enforced
- [ ] Security headers present

### Tools

- **OWASP ZAP**: Automated security scanner
- **Burp Suite**: Manual penetration testing
- **npm audit**: Check for vulnerable dependencies

```bash
# Run security audit
npm audit
pip-audit

# Check Docker images
docker scan construction-backend:latest
```

---

## Troubleshooting

### Common Issues

**Issue: Upload fails**
- Check file size limits
- Verify S3/MinIO credentials
- Check CORS configuration
- Verify network connectivity

**Issue: Analysis never completes**
- Check AI service logs
- Verify Redis connectivity
- Check job queue status
- Verify AI service has resources

**Issue: PDF download broken**
- Check presigned URL expiry
- Verify S3/MinIO access
- Check CORS on storage bucket
- Verify PDF generation succeeded

**Issue: Authentication fails**
- Verify JWT_SECRET matches
- Check token expiry
- Verify database connectivity
- Check CORS headers

---

## Reporting Issues

When reporting issues from E2E testing, include:

1. **Environment**: Production/Staging/Local
2. **Browser/Device**: Chrome on Windows, Safari on iPhone, etc.
3. **Steps to Reproduce**: Detailed steps
4. **Expected Behavior**: What should happen
5. **Actual Behavior**: What actually happened
6. **Screenshots/Videos**: Visual evidence
7. **Console Errors**: Browser console output
8. **Network Logs**: Failed requests
9. **Timestamps**: When the issue occurred
10. **User Account**: Test account used (if applicable)

---

## Success Criteria

The deployment is considered successful when:

- ✅ All E2E scenarios pass
- ✅ No critical bugs found
- ✅ Performance targets met
- ✅ Security tests pass
- ✅ User testing feedback positive
- ✅ All devices/browsers supported
- ✅ Documentation complete
- ✅ Monitoring alerts configured
- ✅ Backup/restore tested
- ✅ Rollback plan validated

---

## Next Steps After Validation

1. ✅ Document any issues found
2. ✅ Create tickets for non-critical bugs
3. ✅ Set up continuous E2E testing in CI/CD
4. ✅ Configure production monitoring
5. ✅ Train support team on common issues
6. ✅ Prepare user documentation
7. ✅ Plan for first production users
8. ✅ Set up feedback collection
