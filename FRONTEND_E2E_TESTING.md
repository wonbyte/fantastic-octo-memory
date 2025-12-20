# Frontend E2E Testing Setup

This document describes the E2E testing infrastructure for the Construction Estimator frontend application.

## Overview

The frontend uses Playwright for end-to-end testing. Playwright is a modern testing framework that supports multiple browsers and provides excellent developer experience.

## Prerequisites

- Node.js 24 LTS
- npm 11.6+
- Playwright browsers installed

## Installation

E2E testing dependencies are already included in the root `package.json`. To install Playwright browsers:

```bash
npx playwright install
```

## Running Tests

From the root directory:

```bash
# Run all E2E tests
npm run test:e2e

# Run tests with UI mode (interactive)
npm run test:e2e:ui

# Run tests in headed mode (see browser)
npm run test:e2e:headed

# View test report
npm run test:e2e:report
```

## Test Structure

E2E tests are located in the `/e2e` directory at the root of the repository:

```
/e2e
  ├── basic.spec.ts           # Basic functionality tests
  └── user-journey.spec.ts    # Complete user flow tests
```

## Test Scenarios

### Authentication Tests (`basic.spec.ts`)
- Display login page
- Navigate to register page
- Keyboard navigation
- Accessibility labels

### Dark Mode Tests
- Theme toggle functionality
- Theme persistence
- Visual appearance changes

### Offline Mode Tests
- Offline indicator display
- Network status detection
- Offline behavior

### User Journey Tests (`user-journey.spec.ts`)
- Complete signup → login → project creation flow
- Responsive design on mobile viewport
- Responsive design on tablet viewport
- Page load performance

## Writing Tests

### Basic Test Structure

```typescript
import { test, expect } from '@playwright/test';

test.describe('Feature Name', () => {
  test('should do something', async ({ page }) => {
    await page.goto('/');
    
    // Test actions
    await page.click('button:has-text("Click Me")');
    
    // Assertions
    await expect(page.locator('text=Success')).toBeVisible();
  });
});
```

### Best Practices

1. **Use meaningful test descriptions**: Describe what the test validates
2. **Keep tests independent**: Each test should be able to run alone
3. **Use appropriate timeouts**: Network operations may take time
4. **Clean up after tests**: Reset state when needed
5. **Use accessibility selectors**: Prefer role-based and text selectors

### Accessibility Testing

Tests include accessibility checks:

```typescript
test('should have proper ARIA labels', async ({ page }) => {
  await page.goto('/');
  
  const buttons = await page.locator('button').all();
  for (const button of buttons) {
    const ariaLabel = await button.getAttribute('aria-label');
    const text = await button.textContent();
    
    expect(ariaLabel || text).toBeTruthy();
  }
});
```

### Responsive Design Testing

```typescript
test('should be usable on mobile', async ({ page }) => {
  await page.setViewportSize({ width: 375, height: 667 });
  await page.goto('/');
  
  // Verify no horizontal scroll
  const hasHorizontalScroll = await page.evaluate(() => {
    return document.documentElement.scrollWidth > 
           document.documentElement.clientWidth;
  });
  
  expect(hasHorizontalScroll).toBe(false);
});
```

## Configuration

The Playwright configuration is in `playwright.config.ts` at the root level:

```typescript
export default defineConfig({
  testDir: './e2e',
  use: {
    baseURL: 'http://localhost:8080',
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

## CI/CD Integration

E2E tests can be integrated into CI/CD pipelines:

```yaml
# Example GitHub Actions workflow
- name: Install Playwright Browsers
  run: npx playwright install --with-deps

- name: Run E2E tests
  run: npm run test:e2e
  env:
    BASE_URL: ${{ secrets.STAGING_URL }}

- name: Upload test results
  if: always()
  uses: actions/upload-artifact@v3
  with:
    name: playwright-report
    path: playwright-report/
```

## Debugging Tests

### Visual Debugging

Run tests in UI mode to debug visually:

```bash
npm run test:e2e:ui
```

### Headed Mode

Run tests with visible browser:

```bash
npm run test:e2e:headed
```

### Using Playwright Inspector

```bash
PWDEBUG=1 npx playwright test
```

### View Trace Files

After a test failure, view the trace:

```bash
npx playwright show-trace trace.zip
```

## Common Issues

### Port Already in Use

If port 3000 is already in use:

```bash
# Change BASE_URL in test or kill the process
lsof -ti:3000 | xargs kill -9
```

### Browser Not Installed

```bash
npx playwright install
```

### Tests Timing Out

Increase timeout in test:

```typescript
test('slow operation', async ({ page }) => {
  test.setTimeout(60000); // 60 seconds
  // ... test code
});
```

## Performance Testing

Tests include basic performance checks:

```typescript
test('should load within acceptable time', async ({ page }) => {
  const startTime = Date.now();
  await page.goto('/');
  await page.waitForLoadState('networkidle');
  const loadTime = Date.now() - startTime;
  
  expect(loadTime).toBeLessThan(5000); // 5 seconds
});
```

## Coverage

E2E tests cover:

- ✅ Authentication flow
- ✅ Dark mode toggle
- ✅ Offline mode indicator
- ✅ Accessibility features
- ✅ Responsive design
- ✅ Performance metrics
- ✅ Keyboard navigation

## Future Enhancements

- [ ] Visual regression testing
- [ ] API mocking for isolated tests
- [ ] Component-level E2E tests
- [ ] Cross-browser screenshot comparisons
- [ ] Extended user journey tests
- [ ] Load testing integration

## Resources

- [Playwright Documentation](https://playwright.dev/)
- [Best Practices](https://playwright.dev/docs/best-practices)
- [Debugging Guide](https://playwright.dev/docs/debug)
- [CI/CD Guide](https://playwright.dev/docs/ci)

## Support

For issues with E2E tests, please:

1. Check the test output and screenshots
2. Review Playwright documentation
3. Check for known issues in the repository
4. Open an issue with test failure details
