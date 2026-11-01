# LearnX Backend – Auth, Profiles, and Striver A2Z DSA Roadmap

High-performance Go REST backend for LearnX. It reuses the existing Supabase Auth + PostgreSQL stack and adds a DSA roadmap, central problem catalogue, per-student progress, and analytics seeded from Striver's A2Z DSA Sheet (TakeUForward).

## Features

- **External OTP Service Integration:** Leverages the external OTP service (`https://otp-latest-0k73.onrender.com`) for sending and verifying 6-digit one-time passwords for signup and password reset.
- **Signup Fields & Validation:** Full Name, Email Address, Password, Confirm Password, and 6-digit Email OTP verification.
- **Login Fields & Authentication:** Email Address and Password with Supabase Auth token issuance.
- **Forgot Password Fields & Single-Use Reset:** Email Address, 6-digit Email OTP, New Password, and Confirm New Password with replay-protected reset tokens.
- **Student Onboarding Flow:** Dedicated onboarding endpoint capturing College/University, Course, Department, Academic Year, Preferred Language (C, C++, Java, Python, Go), Learning Goal, and DSA Level.
- **Complete Student Profile Model (14 Fields):** `id`, `full_name`, `email`, `college`, `course`, `department`, `current_year`, `preferred_language`, `learning_goal`, `dsa_level`, `avatar_url`, `onboarding_completed`, `created_at`, `updated_at`.
- **Sensitive Credential Encryption at Rest:** Passwords pending OTP verification are encrypted with AES-256-GCM.
- **Dedicated Folders:** Schemas in `schema/` and all test suites in `tests/`.
- **Striver A2Z DSA catalogue:** Official sections, problem names, order, and canonical practice URLs. No editorials or solutions are stored.
- **Student progress isolation:** Attempted / Solved / Needs Revision records are scoped to the authenticated student, including Supabase RLS policies.
- **Roadmap analytics:** Overall, topic-wise, and difficulty stats are calculated from database rows.

## Project Structure

```text
backend/
├── cmd/
│   ├── migrate/main.go       # Apply PostgreSQL schema
│   └── import-a2z/main.go    # Upsert Striver A2Z catalogue
├── data/
│   ├── strivers-a2z-sheet.json         # Full official metadata seed
│   └── strivers-a2z-sheet.sample.json  # Truncated sample for tests
├── schema/
│   ├── 001_create_profiles.sql
│   ├── 002_create_pending_signups.sql
│   ├── 003_create_password_resets.sql
│   ├── 004_create_roadmaps.sql
│   ├── 005_create_topics.sql
│   ├── 006_create_subtopics.sql
│   ├── 007_create_problems.sql
│   ├── 008_create_student_progress.sql
│   ├── 009_rls_dsa.sql
│   └── schema.sql
├── tests/
├── internal/
│   ├── auth/
│   ├── roadmap/
│   ├── topics/
│   ├── problems/
│   ├── progress/
│   ├── analytics/
│   ├── importer/
│   ├── dsa/
│   └── server/
├── docs/
│   ├── API.md
│   └── openapi.yaml
└── main.go
```

## Getting Started

### 1. Environment Configuration

Copy `.env.example` to `.env`:

```bash
cp .env.example .env
```

### 2. Run Tests

All unit tests and end-to-end HTTP integration tests run from the `tests/` folder:

```bash
go test -v ./tests/...
```

### 3. Run the Server

```bash
go run main.go
```

The server starts on port `8080`. Without `SUPABASE_DB_URL` it uses an in-memory store and loads `data/strivers-a2z-sheet.json` when present.

```bash
curl http://localhost:8080/health
```

### 4. Import Striver's A2Z DSA Sheet

Apply schema (also runs on server start when the database is reachable):

```bash
go run ./cmd/migrate
```

Import or re-import the official catalogue. The command runs in a transaction, upserts by `(source_sheet, source_problem_id)`, and never deletes student progress:

```bash
go run ./cmd/import-a2z -file data/strivers-a2z-sheet.json
```

Source attribution is stored on the roadmap (`TakeUForward` + official sheet URL). Refresh the seed from TakeUForward's public syllabus if the sheet changes; do not scrape paywalled editorials.

### 5. Example authenticated calls

```bash
TOKEN="<supabase_access_token>"

curl -H "Authorization: Bearer $TOKEN" http://localhost:8080/api/v1/roadmaps
curl -H "Authorization: Bearer $TOKEN" "http://localhost:8080/api/v1/problems?difficulty=Easy&page=1&limit=20&sort=problem_order"
curl -H "Authorization: Bearer $TOKEN" "http://localhost:8080/api/v1/problems/search?q=two-sum"
curl -X POST -H "Authorization: Bearer $TOKEN" http://localhost:8080/api/v1/student/progress/problems/<problem_id>/attempt
curl -X PATCH -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" \
  -d '{"status":"Solved"}' \
  http://localhost:8080/api/v1/student/progress/problems/<problem_id>
curl -H "Authorization: Bearer $TOKEN" http://localhost:8080/api/v1/student/analytics/overview
```

### 6. API Documentation

- [`docs/API.md`](docs/API.md) – request/response examples
- [`docs/openapi.yaml`](docs/openapi.yaml) – OpenAPI 3 spec

Catalogue HTTP APIs are read-only for students. Shared roadmap/problem rows are readable by authenticated users; only the service-role import path mutates them. Student progress is protected with JWT verification, parameterized SQL, pagination caps (`limit` max 100), and the existing IP rate limiter.
