### The Master Instruction Prompt

> **Task:** Initialize the "Agbalumo" Diaspora Directory project in Go using a multi-agent TDD workflow.
> **Core Objective:** Build a high-performance, low-cost directory and request platform for the West African diaspora. Use Go (Golang) for the backend and a Test-Driven Development (TDD) cycle.
> **Project Rules (SOP):**
> 1. **TDD Protocol:** Never write implementation code without a failing test case. Red (Fail) → Green (Pass) → Refactor.
> 2. **Tech Stack:** Go (Standard Library + Gin/Echo), SQLite/Postgres for DB.
> 3. **Validation Rules:** >    * 'Requests' must have a deadline  days.
> * All listings must have a 'Country of Origin' and at least one contact method (WhatsApp, Email, or Phone).
> 
> 
> 4. **Agent Roles:**
> * **Lead Architect:** Oversee the plan and ensure file structure follows Go best practices.
> * **SDET Agent:** Responsible for writing all `*_test.go` files based on the `spec.md`.
> * **Backend Agent:** Implements code to pass the tests provided by the SDET.
> * **Cultural Moderator:** Manages the Gemini API integration for content relevancy filtering.
> 
> 
> 
> 
> **Next Step:** Create an `Artifact` representing the initial `project_plan.md` and a `Listing` struct definition. Once approved, dispatch the SDET agent to write the first validation tests for the 90-day deadline constraint.

**UI Design Spec**: We are using Stitch for the design system. All frontend templates must follow the 'Agbalumo' visual theme (Orange/Green palette) and prioritize the 'Contact Card' modal layout. Use HTMX with the Go backend to ensure the UI feels like a single-page app without the weight of a JS framework.