## spec.md: Agbalumo MVP (Golang Edition)

### 1. Project Overview

**Agbalumo** is a high-performance directory and request platform for the West African diaspora.

* **Goal:** users can find West African businesses, services, products, food, jobs, and events that meet their needs.
* **Architecture:** **Go (Golang)** for the backend, utilizing `Echo` for high-concurrency routing, `HTMX` for dynamic frontend interactions, and `SQLite` for local data persistence.
* **TDD Strategy:** Leverage Go’s built-in `testing` package with `testify` for assertions. RED GREEN REFACTOR cycle.

### 2. Data Models (Go Structs)

The system use strict types to ensure data integrity.

```go
type Category string

const (
    Business Category = "Business"
    Service  Category = "Service"
    Product  Category = "Product"
    Job      Category = "Job"
    Request  Category = "Request"
    Food     Category = "Food"
    Event    Category = "Event"
)

type Listing struct {
    ID              string    `json:"id"`
    OwnerID         string    `json:"owner_id"`         // Link to User.ID
    OwnerOrigin     string    `json:"owner_origin"`     // Required: Country of Origin
    Type            Category  `json:"type"`
    Title           string    `json:"title"`
    Description     string    `json:"description"`
    
    // Location & Contact
    City            string    `json:"city"`
    Address         string    `json:"address"`
    ContactEmail    string    `json:"contact_email"`
    ContactPhone    string    `json:"contact_phone"`
    ContactWhatsApp string    `json:"contact_whatsapp"`
    WebsiteURL      string    `json:"website_url"`
    
    // Media
    ImageURL        string    `json:"image_url"`
    
    // Metadata
    CreatedAt       time.Time `json:"created_at"`
    IsActive        bool      `json:"is_active"`
    
    // Type Specific
    Deadline        time.Time `json:"deadline"`       // Request
    EventStart      time.Time `json:"event_start"`    // Event
    EventEnd        time.Time `json:"event_end"`      // Event
    Skills          string    `json:"skills"`         // Job
    JobStartDate    time.Time `json:"job_start_date"` // Job
    Company         string    `json:"company"`        // Job
    PayRange        string    `json:"pay_range"`      // Job
    JobApplyURL     string    `json:"job_apply_url"`  // Job
}
```

### 3. functional Requirements & Constraints

* **Concurrency:** Use Goroutines for non-blocking Gemini AI moderation tasks.
* **Authentication:** Google OAuth2 for user authentication. **Posting**, editing, and deleting listings requires authentication.
* **Deadline Validation:** 
    * `Request` types must have a `Deadline` within 90 days.
    * `Event` types must have valid start/end times.
    * `Job` types must have a valid start date.
* **Contact Integrity:** Every listing must contain at least one valid communication method.

### 4. TDD Specification

#### 4.1 Unit Tests (`listing_test.go`)

* `TestValidateDeadline`: Verify deadline constraints.
* `TestContactRequirement`: Ensure contact info presence.
* `TestOriginValidation`: Verify `OwnerOrigin` against West African nations.

#### 4.2 Integration Tests

* `TestExpirationLogic`: A background ticker service must find listings where `Deadline` or `EventEnd` has passed and set `IsActive = false`.

### 5. Deployment Strategy

* **Compute:** Dockerized application deployable to Fly.io or similar container platforms.
* **Database:** SQLite (Embedded) for MVP simplicity and portability.
* **CI/CD:** GitHub Actions / Pre-commit hooks to run tests before functionality changes.

### 6. Admin & Moderation

* **Admin Dashboard:** A secured area for administrators to view system metrics and moderate content.
* **Moderation:** Listings can be approved or rejected.
* **Access Control:** Admin access is protected via Google Auth + a secondary Access Code.

### 7. UI Design System (Brand Guidelines)

* **Colors:** Consistent use of `stone-*` tokens for neutrals and `primary` / `secondary` brand colors for accents.
* **Shapes:** Unified shape language utilizing `rounded-3xl` (24px) for cards/modals and `rounded-xl` for form inputs.
* **Component-Based Styling:** All templates link to a compiled `output.css` to ensure consistent utility classes and theme tokens are applied.

### 8. Codebase Critique & Improvements (Self-Correction)

* **UI Consistency:** The administrative interface was refactored from "programmer art" (generic grays and small radii) to a premium, brand-aligned experience. This involved standardizing on `stone` tokens and `rounded-3xl` across all templates.
* **Loose Coupling**: The `TemplateRenderer` was refactored to isolate page templates, preventing namespace collisions.
* **Security**: `DevLogin` in `auth.go` currently bypasses the Admin Claim flow. This should be tightened in future iterations.
* **Testing**: Browser subagent tests proved critical in catching template errors (like unterminated strings) and visual regressions during the rebranding process.
* **Bulk Upload**: Admin bulk upload requires confirmation and gracefully handles both malformed files (by redirecting to the dashboard with flash messages) and invalid categories (by falling back to the `Business` type).