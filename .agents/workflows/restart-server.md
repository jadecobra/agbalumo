---
description: restart the application server
---
This workflow stops any running instance of the application server, builds the CSS and server binary, and starts a new instance on port 8443 with HTTPS, followed by verifying availability.

1. Restart the server and verify it is running using the dedicated verification script.
// turbo
2. Run `bash scripts/verify_restart.sh`
