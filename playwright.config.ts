import { defineConfig, devices } from '@playwright/test';

export default defineConfig({
  testDir: './tests/e2e',
  fullyParallel: true,
  forbidOnly: !!process.env.CI,
  retries: process.env.CI ? 2 : 0,
  workers: process.env.CI ? 2 : '50%',
  reporter: [['html', { open: 'never' }]],
  expect: {
    toHaveScreenshot: { maxDiffPixelRatio: 0.01 },
  },
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
      name: 'Desktop',
      use: { viewport: { width: 1440, height: 900 } },
    },
  ],
  webServer: {
    command: `sh -c "${process.env.AGBALUMO_TEST_SERVER_COMMAND || 'go run main.go serve'}"`,
    url: 'https://localhost:8443',
    reuseExistingServer: true,
    ignoreHTTPSErrors: true,
  },
});
