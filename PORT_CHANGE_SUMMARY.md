# Port Configuration Change Summary

## Changes Made

To improve user experience and match common expectations, we've updated the port configuration for development mode:

### Before (Old Configuration)
- **Frontend (UI)**: `http://localhost:3000`
- **Backend (API)**: `http://localhost:8080`
- **AI Service**: `http://localhost:8000`

### After (New Configuration)
- **Frontend (UI)**: `http://localhost:8080` ✨
- **Backend (API)**: `http://localhost:8081`
- **AI Service**: `http://localhost:8000`

## Why This Change?

Users commonly expect web applications to be available at `localhost:8080`. Previously, accessing `localhost:8080` would show the backend API JSON response instead of the UI, which was confusing.

## What Was Updated

### Configuration Files
- ✅ `docker-compose.yml` - Port mappings updated
- ✅ `.env.example` - Updated FRONTEND_URL
- ✅ `app/.env.example` - Updated EXPO_PUBLIC_API_URL
- ✅ `app/src/api/client.ts` - Updated default API URL

### Documentation Files
- ✅ `README.md` - Updated all port references
- ✅ `app/README.md` - Updated port information
- ✅ `backend/README.md` - Added clarification about docker-compose port
- ✅ `BID_EXPORT_DEMO.md` - Updated API endpoint examples
- ✅ `TESTING_SUMMARY.md` - Updated test URLs
- ✅ `DEPLOYMENT.md` - Updated deployment examples
- ✅ `REVISION_COMPARISON_DEMO.md` - Updated API examples
- ✅ `M6_IMPLEMENTATION_SUMMARY.md` - Updated examples
- ✅ `backend/API_TEST_GUIDE.md` - Updated all curl examples

### Test Files
- ✅ `app/__tests__/api.test.ts` - Updated expected base URL

## How to Use

After pulling these changes:

1. **Start the services:**
   ```bash
   make dev
   ```

2. **Access the application:**
   - **Web UI**: http://localhost:8080 (main interface)
   - **Backend API**: http://localhost:8081 (for direct API calls)
   - **AI Service**: http://localhost:8000 (AI endpoints)

3. **API Calls:**
   All API endpoint examples in documentation now use `http://localhost:8081`

## Backward Compatibility

This change only affects **development mode** (docker-compose.yml). 

**Production mode** (docker-compose.production.yml) remains unchanged:
- Frontend: Port 80 (configurable via `FRONTEND_PORT`)
- Backend: Port 8080 (configurable via `BACKEND_PORT`)
- AI Service: Port 8000 (configurable via `AI_SERVICE_PORT`)

## Migration Guide

If you have local `.env` files or scripts that reference the old ports:

1. **Update your local `.env` file:**
   ```bash
   # Old
   EXPO_PUBLIC_API_URL=http://localhost:8080
   
   # New
   EXPO_PUBLIC_API_URL=http://localhost:8081
   ```

2. **Update any custom scripts or bookmarks:**
   - Change UI bookmarks from `:3000` to `:8080`
   - Change API calls from `:8080` to `:8081`

## Troubleshooting

### "Can't connect to API"
- Verify backend is running on port 8081: `curl http://localhost:8081/health`
- Check `EXPO_PUBLIC_API_URL` is set to `http://localhost:8081`

### "Port already in use"
- Stop any services using port 8080 or 8081
- Run `docker compose down` and `make dev` again

### "404 Not Found on API calls"
- Ensure you're using `http://localhost:8081` for API endpoints
- Check docker-compose logs: `docker compose logs backend`

