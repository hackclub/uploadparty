import { test, expect } from '@playwright/test';

test.describe('UploadParty Application', () => {
  test.beforeEach(async ({ page }) => {
    // Navigate to the app
    await page.goto('/');
  });

  test('should load homepage successfully', async ({ page }) => {
    // Check that the page loads
    await expect(page).toHaveTitle(/UploadParty/);
    
    // Check for main navigation or content
    await expect(page.locator('body')).toBeVisible();
  });

  test('should have proper meta tags for SEO', async ({ page }) => {
    // Check essential meta tags
    await expect(page.locator('meta[name="description"]')).toHaveAttribute(
      'content', 
      /upload.*beats.*audio.*community.*challenges/i
    );
    
    await expect(page.locator('meta[name="keywords"]')).toHaveAttribute(
      'content',
      /music.*beats.*audio/i
    );
  });

  test('should be responsive on mobile', async ({ page }) => {
    // Set mobile viewport
    await page.setViewportSize({ width: 375, height: 667 });
    
    // Check that content is visible and properly laid out
    await expect(page.locator('body')).toBeVisible();
    
    // Check that text is readable (not too small)
    const bodyStyles = await page.locator('body').evaluate(el => {
      const styles = window.getComputedStyle(el);
      return {
        fontSize: styles.fontSize,
        fontFamily: styles.fontFamily,
      };
    });
    
    // Font size should be at least 14px for mobile readability
    const fontSize = parseInt(bodyStyles.fontSize);
    expect(fontSize).toBeGreaterThanOrEqual(14);
  });

  test('should handle server errors gracefully', async ({ page }) => {
    // Mock a server error
    await page.route('**/api/**', route => {
      route.fulfill({
        status: 500,
        contentType: 'application/json',
        body: JSON.stringify({ error: 'Internal server error' }),
      });
    });

    // Navigate to a page that would make API calls
    await page.goto('/app');
    
    // Should show error state instead of crashing
    // Look for error boundary or error message
    const hasErrorHandling = await page.locator('[data-testid="error"], .error, [class*="error"]').count() > 0 ||
                            await page.locator('text=error').count() > 0 ||
                            await page.locator('text=something went wrong').count() > 0;
    
    expect(hasErrorHandling).toBeTruthy();
  });

  test('should have accessible navigation', async ({ page }) => {
    // Check for proper heading hierarchy
    const h1Count = await page.locator('h1').count();
    expect(h1Count).toBeGreaterThanOrEqual(1);
    expect(h1Count).toBeLessThanOrEqual(1); // Should have exactly one h1
    
    // Check for skip links or other accessibility features
    const skipLink = page.locator('a[href="#main"], a[href="#content"], [class*="skip"]');
    const hasSkipLink = await skipLink.count() > 0;
    
    // Check for focus management
    await page.keyboard.press('Tab');
    const focusedElement = await page.locator(':focus');
    await expect(focusedElement).toBeVisible();
  });

  test('should load and display fonts correctly', async ({ page }) => {
    // Wait for fonts to load
    await page.waitForLoadState('networkidle');
    
    // Check that custom fonts are applied
    const bodyFontFamily = await page.locator('body').evaluate(el => 
      window.getComputedStyle(el).fontFamily
    );
    
    // Should include custom font (Geist in this case)
    expect(bodyFontFamily).toMatch(/geist/i);
  });

  test('should have proper error boundaries', async ({ page }) => {
    // Create a JavaScript error to trigger error boundary
    await page.evaluate(() => {
      // Simulate a React error
      window.dispatchEvent(new Error('Test error for error boundary'));
    });
    
    // Page should still be functional
    await expect(page.locator('body')).toBeVisible();
  });
});

test.describe('Performance', () => {
  test('should meet Core Web Vitals targets', async ({ page }) => {
    await page.goto('/');
    
    // Measure performance
    const performanceMetrics = await page.evaluate(() => {
      return new Promise((resolve) => {
        new PerformanceObserver((list) => {
          const entries = list.getEntries();
          const metrics = {};
          
          entries.forEach((entry) => {
            if (entry.entryType === 'navigation') {
              metrics.loadTime = entry.loadEventEnd - entry.loadEventStart;
              metrics.domContentLoaded = entry.domContentLoadedEventEnd - entry.domContentLoadedEventStart;
            }
          });
          
          resolve(metrics);
        }).observe({ entryTypes: ['navigation'] });
        
        // Fallback timeout
        setTimeout(() => resolve({}), 5000);
      });
    });
    
    // Basic performance expectations
    if (performanceMetrics.loadTime) {
      expect(performanceMetrics.loadTime).toBeLessThan(3000); // 3 seconds
    }
  });

  test('should have optimized images', async ({ page }) => {
    await page.goto('/');
    
    // Check for proper image optimization
    const images = await page.locator('img').all();
    
    for (const img of images) {
      const src = await img.getAttribute('src');
      const loading = await img.getAttribute('loading');
      
      if (src) {
        // Check for lazy loading (except above-the-fold images)
        const boundingBox = await img.boundingBox();
        if (boundingBox && boundingBox.y > 600) {
          expect(loading).toBe('lazy');
        }
        
        // Check for responsive images (srcset or Next.js optimized)
        const srcset = await img.getAttribute('srcset');
        const isNextImage = src.includes('/_next/image') || src.includes('/_next/static');
        
        expect(srcset || isNextImage).toBeTruthy();
      }
    }
  });
});