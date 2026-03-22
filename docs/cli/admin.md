# Admin CLI

Administrative operations.

```bash
agbalumo admin [command]
```

## Subcommands

### approve

Approve a pending listing.

```bash
agbalumo admin approve [listing-id]
```

**Example:**

```bash
agbalumo admin approve cli-1234567890
# Output: Listing approved: cli-1234567890
```

### reject

Reject a listing.

```bash
agbalumo admin reject [listing-id]
```

**Example:**

```bash
agbalumo admin reject cli-1234567890
# Output: Listing rejected: cli-1234567890
```

### featured

Toggle featured status of a listing.

```bash
agbalumo admin featured [listing-id]
```

**Example:**

```bash
agbalumo admin featured cli-1234567890
# Output: Listing featured: cli-1234567890
# Run again to unfeature:
# Output: Listing unfeatured: cli-1234567890
```

### pending-claims

List all pending claim requests.

```bash
agbalumo admin pending-claims
```

**Example:**

```bash
agbalumo admin pending-claims
# Output:
# Found 2 pending claims:
#
# [claim-123] User 'john@example.com' claiming 'Lagos Restaurant'
# [claim-456] User 'jane@example.com' claiming 'Accra Market'
```

### users

List all registered users.

```bash
agbalumo admin users
```

**Example:**

```bash
agbalumo admin users
# Output:
# Found 10 users:
#
# [user-123] john@example.com - user
# [admin-456] admin@example.com - admin
```

### promote

Promote a user to admin role.

```bash
agbalumo admin promote [user-id]
```

**Example:**

```bash
agbalumo admin promote user-123
# Output: User promoted to admin: user-123
```
