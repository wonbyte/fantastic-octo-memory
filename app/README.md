# Construction Estimator - Frontend App

React Native + Expo application for the Construction Estimation & Bidding Automation platform.

## 🚀 Tech Stack

- **Expo SDK**: 54.0+
- **React Native**: 0.82
- **React**: 19.0+
- **TypeScript**: 5.9+
- **Node.js**: 24 LTS
- **Expo Router**: 6.0+ (File-based routing)
- **React Query**: 5.60+ (Data fetching and caching)
- **Axios**: HTTP client
- **React Native Reanimated**: Animations
- **React Native Gesture Handler**: Gestures

## 📁 Project Structure

```
/app
  ├── app/                    # Expo Router pages
  │   ├── (auth)/            # Auth flow (login, register)
  │   ├── (main)/            # Main app flow (projects, blueprints)
  │   ├── _layout.tsx        # Root layout with providers
  │   └── index.tsx          # Entry point with auth redirect
  ├── src/
  │   ├── api/               # API client and endpoint functions
  │   ├── components/        # Reusable UI components
  │   ├── contexts/          # React contexts (Auth, Theme)
  │   ├── hooks/             # Custom hooks with React Query
  │   ├── types/             # TypeScript type definitions
  │   └── utils/             # Utility functions and constants
  ├── __tests__/             # Test files
  └── assets/                # Images, fonts, etc.
```

## 🛠️ Development

### Prerequisites

- Node.js 22 LTS
- npm 11.6+
- Expo CLI

### Installation

```bash
npm install
```

### Running the App

```bash
# Start Expo development server
npm start

# Run on Android
npm run android

# Run on iOS
npm run ios

# Run on Web
npm run web
```

### Environment Variables

Copy `.env.example` to `.env` and configure:

```env
EXPO_PUBLIC_API_URL=http://localhost:8080
EXPO_PUBLIC_ENV=development
```

For Docker, use the backend container name:

```env
EXPO_PUBLIC_API_URL=http://backend:8080
```

## 🧪 Testing

```bash
# Run all tests
npm test

# Run tests in watch mode
npm test -- --watch

# Type checking
npm run type-check

# Linting
npm run lint

# E2E tests (from root directory)
cd .. && npm run test:e2e

# E2E tests with UI
cd .. && npm run test:e2e:ui
```

## 🏗️ Features Implemented

### Authentication
- ✅ Login screen with email/password
- ✅ Register screen
- ✅ Auth context with token management
- ✅ Protected routes with automatic redirect

### Projects
- ✅ Project list with pull-to-refresh
- ✅ Create new project
- ✅ Project detail view
- ✅ Empty states and loading states

### Blueprints
- ✅ Blueprint upload with file picker
- ✅ S3 pre-signed URL upload flow
- ✅ Upload progress tracking
- ✅ Blueprint detail view
- ✅ Blueprint list in project

### Analysis
- ✅ Trigger blueprint analysis
- ✅ Job status polling with React Query
- ✅ Analysis results display
- ✅ Real-time status updates
- ✅ Progress indicators

### UI Components
- ✅ Button (primary, secondary, danger variants)
- ✅ Card
- ✅ Input with validation
- ✅ Loading spinner
- ✅ Error state with retry

## 📱 Navigation Structure

```
├── / (index)
│   ├── Redirect to login or main based on auth
├── (auth)
│   ├── login
│   └── register
└── (main)
    └── projects (Tab)
        ├── index (List)
        └── [id] (Detail)
            └── blueprints
                ├── upload
                ├── [blueprintId] (Detail)
                └── [blueprintId]/analysis
```

## 🐳 Docker Support

The app is configured to run in Docker with:

- Hot reload enabled
- Metro bundler configured for container environment
- CORS headers for development
- API calls to backend container

### Running in Docker

```bash
# From the root of the monorepo
docker-compose up frontend

# Or use make
make dev
```

Access the app at:
- Web: http://localhost:3000
- Metro: http://localhost:19000

## 🔧 Configuration Files

- `app.json` - Expo app configuration
- `babel.config.js` - Babel with Reanimated plugin
- `metro.config.js` - Metro bundler config for Docker
- `tsconfig.json` - TypeScript configuration
- `jest.config.js` - Jest test configuration
- `eslint.config.js` - ESLint configuration

## 📝 Code Quality

### Type Safety
All code is written in TypeScript with strict type checking enabled.

### Linting
ESLint is configured with React and TypeScript rules.

### Testing
Jest and React Testing Library are used for unit and component tests.

## 🚧 Future Enhancements

- [x] Offline mode with local caching
- [x] Push notifications for job completion
- [x] Dark mode support
- [x] Accessibility improvements
- [x] Performance optimizations
- [x] E2E tests with Playwright
- [ ] Storybook for component documentation

## ✨ Recent Additions (v1.1.0)

### Dark Mode Support
- Light, dark, and auto themes
- Persistent theme preference
- System theme detection
- Theme toggle component

### Offline Mode & Caching
- React Query persistent cache
- Network status detection
- Offline indicator UI
- Data persistence with AsyncStorage

### Push Notifications
- Local notification support
- Job completion notifications
- Error notifications
- Blueprint analysis complete notifications
- Configurable notification permissions

### Accessibility Improvements (WCAG Compliance)
- Accessibility labels on all interactive elements
- Screen reader support
- Focus management
- Proper semantic HTML/ARIA roles
- Keyboard navigation support

### E2E Testing
- Playwright test infrastructure
- Basic E2E test suite
- User journey tests
- Responsive design tests
- Performance tests

## 📄 API Integration

The app integrates with the Go backend API:

- **Authentication**: POST `/auth/login`, POST `/auth/register`
- **Projects**: GET/POST/PUT/DELETE `/projects`
- **Blueprints**: GET/POST `/projects/{id}/blueprints`
- **Upload**: POST `/projects/{id}/blueprints/upload-url`
- **Analysis**: POST `/blueprints/{id}/analyze`
- **Jobs**: GET `/jobs/{id}`

See the backend API documentation for full endpoint details.

## 🐛 Troubleshooting

### Metro bundler not starting
```bash
# Clear cache and restart
npm start -- --reset-cache
```

### Type errors
```bash
# Rebuild TypeScript
npm run type-check
```

### Docker networking issues
Ensure the API URL is set correctly:
- Local dev: `http://localhost:8080`
- Docker: `http://backend:8080`

## 📞 Support

For issues or questions, please open a GitHub issue.
