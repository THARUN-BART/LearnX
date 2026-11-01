# LearnX Backend – Supabase Authentication + Custom OTP Integration

High-performance, secure Go authentication backend for LearnX integrating Supabase Auth, PostgreSQL, and an external OTP service.

## Features

- **External OTP Service Integration:** Leverages the external OTP service (`https://otp-latest-0k73.onrender.com`) for sending and verifying 6-digit one-time passwords for signup and password reset.
- **Signup Fields & Validation:** Full Name, Email Address, Password, Confirm Password, and 6-digit Email OTP verification.
- **Login Fields & Authentication:** Email Address and Password with Supabase Auth token issuance.
- **Forgot Password Fields & Single-Use Reset:** Email Address, 6-digit Email OTP, New Password, and Confirm New Password with replay-protected reset tokens.
- **Student Onboarding Flow:** Dedicated onboarding endpoint capturing College/University, Course, Department, Academic Year, Preferred Language (C, C++, Java, Python, Go), Learning Goal, and DSA Level.
- **Complete Student Profile Model (14 Fields):** `id`, `full_name`, `email`, `college`, `course`, `department`, `current_year`, `preferred_language`, `learning_goal`, `dsa_level`, `avatar_url`, `onboarding_completed`, `created_at`, `updated_at`.
- **Sensitive Credential Encryption at Rest:** Passwords pending OTP verification are encrypted with AES-256-GCM.
- **Dedicated Folders:** Schemas in `schema/` and all test suites in `tests/`.

## Project Structure

```text
LearnX/
├── cmd/
│   └── server/
│       └── main.go           # Server entry point
├── schema/                   # Dedicated SQL schema directory
│   ├── 001_create_profiles.sql
│   ├── 002_create_pending_signups.sql
│   ├── 003_create_password_resets.sql
│   └── schema.sql            # Complete unified schema
├── tests/                    # Dedicated test suites directory
│   ├── auth_test.go          # Signup, password match, login, forgot pass
│   ├── crypto_test.go        # AES-256-GCM, SHA-256, tokens
│   ├── onboarding_test.go    # Student onboarding validations & fields
│   ├── otp_test.go           # External OTP client & cooldown service
│   ├── profile_test.go       # Student profile retrieval & update
│   ├── ratelimit_test.go     # Token bucket IP rate limiter
│   └── router_test.go        # End-to-end HTTP integration test
├── internal/
│   ├── auth/                 # Signup, login, password reset service & handler
│   ├── config/               # Environment & configuration loading
│   ├── crypto/               # AES-256-GCM encryption, hashing, tokens
│   ├── database/             # PostgreSQL pgxpool & schema migration
│   ├── middleware/           # Auth JWT, rate limit, CORS, recovery, logging
│   ├── otp/                  # External OTP client, rate limiting & models
│   ├── profile/              # Student profile & onboarding service & handler
│   ├── server/               # Router configuration & HTTP routes
│   └── supabase/             # Supabase GoTrue Admin & Auth client
├── docs/
│   └── API.md                # Full API specification & field definitions
├── .env.example              # Environment variables template
├── go.mod
├── go.sum
└── main.go                   # Root entry point
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
# or
go run cmd/server/main.go
```

The server will start on port `8080`.

Test the health check:
```bash
curl http://localhost:8080/health
```

### 4. API Documentation

Detailed endpoint schemas, request/response payloads, and security guidelines are documented in [`docs/API.md`](docs/API.md).
