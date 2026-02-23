# Agbalumo HTTP API Reference

REST API documentation for the Agbalumo web application.

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
| `category` | string | Filter by category |
| `search` | string | Search term |
| `city` | string | Filter by city |
| `origin` | string | Filter by country of origin |

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
  "deadline": "datetime (Request type)",
  "event_start": "datetime (Event type)",
  "event_end": "datetime (Event type)",
  "skills": "string (Job type)",
  "job_start_date": "datetime (Job type)",
  "job_apply_url": "string (Job type)",
  "company": "string (Job type)",
  "pay_range": "string (Job type)",
  "image": "file (multipart)"
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
| POST | `/admin/listings/:id/approve` | Approve listing |
| POST | `/admin/listings/:id/reject` | Reject listing |
| POST | `/admin/listings/:id/featured` | Toggle featured |
| POST | `/admin/listings/bulk` | Bulk action |
| GET | `/admin/listings/delete-confirm` | Delete confirmation |
| POST | `/admin/listings/delete` | Delete listing |
| POST | `/admin/upload` | Bulk CSV upload |

### Bulk Action Request

```json
{
  "action": "approve|reject|delete",
  "ids": ["id1", "id2", "id3"]
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
