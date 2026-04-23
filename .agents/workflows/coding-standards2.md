---
description: Coding Standards and Guidelines (Go, HTMX, Tailwind)
---

Coding Standards & Guidelines

This document defines the technical and product constraints for the Agbalumo project. It is optimized for a Senior Product Engineer persona, prioritizing user utility, system transparency, and ruthless simplicity.

Product-Centric Performance (The 60-Second Goal)

The "North Star" for this project is for a user to find African food in any city in under 60 seconds.

Complexity Kill-Switch: You MUST justify the existence of any new feature or abstraction. If it increases UI friction or DB latency without a 2x increase in user utility, it must be deleted or simplified.

Bottleneck-Aware Growth: Features designed solely for user acquisition (sharing, referrals, social loops) are considered UI Bloat and MUST be rejected until listing quality (Accuracy/Verification) is no longer the primary bottleneck.

Performance Budget: The target Time to First Result (TTFR) is < 500ms on a standard mobile connection.

Latency Guardrail: Any change estimated to add >100ms to the critical search path requires a formal ADR and a justification of why the "User Value" outweighs the speed penalty.

Data Integrity & Trust Mandate

Speed is irrelevant if the data is wrong. Trust is our most expensive asset.

The Hours-to-Pulse Pipeline:

We do not call blind. To minimize nuisance and maximize "Proof of Life" accuracy:

Scraper-First Hours: The Menu URL scraper MUST prioritize extracting "Hours of Operation" text from the primary Menu URL/Official Website.

LLM Extraction: Use a lightweight LLM prompt to normalize messy "Hours" text into a standard JSON schedule.

The Scheduler: The "Phone Pulse" system MUST only initiate calls during the extracted "Open" windows.

The No-Vision Rule: If hours are embedded in images/flyers, do NOT use OCR/Vision extraction. This is considered Complexity Creep. Fall back immediately to the "Phone Pulse."

Zero-Data Fallback: If no hours are found via text scraping, use a global "Safe Window" (1 PM - 6 PM local time) for the first "Phone Pulse."

The Escalation Pulse (NLP Curation):

If the primary site scraper fails to find hours, the bot script graduates to a curation tool.

The NLP Script: "Hi, I'm from Agbalumo. We couldn't find your hours online—what are your opening hours today?"

Ambiguity Handling (The "Honest Failure" Rule): If the LLM parser confidence score is low (e.g., < 0.8), mark as "Help Us Verify".

UI Implementation: Display a "We tried to verify hours but weren't 100% sure. Can you help us?" prompt on the listing.

Zero-Cognitive-Load Curation:

The Single-Tap Rule: Interaction for "Help Us Verify" must be a binary confirmation (e.g., "Are they open right now? [Yes] [No]").

Existence Verification & Proxy Signals:

The "Phone Pulse" Protocol:

Frequency: Limit successful pulses to once every 14 days.

Success Definition: Human or IVR pickup counts as success.

Multi-Day Retry Logic: Soft failures (Busy/No Answer) require three (3) retry attempts on different days and windows within a 1-week period.

Hard Failure Action: If all three multi-day retries fail, immediately flag as "Menu Unavailable" and deprioritize.

Automated Trust Scoring (Verified Badge):

A listing is "Verified" if it has:

Freshness: Successful "Proof of Life" signal within the last 7 days.

Consistency: Zero "Broken Link" or "Closed" reports in last 30 days.

Completeness: Valid address, phone, and verified hours.

Partial Failure Honesty: If a critical data point (like a Menu URL) is broken but the restaurant is verified open, DO NOT hide the restaurant. Display the listing with a clear "Menu Unavailable" status.

Scaling Skepticism (Conflict Resolution):

If a user-tap conflicts with a recent system verification:

Early Stage (Full Trust): While users < 100, the user-tap takes immediate precedence.

Growth Stage (Skeptical): Once users > 100, require a threshold ($N > 1$) before overriding.

Code Style & Architecture (Hexagonal)

Imports & Naming

Order: Standard library, blank line, third-party, blank line, local packages.

Packages: Lowercase single word (domain, handler, service, repository).

Hexagonal Boundaries

Domain: Core types only.

Handler: HTTP logic and user-facing friction reduction.

Service: Pure business logic (The "Product Engine").

Repository: Data access (Production: SQLite).

SQLite & Storage Parity

Disk-Parity Requirement: While :memory: is used for fast TDD, critical repo changes MUST be verified against a file-backed SQLite DB.

Context Cost & Agentic Efficiency

Token Density: Target TokenRMS < 110.

Testing & Security

TDD: Write failing tests (RED) FIRST.

Strict Lessons

Zombie Data Prevention: Every repository method that fetches listings for public display MUST explicitly include a filter for last_verified_at.

JavaScript Reliability: Scripts MUST be validated with go run ./cmd/verify js-syntax.

UI Regression: Always verify UI changes with the browser subagent.

Idempotent Global Listeners: Ensure event listeners are bound once to avoid double-triggering during HTMX swaps.