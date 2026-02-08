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
* **Authentication:** Google OAuth2 for user authentication. Users can only edit/delete their own listings.
* **Deadline Validation:** 
    * `Request` types must have a `Deadline` within 90 days.
    * `Event` types must have valid start/end times.
    * `Job` types must have a valid start date.
* **Contact Integrity:** Every listing must contain at least one valid communication method.
* **Cultural Filter:** An internal service `Moderator` will interface with the Gemini API to verify the cultural relevancy of the `Description`.

### 4. TDD Specification

#### 4.1 Unit Tests (`listing_test.go`)

* `TestValidateDeadline`: Verify deadline constraints.
* `TestContactRequirement`: Ensure contact info presence.
* `TestOriginValidation`: Verify `OwnerOrigin` against West African nations.

#### 4.2 Integration Tests

* `TestExpirationLogic`: A background ticker service must find listings where `Deadline` or `EventEnd` has passed and set `IsActive = false`.
* `TestGeminiResponse`: Verify moderation logic handles "DENY" and "PERMIT" correctly.

### 5. Deployment Strategy

* **Compute:** Dockerized application deployable to Fly.io or similar container platforms.
* **Database:** SQLite (Embedded) for MVP simplicity and portability.
* **CI/CD:** GitHub Actions / Pre-commit hooks to run tests before functionality changes.