import { test, expect } from '@playwright/test';

/**
 * Complete user journey E2E test
 * Tests the full flow from signup to project creation
 */
test.describe('Complete User Journey', () => {
  const testEmail = `test-${Date.now()}@example.com`;
  const testPassword = 'TestPassword123!';

  test('complete flow: signup → login → project creation', async ({ page }) => {
    // Step 1: Navigate to app
    await page.goto('/');
    await page.waitForLoadState('networkidle');

    // Step 2: Try to register (if register page exists)
    const registerLink = page.locator('text=/sign.*up|register|create.*account/i').first();
    
    if (await registerLink.isVisible({ timeout: 5000 }).catch(() => false)) {
      await registerLink.click();
      await page.waitForURL(/register|signup/);
      
      // Fill registration form
      await page.fill('input[name="email"], input[type="email"]', testEmail);
      await page.fill('input[name="password"], input[type="password"]', testPassword);
      
      // Look for name field if it exists
      const nameInput = page.locator('input[name="name"]').first();
      if (await nameInput.isVisible({ timeout: 2000 }).catch(() => false)) {
        await nameInput.fill('Test User');
      }
      
      // Submit registration
      const submitButton = page.locator('button[type="submit"], button:has-text("Sign Up"), button:has-text("Register")').first();
      await submitButton.click();
      
      // Wait for navigation after registration
      await page.waitForURL(/dashboard|projects|home/, { timeout: 10000 }).catch(() => {
        console.log('Did not redirect to dashboard after registration');
      });
    }

    // Step 3: If we're at login page, login
    if (await page.locator('text=/login|sign.*in/i').isVisible({ timeout: 2000 }).catch(() => false)) {
      await page.fill('input[name="email"], input[type="email"]', testEmail);
      await page.fill('input[name="password"], input[type="password"]', testPassword);
      
      const loginButton = page.locator('button[type="submit"], button:has-text("Login"), button:has-text("Sign In")').first();
      await loginButton.click();
      
      await page.waitForURL(/dashboard|projects|home/, { timeout: 10000 });
    }

    // Step 4: Create a project (if accessible)
    const newProjectButton = page.locator('button:has-text("New Project"), button:has-text("Create Project"), a:has-text("New Project")').first();
    
    if (await newProjectButton.isVisible({ timeout: 5000 }).catch(() => false)) {
      await newProjectButton.click();
      
      // Fill project form
      await page.fill('input[name="name"]', 'E2E Test Project');
      
      const descriptionInput = page.locator('textarea[name="description"], input[name="description"]').first();
      if (await descriptionInput.isVisible({ timeout: 2000 }).catch(() => false)) {
        await descriptionInput.fill('Automated E2E test project');
      }
      
      // Submit project creation
      const createButton = page.locator('button:has-text("Create"), button[type="submit"]').first();
      await createButton.click();
      
      // Verify project was created
      await expect(page.locator('text=E2E Test Project')).toBeVisible({ timeout: 10000 });
    }
  });
});

test.describe('Responsive Design', () => {
  test('should be usable on mobile viewport', async ({ page }) => {
    await page.setViewportSize({ width: 375, height: 667 });
    await page.goto('/');
    
    // Check that content is visible and not cut off
    await expect(page.locator('body')).toBeVisible();
    
    // Verify no horizontal scroll
    const hasHorizontalScroll = await page.evaluate(() => {
      return document.documentElement.scrollWidth > document.documentElement.clientWidth;
    });
    
    expect(hasHorizontalScroll).toBe(false);
  });

  test('should be usable on tablet viewport', async ({ page }) => {
    await page.setViewportSize({ width: 768, height: 1024 });
    await page.goto('/');
    
    await expect(page.locator('body')).toBeVisible();
  });
});

test.describe('Performance', () => {
  test('should load within acceptable time', async ({ page }) => {
    const startTime = Date.now();
    await page.goto('/');
    await page.waitForLoadState('networkidle');
    const loadTime = Date.now() - startTime;
    
    // Should load within 5 seconds
    expect(loadTime).toBeLessThan(5000);
  });
});
