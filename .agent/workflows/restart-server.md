---
description: restart the application server
---
This workflow stops any running instance of the application server and starts a new one on port 8443 with HTTPS, followed by verifying availability.

1. Stop any existing server process running on port 8443.
// turbo
2. Run `lsof -t -i :8443 | xargs kill -9 2>/dev/null || true`

3. Start the server in the background. Note: this command should be sent to the background asynchronously (e.g., set `WaitMsBeforeAsync` to around `2000`).
4. Run `export PATH=$PATH:/opt/homebrew/bin && go run main.go serve`

5. Verify that the server started successfully and is listening.
// turbo
6. Run `curl -k -I https://localhost:8443`
