# Agbalumo CLI Documentation

Command-line interface for managing the Agbalumo directory platform.

## Commands

### listing

Manage business listings.

```bash
agbalumo listing [command]
```

#### Subcommands

##### create

Create a new listing.

```bash
agbalumo listing create -t "Business Title" -d "Description" -c "City" [flags]
```

**Flags:**

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--title` | `-t` | (required) | Listing title |
| `--type` | `-y` | Business | Listing type (Business, Service, Product, Food, Event, Job, Request) |
| `--origin` | `-o` | Nigeria | Owner origin/country |
| `--description` | `-d` | "" | Listing description |
| `--city` | `-c` | "" | City |
| `--address` | `-a` | "" | Address |
| `--email` | `-e` | "" | Contact email |
| `--phone` | `-p` | "" | Contact phone |
| `--whatsapp` | `-w` | "" | WhatsApp number |
| `--website` | `-s` | "" | Website URL |
| `--owner-id` | | "" | Owner user ID |

**Example:**

```bash
agbalumo listing create \
  --title "Lagos Restaurant" \
  --type "Business" \
  --origin "Nigeria" \
  --description "Authentic Nigerian cuisine" \
  --city "Lagos" \
  --address "123 Main Street" \
  --email "info@lagosrestaurant.com" \
  --phone "+2341234567890"
```

##### list

List all listings.

```bash
agbalumo listing list
```

**Example:**

```bash
agbalumo listing list
# Output:
# Found 50 listings:
#
# [abc12345] Lagos Restaurant - Business (Lagos) [Approved]
# [def67890] Accra Market - Business (Accra) [Approved]
# ...
```

##### get

Get a listing by ID.

```bash
agbalumo listing get [id]
```

**Example:**

```bash
agbalumo listing get cli-1234567890
```

##### update

Update a listing.

```bash
agbalumo listing update [id] [flags]
```

**Flags:**

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--title` | `-t` | "" | New title |
| `--description` | `-d` | "" | New description |
| `--city` | `-c` | "" | New city |
| `--address` | `-a` | "" | New address |
| `--email` | `-e` | "" | New email |
| `--phone` | `-p` | "" | New phone |
| `--whatsapp` | `-w` | "" | New WhatsApp |
| `--website` | `-s` | "" | New website |

**Example:**

```bash
agbalumo listing update cli-1234567890 \
  --title "Updated Title" \
  --city "Abuja"
```

##### delete

Delete a listing.

```bash
agbalumo listing delete [id]
```

**Example:**

```bash
agbalumo listing delete cli-1234567890
```

---

### admin

Administrative operations.

```bash
agbalumo admin [command]
```

#### Subcommands

##### approve

Approve a pending listing.

```bash
agbalumo admin approve [listing-id]
```

**Example:**

```bash
agbalumo admin approve cli-1234567890
# Output: Listing approved: cli-1234567890
```

##### reject

Reject a listing.

```bash
agbalumo admin reject [listing-id]
```

**Example:**

```bash
agbalumo admin reject cli-1234567890
# Output: Listing rejected: cli-1234567890
```

##### featured

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

##### pending

List all pending listings awaiting approval.

```bash
agbalumo admin pending
```

**Example:**

```bash
agbalumo admin pending
# Output:
# Found 5 pending listings:
#
# [abc12345] New Restaurant - Business (Lagos) [Pending]
# [def67890] Service Company - Service (Accra) [Pending]
# ...
```

##### users

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

##### promote

Promote a user to admin role.

```bash
agbalumo admin promote [user-id]
```

**Example:**

```bash
agbalumo admin promote user-123
# Output: User promoted to admin: user-123
```

---

### seed

Seed the database with initial data.

```bash
agbalumo seed [database-path]
```

**Example:**

```bash
agbalumo seed
agbalumo seed custom.db
```

---

### serve

Start the Agbalumo web server.

```bash
agbalumo serve [flags]
```

**Flags:**

| Flag | Default | Description |
|------|---------|-------------|
| `--env` | development | Environment (development, production) |
| `--port` | 8080 | Server port |
| `--cert` | certs/cert.pem | TLS certificate file |
| `--key` | certs/key.pem | TLS key file |

**Example:**

```bash
# Development server
agbalumo serve

# Production server
agbalumo serve --env production --port 443 --cert /path/to/cert.pem --key /path/to/key.pem
```

---

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `DATABASE_URL` | agbalumo.db | SQLite database path |
| `SESSION_SECRET` | dev-secret-key | Session encryption secret |
| `ADMIN_CODE` | (generated) | Admin access code |
| `UPLOAD_DIR` | ./uploads | Image upload directory |
| `GOOGLE_CLIENT_ID` | - | Google OAuth client ID |
| `GOOGLE_CLIENT_SECRET` | - | Google OAuth client secret |
| `RATE_LIMIT_RATE` | 60 | Requests per minute |
| `RATE_LIMIT_BURST` | 10 | Burst limit |

---

## Quick Reference

```bash
# List all listings
agbalumo listing list

# Create a listing
agbalumo listing create -t "My Business" -c "Lagos"

# Get listing details
agbalumo listing get <id>

# Update a listing
agbalumo listing update <id> -t "New Title"

# Delete a listing
agbalumo listing delete <id>

# Admin: Approve listing
agbalumo admin approve <id>

# Admin: List pending
agbalumo admin pending

# Admin: List users
agbalumo admin users

# Admin: Promote user
agbalumo admin promote <user-id>
```
