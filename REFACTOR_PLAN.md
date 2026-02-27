# agbalumo Refactoring Plan

## Goal
Refactor the agbalumo codebase to remove duplication, simplify structure, and maintain all functionality while ensuring all tests pass and API specification is met.

## Current State Analysis (Last Updated: Feb 2026)

### ✅ COMPLETED ITEMS

#### 1. Pagination Helper
- **Status**: ✅ DONE
- **File**: `internal/handler/pagination.go`
- **Implementation**: `GetPagination(c echo.Context, defaultLimit int) Pagination` function created
- **Usage**: Used in `admin.go` and `listing.go` (all pagination standardized)

#### 2. Error Handling Standardization
- **Status**: ✅ DONE
- **File**: `internal/handler/error.go`
- **Implementation**: `RespondError(c echo.Context, err error) error` function created
- **Usage**: Used across all handlers for consistent error responses

#### 3. Route Grouping (Partial)
- **Status**: ✅ PARTIALLY DONE
- **Implementation**: Admin routes use `e.Group("/admin")` with pre-applied middleware in `cmd/server.go`
- **Note**: Public and authenticated routes still use individual middleware per route

#### 4. Pagination Standardization (Complete)
- **Status**: ✅ DONE (Feb 2026)
- **File**: `internal/handler/listing.go`
- **Changes**: Replaced manual pagination in `HandleHome` and `HandleFragment` with `GetPagination(c, 20)`
- **Result**: Pagination now fully standardized across all handlers

---

### ⚠️ INCOMPLETED / OUTDATED ITEMS

#### ❌ Response Helpers
- **Status**: NOT IMPLEMENTED
- **Planned**: `RenderSuccess()`, `JSONSuccess()`, `ErrorResponse()` helpers
- **Current**: Handlers still use direct `c.Render()` calls
- **Priority**: Medium - Low impact, works fine as-is

#### ❌ Base Handler Struct
- **Status**: NOT IMPLEMENTED
- **Planned**: `BaseHandler` with common utilities (`getPagination`, `getCurrentUser`, `requireAuth`)
- **Current**: Each handler has its own structure
- **Priority**: Medium - Could reduce duplication but current duplication is minimal

#### ❌ Service Layer Consolidation
- **Status**: NOT IMPLEMENTED
- **Planned**: Enhanced service layer with `GetListings`, `CreateListing` methods
- **Current**: Minimal business logic in services
- **Priority**: Low - Current separation works adequately

#### ❌ Validation Middleware
- **Status**: NOT IMPLEMENTED
- **Planned**: `ValidatePagination` middleware
- **Current**: Manual validation in handlers
- **Priority**: Low - Not critical

---

## Current Codebase State

### What We Have:
- **29 API routes** implemented across 3 categories (Public, User, Admin)
- **90% test coverage** with all tests passing
- **Functional application** with Echo framework, SQLite, HTMX
- **Good security practices** (CSRF, session auth, rate limiting)
- **Modular structure** with clear separation of concerns

### Remaining Duplication:
1. Session Access: Repeated session and user extraction patterns in handlers

---

## Success Criteria

### Functional Criteria
- [x] All 29 API routes functional
- [x] 90%+ test coverage maintained
- [x] All tests pass
- [x] No breaking changes to existing functionality

### Code Quality Criteria
- [x] Pagination fully standardized
- [x] Error handling consistent
- [x] Route grouping partially implemented

---

## Notes

- The refactoring has achieved most of its goals for **error handling** and **pagination infrastructure**
- Pagination is now fully standardized across all handlers
- Other planned items (response helpers, base handler, service consolidation) are **nice-to-haves** that don't significantly impact code quality given the current state
