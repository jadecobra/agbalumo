## spec.md: Agbalumo MVP

### 1. Project Overview

**Agbalumo** is a high-performance directory and request platform for the West African diaspora.

* **Goal:** Connect users with West African businesses and services.
* **Architecture:** Python-based, TDD-driven, designed for low-cost serverless deployment.
* **Core Philosophy:** Cultural authenticity (story-driven) with technical minimalism.

### 2. Data Models

#### 2.1 The Listing Entity

* **Owner_Info:** Name, Country of Origin (Required).
* **Category:** [Business, Service, Product, Request].
* **Anchors:** [Food, Professional Services, Retail, Community, Creative].
* **Content:** Title, Description (Gemini-analyzed for AI tagging).
* **Location:** Neighborhood/City, State, Country.
* **Contact_Card:** Email, Phone, WhatsApp link, URL (Optional).
* **Metadata:** `created_at`, `is_active`.
* **Deadline:** Required for 'Request' type (Max 90 days).

### 3. Functional Requirements

* **Global Discovery (Phase 1):** Interest-based search results sorted by recency/relevancy.
* **Lifecycle Management:** Automatic expiration of 'Requests' once the `deadline` is reached.
* **Validation:** * Deadlines cannot exceed 90 days from the current date.
* Contact cards must have at least one valid communication method.
* AI Moderation: Gemini must verify the description aligns with West African cultural context.



### 4. TDD Test Suite (Pytest Framework)

#### 4.1 Unit Tests

* `test_request_deadline_valid`: Pass if date <= 90 days, fail if > 90 or in the past.
* `test_listing_origin_required`: Fail if Country of Origin is missing.
* `test_contact_card_format`: Validate WhatsApp URL formatting and email syntax.

#### 4.2 Integration Tests

* `test_gemini_moderation_logic`: Mock Gemini response to ensure "Naija Suya" passes and "Generic Pizza" fails.
* `test_expiration_logic`: Verify that a post with a deadline of `Yesterday` is marked as `is_active=False`.

### 5. Deployment Strategy

* **Backend:** FastHTML/FastAPI.
* **Database:** SQLite (local development) / Supabase (production).
* **Storage:** AWS S3 for images.

---

### How to use this with Antigravity:

1. Open your Antigravity environment.
2. Upload/Paste this `spec.md`.
3. **Prompt Gemini:** *"Based on the attached spec.md, let's start the TDD cycle. Generate the `pytest` file for the Listing Entity and its validation logic first. Do not write the implementation until the tests are defined."*

**Would you like me to refine any specific field in the data model before you start, or are you ready to jump into the first TDD cycle?**