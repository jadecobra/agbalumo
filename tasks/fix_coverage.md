# Test Coverage Improvement Plan

## Current Coverage Analysis

Based on the latest test coverage report, overall coverage is at 86.5% with several critical gaps. The main issues are:

### Critical Coverage Gaps (<80%):
- `github.com/jadecobra/agbalumo/cmd/security-audit/main.go:20: Run (0.0%)
- `github.com/jadecobra/agbalumo/cmd/security-audit/main.go:27: main (0.0%)
- `github.com/jadecobra/agbalumo/internal/middleware/ratelimit.go:47: purge (0.0%)
- `github.com/jadecobra/agbalumo/internal/middleware/ratelimit.go:58: getVisitor (0.0%)
- `github.com/jadecobra/agbalumo/internal/middleware/session.go:41: NewTestSessionStore (0.0%)
- `github.com/jadecobra/agbalumo/internal/mock/image_service.go:14: UploadImage (0.0%)
- `github.com/jadecobra/agbalumo/internal/mock/image_service.go:19: DeleteImage (0.0%)
- `github.com/jadecobra/agbalumo/internal/mock/renderer.go:11: Render (0.0%)
- `github.com/jadecobra/agbalumo/internal/repository/sqlite/sqlite.go:484: TitleExists (0.0%)
- `github.com/jadecobra/agbalumo/internal/repository/sqlite/sqlite.go:889: GetClaimRequestByUserAndListing (0.0%)
- `github.com/jadecobra/agbalumo/internal/repository/sqlite/sqlite.go:129: GetCategories (0.0%)
- `github.com/jadecobra/agbalumo/internal/repository/sqlite/sqlite.go:140: GetCategory (0.0%)
- `github.com/jadecobra/agbalumo/internal/repository/sqlite/sqlite.go:145: SaveCategory (0.0%)

### Moderate Coverage Gaps (70-85%):
- `github.com/jadecobra/agbalumo/internal/config/config.go:45: getAdminCode (75.0%)
- `github.com/jadecobra/agbalumo/internal/domain/listing.go:171: validateRequest (77.8%)
- `github.com/jadecobra/agbalumo/internal/repository/sqlite/sqlite.go:76: NewSQLiteRepository (70.0%)
- `github.com/jadecobra/agbalumo/internal/repository/sqlite/sqlite.go:308: FindAll (73.5%)
- `github.com/jadecobra/agbalumo/internal/repository/sqlite/sqlite.go:402: SaveUser (83.3%)
- `github.com/jadecobra/agbalumo/internal/repository/sqlite/sqlite.go:459: FindAllByOwner (84.6%)
- `github.com/jadecobra/agbalumo/internal/repository/sqlite/sqlite.go:378: FindByTitle (84.6%)
- `github.com/jadecobra/agbalumo/internal/repository/sqlite/sqlite.go:527: GetLocations (92.3%)
- `github.com/jadecobra/agbalumo/internal/repository/sqlite/sqlite.go:599: GetAllUsers (84.6%)
- `github.com/jadecobra/agbalumo/internal/repository/sqlite/sqlite.go:619: GetFeaturedListings (84.6%)
- `github.com/jadecobra/agbalumo/internal/repository/sqlite/sqlite.go:835: GetPendingClaimRequests (84.6%)
- `github.com/jadecobra/agbalumo/internal/repository/sqlite/sqlite.go:860: UpdateClaimRequestStatus (76.9%)
- `github.com/jadecobra/agbalumo/internal/seeder/category_seeder.go:15: EnsureCategoriesSeeded (88.2%)
- `github.com/jadecobra/agbalumo/internal/service/background.go:21: StartTicker (88.9%)
- `github.com/jadecobra/agbalumo/internal/service/csv.go:18: NewCSVService (100.0%)
- `github.com/jadecobra/agbalumo/internal/service/csv.go:23: ParseAndImport (87.0%)
- `github.com/jadecobra/agbalumo/internal/service/image.go:39: NewLocalImageService (100.0%)
- `github.com/jadecobra/agbalumo/internal/service/image.go:53: UploadImage (87.3%)
- `github.com/jadecobra/agbalumo/internal/service/listing_service.go:25: NewListingService (100.0%)
- `github.com/jadecobra/agbalumo/internal/service/listing_service.go:32: ClaimListing (94.7%)
- `github.com/jadecobra/agbalumo/internal/ui/renderer.go:21: NewTemplateRenderer (92.9%)
- `github.com/jadecobra/agbalumo/internal/ui/renderer.go:108: compileTemplates (85.7%)
- `github.com/jadecobra/agbalumo/internal/ui/renderer.go:137: Render (91.7%)

## Priority 1: Critical Security and Core Functions (Must Fix First)

### 1. Security Audit Main Functions (0.0%)
**Files:** `cmd/security-audit/main.go`
**Functions:** `Run`, `main`

**Why Critical:** Security audit functionality is essential for vulnerability detection.

**Test Cases Needed:**
- Test security audit initialization and configuration
- Test audit execution with various security checks
- Test audit output and reporting
- Test error handling in audit process

### 2. Rate Limiting Core Functions (0.0%)
**Files:** `internal/middleware/ratelimit.go`
**Functions:** `purge`, `getVisitor`

**Why Critical:** Rate limiting prevents abuse and protects system resources.

**Test Cases Needed:**
- Test visitor tracking and identification
- Test rate limit purging logic
- Test concurrent access scenarios
- Test edge cases for rate limiting

### 3. Session Management (0.0%)
**Files:** `internal/middleware/session.go`
**Functions:** `NewTestSessionStore`

**Why Critical:** Session management is essential for user authentication and state.

**Test Cases Needed:**
- Test session store creation
- Test session data storage and retrieval
- Test session expiration and cleanup
- Test concurrent session access

### 4. Mock Functions (0.0%)
**Files:** `internal/mock/`
**Functions:** `UploadImage`, `DeleteImage`, `Render`

**Why Important:** Mocks are used extensively in testing - need proper implementation.

**Test Cases Needed:**
- Test mock behavior for image operations
- Test mock rendering functionality
- Test mock repository operations
- Test mock error scenarios

### 5. SQLite Repository Functions (0.0%)
**Files:** `internal/repository/sqlite/sqlite.go`
**Functions:** `TitleExists`, `GetClaimRequestByUserAndListing`, `GetCategories`, `GetCategory`, `SaveCategory`

**Why Critical:** Database operations are fundamental to application functionality.

**Test Cases Needed:**
- Test category CRUD operations
- Test claim request queries
- Test title existence checks
- Test database transaction handling

## Priority 2: High-Impact Business Logic (Next)

### 6. Configuration Management (75.0%)
**File:** `internal/config/config.go`
**Function:** `getAdminCode`

**Why Important:** Admin code generation is critical for security.

**Test Cases Needed:**
- Test admin code generation logic
- Test code validation and verification
- Test code expiration and rotation
- Test error scenarios for code generation

### 7. Domain Validation (77.8%)
**File:** `internal/domain/listing.go`
**Function:** `validateRequest`

**Why Important:** Input validation prevents security vulnerabilities.

**Test Cases Needed:**
- Test all validation rules and edge cases
- Test invalid input scenarios
- Test boundary conditions
- Test performance with large inputs

### 8. SQLite Repository (70-85%)
**File:** `internal/repository/sqlite/sqlite.go`
**Functions:** `NewSQLiteRepository`, `FindAll`, `SaveUser`, `FindAllByOwner`, `FindByTitle`, `GetLocations`, `GetAllUsers`, `GetFeaturedListings`, `GetPendingClaimRequests`, `UpdateClaimRequestStatus`

**Why Critical:** Database operations need comprehensive testing.

**Test Cases Needed:**
- Test all database query methods
- Test data consistency and integrity
- Test concurrent access
- Test error handling and recovery
- Test performance with large datasets

## Priority 3: Supporting Infrastructure (Lower Priority)

### 9. Seeder Functions (88.2%)
**File:** `internal/seeder/category_seeder.go`
**Function:** `EnsureCategoriesSeeded`

**Why Important:** Data seeding ensures consistent test environments.

**Test Cases Needed:**
- Test category seeding logic
- Test duplicate prevention
- Test seeding order and dependencies
- Test error handling during seeding

### 10. Background Services (88.9%)
**File:** `internal/service/background.go`
**Function:** `StartTicker`

**Why Important:** Background tasks need reliable operation.

**Test Cases Needed:**
- Test ticker initialization and scheduling
- Test background task execution
- Test error handling and recovery
- Test resource cleanup

### 11. CSV Service (87.0%)
**File:** `internal/service/csv.go`
**Function:** `ParseAndImport`

**Why Important:** CSV import functionality for data migration.

**Test Cases Needed:**
- Test CSV parsing with various formats
- Test data validation and error handling
- Test import performance with large files
- Test rollback and recovery scenarios

### 12. Image Service (87.3%)
**File:** `internal/service/image.go`
**Function:** `UploadImage`

**Why Important:** Image handling is critical for user experience.

**Test Cases Needed:**
- Test image upload with various formats
- Test image compression and conversion
- Test error handling for invalid images
- Test performance with large images

### 13. Listing Service (94.7%)
**File:** `internal/service/listing_service.go`
**Function:** `ClaimListing`

**Why Important:** Business logic for listing claims.

**Test Cases Needed:**
- Test claim workflow and validation
- Test concurrent claim attempts
- Test claim expiration and cleanup
- Test error handling for invalid claims

### 14. UI Renderer (92.9%)
**File:** `internal/ui/renderer.go`
**Functions:** `compileTemplates`, `Render`

**Why Important:** Template rendering affects user interface.

**Test Cases Needed:**
- Test template compilation and caching
- Test template rendering with various data
- Test error handling for invalid templates
- Test performance with complex templates

## Implementation Strategy

### Phase 1: Critical Security Functions (Week 1)
1. Complete security audit test coverage
2. Implement rate limiting tests
3. Add session management tests
4. Fix mock implementations
5. Complete SQLite repository core tests

### Phase 2: Business Logic (Week 2)
1. Complete configuration validation tests
2. Add comprehensive domain validation tests
3. Finish SQLite repository testing
4. Complete seeder and background service tests

### Phase 3: Infrastructure (Week 3)
1. Complete CSV service testing
2. Add image service tests
3. Finish listing service testing
4. Complete UI renderer testing

## Testing Guidelines

1. **TDD Approach:** Write tests before implementing new functionality
2. **Edge Cases:** Test all boundary conditions and error scenarios
3. **Concurrency:** Test concurrent access where applicable
4. **Performance:** Include performance tests for critical paths
5. **Integration:** Test database operations with actual SQLite
6. **Mocking:** Use mocks appropriately to isolate units

## Success Criteria

- Achieve 95%+ overall test coverage
- All critical functions at 100% coverage
- All business logic functions at 90%+ coverage
- All infrastructure functions at 80%+ coverage
- All tests pass without errors
- Build process completes successfully

## Next Steps

1. Begin with Phase 1 critical security functions
2. Use TDD approach for all new tests
3. Run tests frequently to catch regressions early
4. Update documentation as tests are added
5. Verify coverage improves with each change

## Notes

- Focus on critical security and core business logic first
- Use existing test patterns in the codebase
- Maintain 90%+ coverage threshold throughout development
- Consider test performance and execution time
- Ensure tests are reliable and deterministic