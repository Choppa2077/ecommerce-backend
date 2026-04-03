# Metrics Report — Test Automation

**Project:** E-Commerce Platform  
**Assignment:** Assignment 2 — Test Automation Implementation  
**Date:** 2026-04-03  
**Author:** QA Team

---

## 1. Automation Coverage

Measures the proportion of HIGH-risk functions that have automated test scripts.

| Module / Feature | High-Risk Function | Automated? | Coverage % | Notes |
|------------------|--------------------|------------|------------|-------|
| Auth | User registration | Yes | 100% | `TestRegister_Success`, `TestRegister_DuplicateEmail` |
| Auth | User login | Yes | 100% | `TestLogin_ValidCredentials`, `TestLogin_InvalidPassword`, `TestLogin_UserNotFound`, `TestLogin_InactiveUser` |
| Auth | JWT token validation | Yes | 100% | `TestValidateToken_Valid`, `TestValidateToken_Invalid`, `TestValidateToken_WrongSecret` |
| Product CRUD | Create product | Yes | 100% | `TestCreateProduct_Success`, `TestCreateProduct_InvalidName`, `TestCreateProduct_InvalidPrice`, `TestCreateProduct_CategoryNotFound` |
| Product CRUD | Get product | Yes | 100% | `TestGetProduct_Success`, `TestGetProduct_NotFound` |
| Product CRUD | Update product | Yes | 100% | `TestUpdateProduct_Success`, `TestUpdateProduct_NotFound` |
| Product CRUD | Delete product | Yes | 100% | `TestDeleteProduct_Success`, `TestDeleteProduct_NotFound` |
| Product CRUD | List products | Yes | 80% | `TestListProducts_WithFilter` — edge cases not yet covered |
| Purchase Flow | Purchase with stock check | Yes | 100% | `TestPurchaseProduct_Success`, `TestPurchaseProduct_InsufficientStock`, `TestPurchaseProduct_ZeroQuantity`, `TestPurchaseProduct_ProductNotFound` |
| Interactions | Like product | Yes | 100% | `TestLikeProduct_Success`, `TestLikeProduct_ProductNotFound` |
| Interactions | View product | Yes | 100% | `TestRecordProductView_Success` |
| Interactions | Purchase history | Yes | 100% | `TestGetUserPurchaseHistory_Success`, `TestGetUserPurchaseHistory_Empty` |
| Interactions | Is liked check | Yes | 100% | `TestIsProductLiked_True`, `TestIsProductLiked_False` |
| Recommendation Engine | Collaborative filtering | No | 0% | Planned for next iteration |
| Frontend Auth | Login page structure | Yes | 100% | Playwright: `e2e/auth.spec.ts` (9 tests) |
| Frontend Auth | Register page structure | Yes | 100% | Playwright: `e2e/auth.spec.ts` |
| Frontend Routes | Protected route redirect | Yes | 100% | Playwright: `e2e/profile.spec.ts` (6 routes tested) |
| Frontend Products | Product page redirect | Yes | 100% | Playwright: `e2e/products.spec.ts` |

**Formula:**  
`Automation Coverage (%) = (Automated high-risk functions / Total high-risk functions) × 100`

**Result:** 15/16 HIGH-risk functions automated = **93.75% automation coverage**

---

## 2. Test Execution Time (TTE)

### Backend Unit Tests

| Module / Feature | Number of Tests | Total Execution Time | Average per Test | Notes |
|-----------------|-----------------|---------------------|-----------------|-------|
| Auth (authService) | 9 | ~0.1s | ~0.01s | Includes JWT signing |
| Products (productService) | 13 | ~0.1s | ~0.01s | Pure business logic, mocked repo |
| Interactions (interactionService) | 9 | ~0.1s | ~0.01s | Mock-based, no I/O |
| **Total** | **31** | **~0.7s** | **~0.02s** | Well within 2-minute threshold |

### Frontend E2E Tests (Playwright)

| Spec File | Number of Tests | Estimated Time | Notes |
|-----------|-----------------|----------------|-------|
| `e2e/auth.spec.ts` | 9 | ~10-15s | Page navigation + DOM checks |
| `e2e/products.spec.ts` | 3 | ~5-8s | Redirect checks |
| `e2e/profile.spec.ts` | 9 | ~12-18s | Protected route checks |
| **Total** | **21** | **~30-45s** | Runs in Chromium only |

---

## 3. Defects Found vs Expected Risk

Defects discovered during automation compared to the expected risk levels from Assignment 1.

| Module | High-Risk Level | Expected Defects | Defects Found | Pass/Fail | Notes |
|--------|----------------|-----------------|---------------|-----------|-------|
| Auth (register) | HIGH | 2 | 0 | PASS | No defects in happy path; duplicate email correctly rejected |
| Auth (login) | HIGH | 2 | 0 | PASS | Invalid credentials correctly handled |
| Auth (JWT) | HIGH | 1 | 0 | PASS | Token validation works; wrong secret correctly rejected |
| Product CRUD | HIGH | 2 | 0 | PASS | Validation errors correctly thrown for empty name, negative price |
| Purchase Flow | HIGH | 3 | 0 | PASS | Stock check works; zero quantity rejected |
| Interactions | MEDIUM | 1 | 0 | PASS | Like/view/purchase flows correct |
| Frontend Routes | HIGH | 2 | 0 | PASS | All protected routes redirect to /login |

**Previously identified known defects (from BASELINE_METRICS.md):**

| Bug ID | Module | Severity | Status |
|--------|--------|----------|--------|
| BUG-001 | Product CRUD | Critical | Not fixed — admin role check not implemented |
| BUG-002 | Category CRUD | Critical | Not fixed — admin role check not implemented |
| BUG-003 | Purchase Flow | High | Not fixed — no atomic transaction |
| BUG-004 | Profile | Medium | Not fixed — JWT not invalidated on password change |
| BUG-005 | Auth | Medium | Not fixed — no rate limiting on login |

> **Note:** These known defects are architectural issues not detectable by unit tests (which test business logic with mocks). They require integration testing to verify.

---

## 4. Test Execution Log

| Test Case ID | Module | Test Name | Execution Date/Time | Result | Execution Time | Notes |
|-------------|--------|-----------|---------------------|--------|----------------|-------|
| TC-AUTH-01 | Auth | TestRegister_Success | 2026-04-03 | PASS | 0.00s | — |
| TC-AUTH-02 | Auth | TestRegister_DuplicateEmail | 2026-04-03 | PASS | 0.00s | — |
| TC-AUTH-03 | Auth | TestLogin_ValidCredentials | 2026-04-03 | PASS | 0.00s | bcrypt MinCost used for test speed |
| TC-AUTH-04 | Auth | TestLogin_InvalidPassword | 2026-04-03 | PASS | 0.00s | — |
| TC-AUTH-05 | Auth | TestLogin_UserNotFound | 2026-04-03 | PASS | 0.00s | Returns ErrInvalidCredentials (not ErrNotFound) |
| TC-AUTH-06 | Auth | TestLogin_InactiveUser | 2026-04-03 | PASS | 0.00s | Status "suspended" → ErrUserInactive |
| TC-AUTH-07 | Auth | TestValidateToken_Valid | 2026-04-03 | PASS | 0.00s | — |
| TC-AUTH-08 | Auth | TestValidateToken_Invalid | 2026-04-03 | PASS | 0.00s | — |
| TC-AUTH-09 | Auth | TestValidateToken_WrongSecret | 2026-04-03 | PASS | 0.00s | — |
| TC-PROD-01 | Products | TestCreateProduct_Success | 2026-04-03 | PASS | 0.00s | IsActive set to true |
| TC-PROD-02 | Products | TestCreateProduct_InvalidName | 2026-04-03 | PASS | 0.00s | Validation catches empty name |
| TC-PROD-03 | Products | TestCreateProduct_InvalidPrice | 2026-04-03 | PASS | 0.00s | Negative price rejected |
| TC-PROD-04 | Products | TestCreateProduct_CategoryNotFound | 2026-04-03 | PASS | 0.00s | — |
| TC-PROD-05 | Products | TestGetProduct_Success | 2026-04-03 | PASS | 0.00s | — |
| TC-PROD-06 | Products | TestGetProduct_NotFound | 2026-04-03 | PASS | 0.00s | Returns ErrNotFound |
| TC-PROD-07 | Products | TestUpdateProduct_Success | 2026-04-03 | PASS | 0.00s | — |
| TC-PROD-08 | Products | TestUpdateProduct_NotFound | 2026-04-03 | PASS | 0.00s | — |
| TC-PROD-09 | Products | TestDeleteProduct_Success | 2026-04-03 | PASS | 0.00s | GetByID called before Delete |
| TC-PROD-10 | Products | TestDeleteProduct_NotFound | 2026-04-03 | PASS | 0.00s | — |
| TC-PROD-11 | Products | TestListProducts_WithFilter | 2026-04-03 | PASS | 0.00s | IsActive default applied |
| TC-INT-01 | Interactions | TestPurchaseProduct_Success | 2026-04-03 | PASS | 0.00s | Stock decremented correctly |
| TC-INT-02 | Interactions | TestPurchaseProduct_InsufficientStock | 2026-04-03 | PASS | 0.00s | Error: "insufficient stock" |
| TC-INT-03 | Interactions | TestPurchaseProduct_ZeroQuantity | 2026-04-03 | PASS | 0.00s | Error: "quantity must be greater than 0" |
| TC-INT-04 | Interactions | TestPurchaseProduct_ProductNotFound | 2026-04-03 | PASS | 0.00s | — |
| TC-INT-05 | Interactions | TestLikeProduct_Success | 2026-04-03 | PASS | 0.00s | — |
| TC-INT-06 | Interactions | TestLikeProduct_ProductNotFound | 2026-04-03 | PASS | 0.00s | — |
| TC-INT-07 | Interactions | TestGetUserPurchaseHistory_Success | 2026-04-03 | PASS | 0.00s | Default limit 50 applied |
| TC-INT-08 | Interactions | TestGetUserPurchaseHistory_Empty | 2026-04-03 | PASS | 0.00s | — |
| TC-INT-09 | Interactions | TestRecordProductView_Success | 2026-04-03 | PASS | 0.00s | — |
| TC-INT-10 | Interactions | TestIsProductLiked_True | 2026-04-03 | PASS | 0.00s | — |
| TC-INT-11 | Interactions | TestIsProductLiked_False | 2026-04-03 | PASS | 0.00s | — |

**Total: 31 tests — 31 PASS, 0 FAIL**

---

## 5. Metrics Summary

| Metric | Value |
|--------|-------|
| Total automated test cases | 31 (unit) + 21 (E2E) = **52** |
| HIGH-risk modules with automation | **5/6** (83%) |
| HIGH-risk functions automated | **15/16** (93.75%) |
| Unit test pass rate | **100%** (31/31) |
| Total unit test execution time | **~0.7 seconds** |
| Known defects (pre-existing) | **5** (architectural, not unit-testable) |
| New defects found by automation | **0** |
| Quality gates passed | **6/7** (QG04 pending CI run) |

---

## 6. Evidence for Reproducibility

| Evidence ID | Module | Type | Description | Location |
|-------------|--------|------|-------------|----------|
| E01 | Auth | Go test output | 9 auth tests passing | `internal/service/authService_test.go` |
| E02 | Products | Go test output | 13 product tests passing | `internal/service/productService_test.go` |
| E03 | Interactions | Go test output | 9 interaction tests passing | `internal/service/interactionService_test.go` |
| E04 | Frontend | Playwright spec | Auth flow E2E tests | `e2e/auth.spec.ts` |
| E05 | Frontend | Playwright spec | Product redirect tests | `e2e/products.spec.ts` |
| E06 | Frontend | Playwright spec | Profile redirect tests | `e2e/profile.spec.ts` |
| E07 | CI/CD | YAML pipeline | Backend pipeline with go test | `.github/workflows/ci.yml` |
| E08 | CI/CD | YAML pipeline | Frontend pipeline with Playwright | `.github/workflows/ci.yml` (frontend repo) |

**To rerun unit tests:**
```bash
go test ./internal/service/... -v -timeout 60s
```

**To rerun E2E tests:**
```bash
cd /path/to/frontend
npx playwright install chromium
npm run test:e2e
```
