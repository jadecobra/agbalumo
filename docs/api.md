# agbalumo HTTP API Reference

REST API documentation for the agbalumo web application.

## Base URL

| Environment | URL |
|-------------|-----|
| Development | `http://localhost:8080` |
| Production | `https://agbalumo.com` |

## Authentication

The API uses session-based authentication with cookies.

| Method | Description |
|--------|-------------|
| **Session Cookie** | `session` cookie set after login |
| **CSRF Token** | Required for state-changing requests via `X-CSRF-Token` header |

### Endpoints

| Method | Path | Description |
|--------|------|-------------|
| GET | `/auth/dev` | Dev login (development only) |
| GET | `/auth/logout` | Clear session |
| GET | `/auth/google/login` | Initiate Google OAuth |
| GET | `/auth/google/callback` | Handle OAuth callback |
| GET | `/healthz` | Health check (returns 200 OK) |
| POST | `/api/metrics` | Ingest user interaction metrics |

## Public Endpoints

No authentication required.

| Method | Path | Description |
|--------|------|-------------|
| GET | `/` | Homepage with featured listings |
| GET | `/about` | About page |
| GET | `/listings/fragment` | HTMX partial for listings |
| GET | `/listings/:id` | Listing detail page |

### Query Parameters

**`/listings/fragment`**

| Parameter | Type | Description |
|-----------|------|-------------|
| `type` | string | Filter by category (category) |
| `q` | string | Search term (search) |
| `page` | integer | Page number for pagination |

## User Endpoints

Requires authentication (session cookie).

### Listings

| Method | Path | Description |
|--------|------|-------------|
| POST | `/listings` | Create new listing |
| GET | `/listings/:id/edit` | Get edit form |
| PUT | `/listings/:id` | Update listing |
| POST | `/listings/:id` | Update listing (alt) |
| DELETE | `/listings/:id` | Delete listing |
| GET | `/profile` | User profile page |
| POST | `/listings/:id/claim` | Claim listing |

### Feedback

| Method | Path | Description |
|--------|------|-------------|
| GET | `/feedback/modal` | Feedback form modal |
| POST | `/feedback` | Submit feedback |

### Request Body: Feedback

```json
{
  "type": "Bug|Feature|Question|Other",
  "content": "string (required)"
}
```

### Request Body: Create/Update Listing

```json
{
  "title": "string (required)",
  "type": "Business|Service|Product|Food|Event|Job|Request (required)",
  "owner_origin": "string (required)",
  "city": "string (required)",
  "description": "string",
  "address": "string",
  "hours_of_operation": "string",
  "contact_email": "string",
  "contact_phone": "string",
  "contact_whatsapp": "string",
  "website_url": "string",
  "deadline_date": "date (Request type)",
  "event_start": "datetime (Event type)",
  "event_end": "datetime (Event type)",
  "skills": "string (Job type)",
  "job_start_date": "datetime (Job type)",
  "job_apply_url": "string (Job type)",
  "company": "string (Job type)",
  "pay_range": "string (Job type)",
  "image": "file (multipart)",
  "remove_image": "boolean (optional, for updates)",
  "heat_level": "integer (0-5)",
  "regional_specialty": "string",
  "top_dish": "string"
}
```

## Admin Endpoints

Requires admin role and session authentication.

| Method | Path | Description |
|--------|------|-------------|
| GET | `/admin` | Dashboard |
| GET | `/admin/login` | Login form |
| POST | `/admin/login` | Login action |
| GET | `/admin/users` | List users |
| GET | `/admin/listings` | List all listings |
| GET | `/admin/listings/:id/row` | Return HTML row for listing |
| POST | `/admin/claims/:id/approve` | Approve claim request |
| POST | `/admin/claims/:id/reject` | Reject claim request |
| POST | `/admin/listings/:id/featured` | Toggle featured (`featured=true/false`) |
| POST | `/admin/listings/bulk` | Bulk action (approve|reject|delete) |
| GET | `/admin/listings/delete-confirm` | Delete confirmation (query param `id`) |
| POST | `/admin/listings/delete` | Delete listings (`admin_code` required) |
| POST | `/admin/upload` | Bulk CSV upload (`csv_file`) |
| GET | `/admin/listings/export` | Export all listings to CSV |
| POST | `/admin/categories` | Add custom category |
| GET | `/admin/modal/charts` | Admin charts modal fragment |
| GET | `/admin/modal/users` | Admin users modal fragment |
| GET | `/admin/modal/bulk` | Admin bulk upload modal fragment |
| GET | `/admin/modal/category` | Admin category management modal fragment |
| GET | `/admin/modal/moderation` | Admin moderation queue modal fragment |

### Admin Listing Filters (GET `/admin/listings`)

| Parameter | Type | Description |
|-----------|------|-------------|
| `category` | string | Filter by type |
| `sort` | string | Sort field (title, created_at, status) |
| `order` | string | Sort order (ASC, DESC) |
| `page` | integer | Page number |

### Bulk Action Request

```json
{
  "action": "approve|reject|delete",
  "selectedListings": ["id1", "id2", "id3"]
}
```

## Response Codes

| Code | Description |
|------|-------------|
| 200 | Success |
| 302 | Redirect (after form submission) |
| 400 | Validation error |
| 401 | Unauthorized |
| 404 | Not found |

## Rate Limits

| Endpoint | Limit |
|----------|-------|
| Default | 60 req/min |
| Admin login | 5 req/min |

## OpenAPI Specification

For complete schema definitions, see [openapi.yaml](openapi.yaml).

### Generate Client

```bash
# Using swagger-codegen
swagger-codegen generate -i docs/openapi.yaml -l go -o client

# Using openapi-generator
openapi-generator generate -i docs/openapi.yaml -g go -o client
```

### Validate Spec

```bash
# Using swagger-cli
swagger-cli validate docs/openapi.yaml
```
