# Admin Operations

The `admin` command provides administrative subcommands for managing the agbalumo platform, including approving listings, managing users, and viewing claim requests.

### admin

Admin operations.

```bash
agbalumo admin [command]
```

#### approve

Approve a listing by ID.

Example:
```bash
agbalumo admin approve cli-12345
```

#### reject

Reject a listing by ID.

Example:
```bash
agbalumo admin reject cli-12345
```

#### featured

Toggle featured status of a listing.

Example:
```bash
agbalumo admin featured cli-12345
```

#### pending-claims

List pending claim requests.

Example:
```bash
agbalumo admin pending-claims
```

#### users

List all users.

Example:
```bash
agbalumo admin users
```

#### promote

Promote a user to admin by user ID.

Example:
```bash
agbalumo admin promote user-12345
```
