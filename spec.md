## spec.md: Agbalumo MVP (Golang Edition)

### 1. Project Overview

**Agbalumo** is a high-performance directory and request platform for the West African diaspora.

* **Goal:** Connect users with West African businesses and services.
* **Architecture:** **Go (Golang)** for the backend, utilizing the `Standard Library` and `Gin/Echo` for high-concurrency routing.
* **TDD Strategy:** Leverage Go’s built-in `testing` package with `testify` for assertions.

### 2. Data Models (Go Structs)

The system will use strict types to ensure data integrity across the global search.

```go
type Category string

const (
    Business Category = "Business"
    Service  Category = "Service"
    Product  Category = "Product"
    Request  Category = "Request"
)

type Listing struct {
    ID                string    `json:"id"`
    OwnerOrigin       string    `json:"owner_origin"` // Required: Country of Origin
    Type              Category  `json:"type"`
    Anchor            string    `json:"anchor"`       // Food, Professional, etc.
    Title             string    `json:"title"`
    Description       string    `json:"description"`
    Neighborhood      string    `json:"neighborhood"`
    ContactEmail      string    `json:"contact_email"`
    ContactWhatsApp   string    `json:"contact_whatsapp"`
    CreatedAt         time.Time `json:"created_at"`
    Deadline          time.Time `json:"deadline"`     // Required for 'Request'
    IsActive          bool      `json:"is_active"`
}

```

### 3. Functional Requirements & Constraints

* **Concurrency:** Use Goroutines for non-blocking Gemini AI tagging and moderation tasks.
* **Deadline Validation:** `Request` types must have a `Deadline` within  days of `CreatedAt`.
* **Contact Integrity:** Every listing must contain at least one valid communication method (WhatsApp, Email, or Phone).
* **Cultural Filter:** An internal service `Moderator` will interface with the Gemini API to verify the cultural relevancy of the `Description`.

### 4. TDD Specification

#### 4.1 Unit Tests (`listing_test.go`)

* `TestValidateDeadline`: Verify that a deadline  days returns a custom `ErrInvalidDeadline`.
* `TestContactRequirement`: Ensure that an empty contact card returns `ErrMissingContact`.
* `TestOriginValidation`: Verify that the `OwnerOrigin` is not empty and matches a supported list of West African nations.

#### 4.2 Integration Tests

* `TestExpirationLogic`: A background ticker service must find listings where `time.Now() > Deadline` and set `IsActive = false`.
* `TestGeminiResponse`: Mock the Gemini client to ensure the system handles "Uncertain" cultural matches by flagging them for manual review instead of failing silently.

### 5. Deployment Strategy (Low Cost)

* **Compute:** AWS App Runner or AWS Lambda (Go 1.x runtime).
* **Database:** Supabase (Postgres) or PlanetScale (MySQL) for the free-tier managed relational database.
* **CI/CD:** GitHub Actions to run `go test ./...` before every deploy to ensure the TDD "Green" state is maintained.