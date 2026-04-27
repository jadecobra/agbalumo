# Decision Log

- **2026-04-27**: Initialized task to investigate and fix unpopulated `menuURL` fields in production listings. Need to verify current scraper behavior and update it to look for variations like `/menu`, `/order`, or `/order-online`.
- **2026-04-27**: Confirmed via browser audit of agbalumo.com that listings like Native Restaurant & Lounge, AGEGE BUKKAH, and Joloff lack menu URLs.
- **2026-04-27**: Identified that `internal/service/scraper.go` only checks `href` attributes for keywords. Propose updating it to inspect anchor text (e.g., `<a href="/foo">Order Online</a>`) to capture missed menus without introducing excessive DOM tree parsing latency.
- **2026-04-27**: Proactively rejected adding `goquery` or similar full DOM tree parsers to maintain minimum memory and latency footprint in the background processing job.

# Execution Plan

- [x] Phase 1: Investigate current scraper and identify where menuURL logic is.
- [x] Phase 2: Add test cases in Go for new URL patterns.
- [x] Phase 3: Implement updated scraper logic to find variations.
- [x] Phase 4: Verify the fix with a live run or automated test suite.
