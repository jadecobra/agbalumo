---
description: Run the production secret deployment protocol
---

When the user asks to deploy new secrets, rotate keys, or push environment configurations to production:

1. **Protocol**: You MUST execute `bash scripts/deploy_secrets.sh` to handle secret deployment natively.
2. **Security**: Stop and ask the user to provide the secret values securely. Do NOT echo or log the secret plaintexts back into the conversation or commit them to disk.
3. **Invariants**: When deploying or modifying the deployment pipeline, you MUST preserve the following hardened patterns:
    - **Vulnerability Patching**: Always upgrade `libcrypto3` and `libssl3` in the Dockerfile to resolve CVE-2024-13176.
    - **Binary PATH**: Ensure `/usr/local/bin:/app` is in the `PATH` for Litestream and App binaries.
    - **Privilege Management**: Start as `root`, `chown` the `/data` volume, and use `su-exec` to drop to `appuser`.
    - **Background Backfill**: Non-critical tasks like `city backfill` MUST be run in the background (`&`) to prevent health check timeouts.
4. **Validation**: Validate that the deployment command succeeded and monitor logs with `fly logs` to ensure `replica sync` is active.
