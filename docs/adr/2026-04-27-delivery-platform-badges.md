# ADR 013: Delivery Platform Badges
**Date**: 2026-04-27 **Status**: Accepted

## 1. Context & User Problem
Ada needs to know how to get food from African restaurants quickly. Currently, there is no visibility into which delivery platforms (UberEats, DoorDash, Grubhub) these businesses use, requiring her to check each app manually.

## 2. Decision
Implement passive delivery badges on listing detail modals. We store normalized delivery platform names as a JSON string array in a `delivery_platforms` SQLite column and use a template helper `hasDelivery` to render visual indicators.

## 3. The Complexity Kill-Switch (Rationale)
* **User Value**: Ada gets immediate visibility into ordering options without clicking out to dead ends.
* **Performance Budget**: Latency impact is <1ms. It's a single database column update and simple string/JSON operations in memory.
* **Minimalism Check**: Avoided introducing complex deep-linking URL requirements or extra tables; we store string tags directly on the listing.

## 4. Consequences
* **Technical Tradeoffs**: Passive badges provide awareness but lack a single-click CTA directly into the delivery app.
* **Observability**: Checked visually via automated browser integration flows.
* **SQLite Impact**: Handled with zero locks; added column with standard default values and no extra foreign keys.

## 5. Alternatives Considered
* **Clickable Deep Links**: Rejected because scraper heuristic reliability was too low for direct deep link resolution, leading to poor UX.
