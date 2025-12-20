# Testing Instructions for Port Change

## Quick Verification

After starting the services with `make dev`, verify the following:

### 1. Check Services are Running

```bash
# Check all containers are up
docker ps | grep construction

# You should see:
# - construction-backend
# - construction-frontend
# - construction-postgres
# - construction-redis
# - construction-minio
# - construction-ai-service
```

### 2. Verify Port Accessibility

```bash
# Test Backend API (should return JSON)
curl http://localhost:8081/
# Expected: {"message":"Construction Estimation & Bidding Automation API","version":"1.0.0"}

# Test Backend Health
curl http://localhost:8081/health
# Expected: {"status":"ok",...}

# Test AI Service
curl http://localhost:8000/health
# Expected: {"status":"healthy"}

# Test Frontend (should return HTML)
curl -I http://localhost:8080/
# Expected: HTTP/1.1 200 OK
```

### 3. Browser Testing

1. **Open the UI:**
   - Navigate to: http://localhost:8080
   - Expected: React/Expo app UI loads (not JSON)

2. **Verify API connectivity:**
   - If the app shows connection errors, check browser console
   - API calls should go to `http://localhost:8081`

3. **Test basic functionality:**
   - Try to sign up or log in
   - Check if API requests appear in Network tab
   - Verify requests go to correct port (8081)

### 4. E2E Tests

```bash
# Run Playwright E2E tests
npm run test:e2e

# Expected: Tests should connect to http://localhost:8080
```

## Troubleshooting

### Port Conflict Errors

If you get "port already in use" errors:

```bash
# Stop all Docker containers
docker compose down

# Check what's using the ports
lsof -ti:8080 | xargs kill -9
lsof -ti:8081 | xargs kill -9

# Restart
make dev
```

### Frontend Shows 404

If the frontend doesn't load:

```bash
# Check frontend container logs
docker logs construction-frontend -f

# Verify the container is running
docker ps | grep frontend

# Restart if needed
docker compose restart frontend
```

### API Connection Errors

If the frontend can't reach the API:

```bash
# Check environment variable in container
docker exec construction-frontend env | grep EXPO_PUBLIC_API_URL
# Should show: EXPO_PUBLIC_API_URL=http://localhost:8081

# Check backend is reachable
docker exec construction-frontend curl http://backend:8080/health

# Check from host
curl http://localhost:8081/health
```

## Success Criteria

✅ All tests pass:
- [ ] Backend responds on port 8081
- [ ] Frontend loads on port 8080
- [ ] API calls from frontend go to port 8081
- [ ] E2E tests pass
- [ ] No console errors in browser
- [ ] Health checks return OK

## Rollback

If you need to revert the changes:

```bash
# Checkout the previous commit
git checkout HEAD~2

# Restart services
docker compose down
make dev
```

The old configuration will be:
- Frontend: http://localhost:3000
- Backend: http://localhost:8080
