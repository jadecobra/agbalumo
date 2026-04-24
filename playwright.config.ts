import { defineConfig, devices } from '@playwright/test';

export default defineConfig({
  testDir: './tests/e2e',
  fullyParallel: true,
  forbidOnly: !!process.env.CI,
  retries: process.env.CI ? 2 : 0,
  workers: process.env.CI ? 1 : undefined,
  reporter: 'html',
  use: {
    baseURL: process.env.BASE_URL || 'https://localhost:8443',
    trace: 'on-first-retry',
    ignoreHTTPSErrors: true,
  },
  projects: [
    {
      name: 'Mobile',
      use: { viewport: { width: 375, height: 812 } },
    },
    {
      name: 'Tablet',
      use: { viewport: { width: 768, height: 1024 } },
    },
    {
      name: 'Desktop',
      use: { viewport: { width: 1440, height: 900 } },
    },
    {
      name: 'Wide',
      use: { viewport: { width: 1920, height: 1080 } },
    },
  ],
});
