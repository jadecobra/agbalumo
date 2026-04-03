# agbalumo CLI: Listing Management

Manage business and community listings on the platform.

## Commands

### listing

Manage listings.

```bash
agbalumo listing [command]
```

#### Subcommands

##### create

Create a new listing.

```bash
agbalumo listing create -t "Title" -d "Description" -c "City" [flags]
```

**Flags:**

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--title` | `-t` | (required) | Listing title |
| `--type` | `-y` | Business | Type (Business, Service, Product, Food, Event, Job, Request) |
| `--origin` | `-o` | Nigeria | Owner origin/country |
| `--description` | `-d` | "" | Listing description |
| `--city` | `-c` | "" | City |
| `--address` | `-a` | "" | Address |
| `--email` | `-e` | "" | Contact email |
| `--phone` | `-p` | "" | Contact phone |
| `--whatsapp` | `-w` | "" | WhatsApp number |
| `--website` | `-s` | "" | Website URL |
| `--image-url` | `-i` | "" | Image URL |
| `--deadline` | | "" | Deadline (YYYY-MM-DD) |
| `--event-start` | | "" | Event start (YYYY-MM-DDTHH:MM) |
| `--event-end` | | "" | Event end (YYYY-MM-DDTHH:MM) |

**Example:**

```bash
agbalumo listing create --title "Lagos Deli" --type "Food" --city "Lagos"
```

##### list

List all listings.

```bash
agbalumo listing list [--json]
```

##### get

Get a listing by ID.

```bash
agbalumo listing get [id]
```

##### update

Update a listing.

```bash
agbalumo listing update [id] [flags]
```

##### delete

Delete a listing.

```bash
agbalumo listing delete [id]
```
##### backfill-cities

Backfill missing city data for listings using geocoding.

```bash
agbalumo listing backfill-cities
```
