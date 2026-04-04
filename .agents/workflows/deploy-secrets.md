---
description: Run the production secret deployment protocol
---

When the user asks to deploy new secrets, rotate keys, or push environment configurations to production:

1. **Protocol**: You MUST execute `bash scripts/deploy_secrets.sh` to handle secret deployment natively.
2. **Security**: Stop and ask the user to provide the secret values securely. Do NOT echo or log the secret plaintexts back into the conversation or commit them to disk.
3. **Validation**: Validate that the deployment command succeeded before continuing feature work.
