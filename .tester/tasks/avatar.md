# Implementation Plan - Defining the First Avatar (Aisha)

This plan defines the primary user avatar for Agbalumo to focus product development and design decisions. By anchoring our building efforts to a specific person with clear goals and pain points, we ensure the platform delivers immediate value to the West African diaspora community.

## User Review Required

> [!IMPORTANT]
> The selection of the **first** avatar is a strategic decision that will dictate the priority of upcoming features (e.g., photo galleries, WhatsApp integration, regional filtering).
> Please review the avatar profile and the "Strategic Critique" to ensure it aligns with your vision for Agbalumo.

## Proposed Avatar: Ada, the Quality-Obsessed Consultant

| Attribute | Details |
| :--- | :--- |
| **Name** | Ada |
| **Occupation** | Management Consultant (Frequent Traveler) |
| **Origin** | Nigeria |
| **Location** | Dallas, TX (Downtown) |
| **Core Goal** | To find the top 3 high-quality Nigerian food locations in any city in under 60 seconds. |

### Pain Point Mapping (Ada)

1.  **Quality Lottery**: Ada knows how to cook; she has high expectations. Currently, trying a new place is a gamble. She needs a **Quality Signal** (consistent peer recommendations) that she can trust *before* she orders.
2.  **Time-to-Comfort**: After a long travel day (e.g. Tuesday night in a new city), she needs immediate "home" comfort. The discovery process must be ultra-fast (< 60 seconds).
3.  **Community Isolation**: She wants to see "people like her" ("Nigerians are everywhere"). The platform should signal not just food, but a safe/authentic community space.

By focusing on the **"60-Second Value"** for Ada, we prioritize:
- **Location-Aware Speed**: Finding what's closest and highest quality immediately.
- **Default "Food" Experience**: Setting the homepage to default to Food listings to eliminate search friction (per user feedback).
- **Consistency over Volume**: Showing only the "Top 3" ensures Ada isn't overwhelmed and we maintain a high-quality bar.

## Technical Alignment

### [Component Name] Domain & Content

#### [MODIFY] [listing.go](file:///Users/johnnyblase/gym/agbalumo/internal/module/listing/listing.go)
- **Default Homepage Filter**: Update `HandleHome` to default the `Category` filter to `Food` (Nigeria-Origin) to minimize Ada's time-to-value.
- **Latency Instrumentation**: Add `time.Since` logging in `HandleFragment` to monitor and enforce the < 200ms threshold for Ada's "60-Second Hunt".

#### [MODIFY] [user_journeys.yaml](file:///Users/johnnyblase/gym/agbalumo/.agents/user_journeys.yaml)
- Add "Ada's Tuesday Night Hunt": A journey specifically testing the < 60s discovery of top-tier Nigerian food in a new city.

## Knowledge Alignment

- **Existing Patterns**: Uses the existing `Listing` and `Category` models (`internal/domain/listing.go`).
- **Memory Search**: Confirmed that the project core is "West African diaspora", which aligns perfectly with Aisha's Nigerian origin.

## Requirements & Constraints

1.  **Speed**: The `/listings/fragment` endpoint MUST respond in < 200ms.
2.  **Default State**: The homepage (index) MUST default to the `Food` category for unauthenticated visitors.
3.  **Performance Measurement**: We will measure in-browser performance using the `browser_subagent` to evaluate `PerformanceNavigationTiming` and HTMX "request-to-swap" latency via console metrics.

## Verification Plan

### Automated Tests
- `task test`: Ensure existing listing validation still passes with "Food" and "Business" types.
- `harness verify api-spec`: Confirm no contract regressions.

### Manual Verification
- Walk through the **"Ada's Journey"** using the browser subagent to ensure the UI "feels" right for her specific needs (speed, quality-first, home comfort).
