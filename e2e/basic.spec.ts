import { test, expect } from '@playwright/test';

test.describe('Authentication Flow', () => {
  test('should display login page', async ({ page }) => {
    await page.goto('/');
    
    // Check for login elements
    await expect(page.locator('text=Login')).toBeVisible({ timeout: 10000 });
  });

  test('should navigate to register page', async ({ page }) => {
    await page.goto('/');
    
    // Look for sign up or register link
    const signupLink = page.locator('text=/sign.*up|register/i').first();
    if (await signupLink.isVisible()) {
      await signupLink.click();
      await expect(page).toHaveURL(/register|signup/);
    }
  });
});

test.describe('Accessibility', () => {
  test('should have proper ARIA labels on buttons', async ({ page }) => {
    await page.goto('/');
    
    // Check for buttons with accessibility labels
    const buttons = await page.locator('button').all();
    for (const button of buttons) {
      const ariaLabel = await button.getAttribute('aria-label');
      const text = await button.textContent();
      
      // Button should have either aria-label or text content
      expect(ariaLabel || text).toBeTruthy();
    }
  });

  test('should support keyboard navigation', async ({ page }) => {
    await page.goto('/');
    
    // Tab through interactive elements
    await page.keyboard.press('Tab');
    const focusedElement = await page.evaluate(() => document.activeElement?.tagName);
    
    // Should be able to focus on interactive elements
    expect(['BUTTON', 'INPUT', 'A', 'TEXTAREA']).toContain(focusedElement);
  });
});

test.describe('Dark Mode', () => {
  test('should be able to toggle dark mode', async ({ page }) => {
    await page.goto('/');
    
    // Look for theme toggle if accessible without login
    const themeToggle = page.locator('button:has-text("Dark"), button:has-text("Light")').first();
    
    if (await themeToggle.isVisible({ timeout: 5000 }).catch(() => false)) {
      // Get initial background color
      const initialBg = await page.evaluate(() => 
        window.getComputedStyle(document.body).backgroundColor
      );
      
      await themeToggle.click();
      await page.waitForTimeout(500);
      
      // Check if background changed
      const newBg = await page.evaluate(() => 
        window.getComputedStyle(document.body).backgroundColor
      );
      
      expect(initialBg).not.toBe(newBg);
    }
  });
});

test.describe('Offline Mode', () => {
  test('should display offline indicator when offline', async ({ page, context }) => {
    await page.goto('/');
    
    // Go offline
    await context.setOffline(true);
    await page.waitForTimeout(1000);
    
    // Look for offline indicator
    const offlineIndicator = page.locator('text=/offline|no.*connection/i');
    await expect(offlineIndicator).toBeVisible({ timeout: 5000 });
    
    // Go back online
    await context.setOffline(false);
  });
});
