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

### 7. UI Design System (TOON Tokens)

Agbalumo uses a "Juicy Fruit" design aesthetic, inspired by Yoruba roots and community-first values.

* **Colors:**
  - `primary`: #FF8A00 (Agbalumo orange)
  - `secondary`: #689F38 (Leaf green for CTAs)
  - `accent-star`: #C2185B (Seed magenta for badges)
  - `surface`: #FFF8F0 (Warm cream)
  - `background`: #FFFBF5 (Inviting off-white)
* **Typography:** `Playfair Display` (Serif) for headings to add warmth; `Inter` (Sans) for high-performance UI components.
* **Shapes:** `rounded-3xl` (32px) for a soft, premium feel; `rounded-xl` for interactive elements.
* **Motion:** `juice-bounce` (scale 0.94 → 1.02) on clicks; `gentle-pulse` for verified status and new badges.
* **Shadows:** `shadow-juicy` (orange tinted) and `shadow-lifted` for tactile depth.
* **Accessibility & Ergonomics (The 10x Standard):**
  - Semantic Landmark: `<h1>` logo wrap for screen-reader navigation.
  - Color Contrast: Listing descriptions meet WCAG AA (4.5:1 ratio).
  - Touch Targets: All primary mobile targets (chips, search, nav) are at least 44px.
  - Interaction: High-visibility 3px orange focus rings for keyboard users.
  - ARIA: Placeholders and images include roles and descriptive labels.

### 8. Codebase Critique & Improvements (Self-Correction)

* **Visual Hierarchy:** Unified CTA colors (`btn-leaf`) and mobile heading scaling (H2 > H3) successfully resolved the "chaotic layout" issues.
* **Accessibility First:** The project now adheres to WCAG AA contrast and semantic landmark standards. The "10x" engineer protocol was applied to ensure the UI is robust for all users.
* **Performance:** CSS build steps are integrated via `npm build:css`. Future work includes image optimization (WebP) and further HTMX lazy-loading.
* **Critique (Agbalumo Spirit):** The UI now feels "Sweet like Agbalumo." It is warm, accessible, and premium.