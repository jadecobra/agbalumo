---
name: "Deploy Secrets"
description: "Execute the production secret deployment protocol, key rotation, and hardened Fly.io configuration checks."
triggers:
  - "/deploy-secrets"
  - "deploy secrets"
  - "rotate keys"
  - "push secrets"
mutating: true
---

# Deploy Secrets Skill

Deploy production secrets, rotate encryption keys, and maintain container deployment security standards.

## Execution Protocol
1. Execute `bash scripts/deploy_secrets.sh` to handle secret deployment.
2. Never echo or log secret plaintext values to conversation logs or temporary disk files.
3. Prompt the user directly when secret values must be supplied.

## Production Hardening Invariants
When modifying deployment configurations or Dockerfiles, maintain four required security invariants.
1. Upgrade `libcrypto3` and `libssl3` in Dockerfile layers to address CVE-2024-13176.
2. Maintain `/usr/local/bin:/app` in the binary system PATH for Litestream and application binaries.
3. Start container processes as root, adjust ownership of `/data`, then drop privileges to `appuser` via `su-exec`.
4. Run non-critical backfill jobs in the background to prevent health check timeouts during container launch.

## Verification
1. Confirm the deployment command completed with exit status zero.
2. Inspect deployment logs using `fly logs` to confirm active replica synchronization.
