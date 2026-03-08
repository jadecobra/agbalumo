# Global Squad Rules
These rules apply to ALL agents in the agbalumo squad.

## Rules
- Strict TDD. "Never write implementation code without a failing test case in listing_test.go."
- Go Performance. "Always use native Go concurrency (goroutines) for external API calls like Gemini or Supabase."
- "Always run './scripts/pre-commit.sh' to verify code quality. Do not use raw 'go test'."
- "Contact cards must include at least one valid method: WhatsApp, Email, or Phone."
- "Neighborhood data should prioritize Dallas-area locations for the MVP."
- "Strict Rule: 'Requests' (and all listings) MUST have a valid 'OwnerOrigin' (West African country). No exceptions."
- Cultural Tone. "All placeholder data must reflect West African naming conventions and Dallas-area geography."
- 10x Standard: "Validate first. Security always. Minimal changes only."
- **Traceability**: "Major changes must explicitly reference the specific Standard they satisfy (e.g., 'Ref: Standard 8.2')."
- **Pre-flight Check**: "Before executing, verify your plan against your specific `instructions` in this file."
