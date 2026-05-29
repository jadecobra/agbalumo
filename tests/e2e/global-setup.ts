/**
 * Playwright Global Setup — clears stale claim state from prior E2E runs.
 * Runs ONCE before all projects, ensuring clean DB for owner-claim tests.
 */
import { execSync } from 'child_process';
import * as path from 'path';

export default function globalSetup() {
  const dbPath = path.resolve(process.cwd(), '.tester', 'data', 'agbalumo.db');
  console.log(`[globalSetup] Cleaning DB at: ${dbPath}`);
  try {
    execSync(
      `sqlite3 "${dbPath}" "DELETE FROM claim_requests; UPDATE listings SET owner_id = '' WHERE owner_id != '';"`,
      { timeout: 5000, stdio: 'inherit' }
    );
    console.log('[globalSetup] ✓ Cleared claim_requests and listing owners');
  } catch (err) {
    // Non-fatal: table may not exist on first run
    console.warn('[globalSetup] ⚠ DB cleanup skipped:', (err as Error).message);
  }
}
