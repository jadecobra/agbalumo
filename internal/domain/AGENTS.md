# Domain Layer (internal/domain)

Core types, interfaces, and business rules. Zero external dependencies.

## Constraints
- **No imports** from `internal/handler`, `internal/repository`, `internal/service`, or any infrastructure package.
- **Interfaces**: All store/service contracts are defined here (`ListingStore`, `UserStore`, `ListingService`, etc.).
- **CQRS**: `ListingReader` (reads) and `ListingWriter` (writes) compose into `ListingStore`. Prefer focused interfaces in new code.
- **Validation**: Business validation rules live in `listing_validation.go`, not in handlers or repositories.
- **Constants**: All env keys, field names, and status enums are centralized in `constants.go`.
