# QA Test Strategy Document

**Project:** E-Commerce Platform
**Version:** 2.0
**Date:** 2026-04-03 (updated for Assignment 2)
**Author:** QA Team

---

## 1. Project Scope & Objectives

### System Under Test

A full-stack e-commerce platform consisting of:

| Component | Technology | Repository |
|-----------|-----------|------------|
| **Backend API** | Go 1.24, Gin, MongoDB | github.com/Choppa2077/ecommerce-backend |
| **Frontend** | React 19, TypeScript, Vite | github.com/Choppa2077/ecommerce-frontend |

### Scope

**In scope:**
- All 28 REST API endpoints (auth, products, categories, profiles)
- Core user flows on frontend (login, browse, purchase, profile)
- JWT authentication and token refresh mechanism
- Purchase flow and stock management
- Recommendation engine correctness

**Out of scope:**
- Performance/load testing (future assignment)
- Security penetration testing (future assignment)
- Mobile responsiveness testing

### Objectives

1. Establish a complete test coverage baseline from 0%
2. Identify and document defects in HIGH-risk modules
3. Set up automated testing infrastructure for CI/CD
4. Provide reproducible test environment for future assignments

---

## 2. Risk Assessment Summary

See full document: [RISK_ASSESSMENT.md](./RISK_ASSESSMENT.md)

| Priority | Count | Key Modules |
|----------|-------|-------------|
| 🔴 HIGH | 5 | Auth, Token Refresh, Product CRUD, Purchase Flow, Recommendations |
| 🟡 MEDIUM | 3 | Profile, Categories, Product Filtering |
| 🟢 LOW | 1 | Like/View Interactions |

**Critical known defect:** Admin role check is not implemented (`// TODO` in source code) — all authenticated users can create/delete products.

---

## 3. Test Approach

### 3.1 Testing Levels

| Level | Type | Tool | Coverage Target |
|-------|------|------|----------------|
| **API Testing** | Manual + Automated | Postman | All 28 endpoints |
| **Unit Testing** | Automated | Go `testing` + `testify` | HIGH-risk services |
| **E2E Testing** | Automated | Playwright | Critical user flows |
| **Smoke Testing** | Manual | Postman / Browser | Core functionality |

### 3.2 Priority Order

Testing is prioritized by risk score (HIGH first):

```
1. Authentication (JWT) ────────────── Score: 15
2. Product CRUD ────────────────────── Score: 12
3. Recommendation Engine ───────────── Score: 12
4. Token Refresh Interceptor ───────── Score: 12
5. Purchase Flow + Stock ───────────── Score: 10
6. Profile Management ──────────────── Score: 6
7. Category Management ─────────────── Score: 6
8. Product Filtering ───────────────── Score: 6
9. Like/View Interactions ──────────── Score: 4
```

### 3.3 Test Approach by Module

#### Auth Module (HIGH)
- **Manual:** Postman — register, login, refresh with valid/invalid inputs
- **Automated:** Unit test `authService.go` — Register(), Login(), ValidateToken(), RefreshToken()
- **E2E:** Playwright — login flow, redirect on invalid credentials, token persistence

#### Product CRUD (HIGH)
- **Manual:** Postman — CRUD operations, verify admin-less access (known bug)
- **Automated:** Unit test `productService.go` — CreateProduct(), GetProduct(), UpdateProduct()
- **E2E:** Playwright — browse products, filter by category/price

#### Purchase Flow (HIGH)
- **Manual:** Postman — purchase with valid/invalid quantity, verify stock decrement
- **Automated:** Unit test `interactionService.go` — PurchaseProduct(), stock validation
- **E2E:** Playwright — full purchase flow from product page

#### Recommendation Engine (HIGH)
- **Manual:** Postman — hit endpoint after seeding interaction data
- **Automated:** Unit test `recommendationService.go` — scoring, cold start fallback

---

## 4. Tool Selection & Configuration

### 4.1 Postman
- **Purpose:** Manual API testing and smoke testing all 28 endpoints
- **Collection:** `tests/ecommerce.postman_collection.json`
- **Environment Variables:** `baseUrl`, `accessToken`, `refreshToken`, `productId`, `categoryId`
- **Usage:** Import collection → run Login → copy token → test endpoints

### 4.2 Go Testing + Testify (Backend Unit Tests)
- **Purpose:** Automated unit tests for service layer
- **Framework:** Go built-in `testing` package + `github.com/stretchr/testify`
- **Location:** `internal/service/*_test.go`
- **Run:** `go test ./...`

### 4.3 Playwright (Frontend E2E)
- **Purpose:** Automated end-to-end testing of user flows
- **Config:** `playwright.config.ts` (baseURL: http://localhost:5173)
- **Tests location:** `e2e/` directory
- **Run:** `npm run test:e2e`
- **Browser:** Chromium (Desktop)

### 4.4 GitHub Actions (CI/CD)
- **Backend pipeline:** `.github/workflows/ci.yml` — triggers on push/PR to main
  - Steps: `go mod download` → `go build ./...` → `go vet ./...` → `go test ./internal/service/... -v`
  - Artifact: `test-results.json` uploaded on every run
- **Frontend pipeline:** `.github/workflows/ci.yml` — triggers on push/PR to main
  - Job 1 (Build & Lint): `npm ci` → `npm run lint` → `npm run build`
  - Job 2 (E2E Tests): `npm ci` → `playwright install chromium` → `npm run test:e2e`
  - Artifact: `playwright-report/` HTML report uploaded on every run

---

## 5. Test Environment

### 5.1 Local Development Environment
| Component | Details |
|-----------|---------|
| Backend | `go run cmd/web/main.go` → http://localhost:8080 |
| Frontend | `npm run dev` → http://localhost:5173 |
| Database | MongoDB via Docker: `docker-compose up -d mongodb` |
| Swagger UI | http://localhost:8080/swagger/index.html |

### 5.2 Test Data
- Seed database before testing: `go run scripts/seed/main.go`
- Test credentials: `admin@example.com` / `password123`

---

## 6. Automation Approach & Results (Assignment 2)

### 6.1 Automation Strategy

**Approach:** Risk-based automation — HIGH-risk modules automated first using unit tests with mock repositories.

| Section | Details |
|---------|---------|
| **Automation Approach** | Risk-based, regression-focused; mocked dependencies for pure service logic testing |
| **Tool Selection** | Go `testing` + `testify/mock` for backend (lightweight, no extra frameworks needed); Playwright for frontend E2E (headless Chromium, simple API) |
| **Scope** | Auth service (9 tests), Product service (13 tests), Interaction service (9 tests); Frontend auth flows (9 tests), product/profile redirects (12 tests) |
| **Reusability** | Shared mock structs (`MockUserRepository`, `MockProductRepository`, `MockInteractionRepository`) reused across test files; Playwright tests organized by feature |

**Why Go testing + testify over alternatives:**
- Go built-in `testing` is zero-dependency and well-integrated with `go test` toolchain
- `testify/mock` provides clean mock generation matching repository interfaces
- Alternatives (gomock, ginkgo) add complexity without benefit for this scale

**Why Playwright over Selenium:**
- Native TypeScript support aligns with the React/TS frontend stack
- Modern auto-waiting eliminates flaky test issues common in Selenium
- Built-in HTML reporter and artifact support in GitHub Actions
- Simpler setup: no WebDriver binaries needed

### 6.2 Quality Gate Definitions

See full document: [QUALITY_GATE_REPORT.md](./QUALITY_GATE_REPORT.md)

| QG ID | Metric | Threshold | Status |
|-------|--------|-----------|--------|
| QG01 | Unit test pass rate | 100% | ✅ 31/31 |
| QG02 | Critical defects | 0 | ✅ 0 found |
| QG03 | Build success | 100% | ✅ Passing |
| QG04 | E2E test pass rate | 100% | ⏳ Pending CI |
| QG05 | Lint violations | 0 errors | ✅ Clean |
| QG06 | Unit test execution time | ≤ 2 min | ✅ ~0.7s |

### 6.3 Initial Results & Coverage Metrics

See full document: [METRICS_REPORT.md](./METRICS_REPORT.md)

| Module | Automated? | Coverage % | Tests |
|--------|------------|------------|-------|
| Auth | Yes | 100% | 9 unit tests |
| Product CRUD | Yes | 80% | 13 unit tests |
| Purchase Flow | Yes | 100% | 4 unit tests |
| Interactions | Yes | 100% | 11 unit tests |
| Recommendation Engine | No | 0% | Planned next |
| Frontend Auth | Yes | 100% | 9 E2E tests |
| Frontend Routes | Yes | 100% | 12 E2E tests |

**Overall: 52 automated tests, 93.75% HIGH-risk function coverage**

---

## 7. Planned Metrics

| Metric | Baseline (Assignment 1) | Current (Assignment 2) | Target |
|--------|------------------------|----------------------|--------|
| API endpoint coverage | 0% | 28 endpoints in Postman | 100% |
| Service unit test coverage | 0% | 31 tests (HIGH-risk) | 80% all services |
| E2E scenario coverage | 4 tests | 21 tests | 5+ critical flows |
| Defects found (HIGH severity) | 0 | 0 new (5 pre-existing) | tracked |
| CI pipeline pass rate | Green (build only) | Green (build + tests) | ≥ 95% |

### 6.1 Defect Severity Classification

| Severity | Definition | Example |
|----------|-----------|---------|
| **Critical** | System unusable | Login completely broken |
| **High** | Major feature broken | Purchase doesn't decrement stock |
| **Medium** | Feature works with workaround | Filter returns wrong results |
| **Low** | Minor visual/UX issue | Button misaligned |

---

## 7. Roles & Responsibilities

| Role | Responsibility |
|------|---------------|
| QA Engineer | Write and maintain tests, report defects |
| Developer | Fix defects, maintain CI pipeline |

---

## 8. Entry & Exit Criteria

### Entry Criteria (start testing)
- Application running locally (backend + frontend + MongoDB)
- Database seeded with test data
- Postman collection imported

### Exit Criteria (Assignment 2 complete)
- 31 unit tests passing (`go test ./internal/service/... -v`)
- 21 E2E tests written and passing in Playwright
- CI/CD pipelines green in both repositories (build + tests)
- QUALITY_GATE_REPORT.md created with all 7 gates defined
- METRICS_REPORT.md created with coverage, TTE, defects tables
- TEST_STRATEGY.md updated with automation details
