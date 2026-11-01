# LearnX Backend Authentication & Student Profile API Specification

This document details the complete API schema, authentication lifecycle, student onboarding, and profile models for the LearnX Backend service, integrating Go, Supabase Auth, PostgreSQL, and the external OTP verification service.

---

## Architecture Overview

```
                      +---------------------------------------+
                      |       Client (Web / Mobile)           |
                      +---------------------------------------+
                             |                     ^
             HTTP Requests   |                     |  JWT Tokens / Profiles
                             v                     |
                      +---------------------------------------+
                      |         LearnX Go Backend             |
                      +---------------------------------------+
                         /         |              \
           AES-256-GCM  /          | pgxpool       \ HTTP REST
           Encryption  /           |                \
                      v            v                 v
           +--------------+  +---------------+  +---------------------+
           | Memory/DB    |  | Supabase      |  | External OTP        |
           | Pending Store|  | Postgres DB   |  | Service             |
           +--------------+  +---------------+  | (otp-latest-0k73)   |
                                   |            +---------------------+
                                   v
                             +---------------+
                             | Supabase Auth |
                             | (GoTrue Admin)|
                             +---------------+
```

---

## 1. Authentication Endpoints

### 1.1 Initiate Signup (Step 1)
Validates full name, email, password, and confirmation password. Encrypts sensitive credentials with AES-256-GCM, stores the pending registration, and dispatches a 6-digit OTP code to the student's email via the external OTP service.

- **Method:** `POST`
- **URL:** `/api/v1/auth/signup`
- **Auth Required:** No
- **Fields:**
  - `full_name`: Full Name (required, min 2 chars)
  - `email`: Email Address (required, valid RFC 5322 format)
  - `password`: Password (required, min 8 chars, >=1 uppercase, >=1 lowercase, >=1 digit)
  - `confirm_password`: Confirm Password (must match `password`)

#### Request Body
```json
{
  "full_name": "Student Name",
  "email": "student@gmail.com",
  "password": "StrongPassword123",
  "confirm_password": "StrongPassword123"
}
```

#### Success Response
- **Status:** `200 OK`
```json
{
  "success": true,
  "message": "OTP sent to your email. Please verify to complete signup.",
  "email": "student@gmail.com",
  "expires_in_seconds": 600
}
```

#### Error Responses
- **Status:** `400 Bad Request` (Password mismatch)
```json
{
  "success": false,
  "error": "password and confirmation password do not match"
}
```
- **Status:** `409 Conflict` (Duplicate account)
```json
{
  "success": false,
  "error": "an account with this email already exists"
}
```
- **Status:** `429 Too Many Requests` (Resend cooldown active)
```json
{
  "success": false,
  "error": "OTP resend cooldown is active, please wait before requesting again"
}
```

---

### 1.2 Verify Signup OTP & Create Account (Step 2)
Verifies the 6-digit email OTP through the external OTP service. Creates the user in Supabase Auth using trusted server-side email verification (`email_confirm: true`). Creates the initial student profile in Supabase PostgreSQL, wipes pending credentials, and issues session tokens.

- **Method:** `POST`
- **URL:** `/api/v1/auth/verify-signup-otp`
- **Auth Required:** No
- **Fields:**
  - `email`: Email Address
  - `code`: Email OTP (6-digit verification)

#### Request Body
```json
{
  "email": "student@gmail.com",
  "code": "123456"
}
```

#### Success Response
- **Status:** `201 Created`
```json
{
  "success": true,
  "message": "Signup and OTP verification completed successfully.",
  "user": {
    "id": "e6a2b847-19fa-47ea-bd57-5ffc1d766324",
    "email": "student@gmail.com",
    "full_name": "Student Name",
    "role": "student",
    "onboarding_completed": false
  },
  "session": {
    "access_token": "eyJhbGciOi...",
    "refresh_token": "8gD3vR...",
    "token_type": "bearer",
    "expires_in": 3600
  }
}
```

---

### 1.3 Login (Step 3)
Authenticates existing users with Supabase Auth via email and password. Returns user summary and access/refresh tokens.

- **Method:** `POST`
- **URL:** `/api/v1/auth/login`
- **Auth Required:** No
- **Fields:**
  - `email`: Email Address
  - `password`: Password

#### Request Body
```json
{
  "email": "student@gmail.com",
  "password": "StrongPassword123"
}
```

#### Success Response
- **Status:** `200 OK`
```json
{
  "success": true,
  "message": "Login successful.",
  "user": {
    "id": "e6a2b847-19fa-47ea-bd57-5ffc1d766324",
    "email": "student@gmail.com",
    "full_name": "Student Name",
    "role": "student",
    "onboarding_completed": false
  },
  "session": {
    "access_token": "eyJhbGciOi...",
    "refresh_token": "8gD3vR...",
    "token_type": "bearer",
    "expires_in": 3600
  }
}
```

---

### 1.4 Forgot Password (Step 4)
Dispatches a 6-digit OTP via the external OTP service if the account exists. Returns a generic response to prevent email enumeration attacks.

- **Method:** `POST`
- **URL:** `/api/v1/auth/forgot-password`
- **Auth Required:** No
- **Fields:**
  - `email`: Email Address

#### Request Body
```json
{
  "email": "student@gmail.com"
}
```

#### Success Response
- **Status:** `200 OK`
```json
{
  "success": true,
  "message": "If an account with that email exists, a password reset OTP has been sent."
}
```

---

### 1.5 Verify Password Reset OTP (Step 5)
Verifies the 6-digit email OTP through the external OTP service. On success, issues a single-use, short-lived reset authorization token (`reset_token`).

- **Method:** `POST`
- **URL:** `/api/v1/auth/verify-reset-otp`
- **Auth Required:** No
- **Fields:**
  - `email`: Email Address
  - `code`: Email OTP (6-digit verification)

#### Request Body
```json
{
  "email": "student@gmail.com",
  "code": "123456"
}
```

#### Success Response
- **Status:** `200 OK`
```json
{
  "success": true,
  "message": "OTP verified successfully. You may now reset your password.",
  "reset_token": "63f13f0426b3c95e1e75db740e716ce0368a41764ebc402120cb95df66dfdfc2",
  "expires_in_seconds": 900
}
```

---

### 1.6 Reset Password (Step 6)
Resets user password using the single-use reset authorization token. Validates `new_password` and `confirm_new_password`. Atomically consumes the token to prevent replay.

- **Method:** `POST`
- **URL:** `/api/v1/auth/reset-password`
- **Auth Required:** No
- **Fields:**
  - `reset_token`: Reset authorization token (from Step 5)
  - `new_password`: New Password
  - `confirm_new_password`: Confirm New Password (must match `new_password`)

#### Request Body
```json
{
  "reset_token": "63f13f0426b3c95e1e75db740e716ce0368a41764ebc402120cb95df66dfdfc2",
  "new_password": "NewStrongPassword123",
  "confirm_new_password": "NewStrongPassword123"
}
```

#### Success Response
- **Status:** `200 OK`
```json
{
  "success": true,
  "message": "Password has been reset successfully. Please log in with your new password."
}
```

---

### 1.7 Logout
Revokes the current session token in Supabase Auth.

- **Method:** `POST`
- **URL:** `/api/v1/auth/logout`
- **Auth Required:** Yes (`Authorization: Bearer <access_token>`)

---

### 1.8 Google OAuth
Initiates Google OAuth authentication flow.

- **Method:** `GET`
- **URL:** `/api/v1/auth/google?redirect_to=http://localhost:3000/auth/callback`
- **Auth Required:** No

---

## 2. Student Onboarding & Profile Endpoints

### 2.1 Complete Student Onboarding
Collects student onboarding details and sets `onboarding_completed: true`.

- **Method:** `POST`
- **URL:** `/api/v1/profile/onboarding` (also available via `/api/v1/onboarding`)
- **Auth Required:** Yes (`Authorization: Bearer <access_token>`)
- **Fields:**
  - `college`: College / University Name (required)
  - `course`: Course (required, e.g. B.Tech, B.E, B.Sc, BCA, M.Tech, MCA)
  - `department`: Department (required, e.g. CSE, IT, ECE)
  - `current_year`: Current Academic Year (required, e.g. "1st Year", "2nd Year", "3rd Year", "4th Year")
  - `preferred_language`: Preferred Programming Language (required, allowed: `C`, `C++`, `Java`, `Python`, `Go`)
  - `learning_goal`: Learning Goal (required, e.g. `DSA`, `Interview Preparation`, `Competitive Programming`)
  - `dsa_level`: Current DSA Level (required, allowed: `Beginner`, `Intermediate`, `Advanced`)

#### Request Body
```json
{
  "college": "Stanford University",
  "course": "B.Tech",
  "department": "Computer Science and Engineering",
  "current_year": "3rd Year",
  "preferred_language": "Go",
  "learning_goal": "Interview Preparation",
  "dsa_level": "Intermediate"
}
```

#### Success Response
- **Status:** `200 OK`
```json
{
  "success": true,
  "message": "student onboarding completed successfully",
  "profile": {
    "id": "e6a2b847-19fa-47ea-bd57-5ffc1d766324",
    "full_name": "Student Name",
    "email": "student@gmail.com",
    "college": "Stanford University",
    "course": "B.Tech",
    "department": "Computer Science and Engineering",
    "current_year": "3rd Year",
    "preferred_language": "Go",
    "learning_goal": "Interview Preparation",
    "dsa_level": "Intermediate",
    "avatar_url": "",
    "onboarding_completed": true,
    "created_at": "2026-11-01T10:00:00Z",
    "updated_at": "2026-11-01T10:05:00Z"
  }
}
```

---

### 2.2 Get Student Profile
Retrieves the full student profile including all 14 fields.

- **Method:** `GET`
- **URL:** `/api/v1/profile`
- **Auth Required:** Yes (`Authorization: Bearer <access_token>`)

#### Success Response
- **Status:** `200 OK`
```json
{
  "success": true,
  "profile": {
    "id": "e6a2b847-19fa-47ea-bd57-5ffc1d766324",
    "full_name": "Student Name",
    "email": "student@gmail.com",
    "college": "Stanford University",
    "course": "B.Tech",
    "department": "Computer Science and Engineering",
    "current_year": "3rd Year",
    "preferred_language": "Go",
    "learning_goal": "Interview Preparation",
    "dsa_level": "Intermediate",
    "avatar_url": "https://cdn.learnx.com/avatars/user.png",
    "onboarding_completed": true,
    "created_at": "2026-11-01T10:00:00Z",
    "updated_at": "2026-11-01T10:05:00Z"
  }
}
```

---

### 2.3 Update Student Profile
Allows students to update their profile information.

- **Method:** `PUT`
- **URL:** `/api/v1/profile`
- **Auth Required:** Yes (`Authorization: Bearer <access_token>`)

#### Request Body
```json
{
  "full_name": "Updated Student Name",
  "avatar_url": "https://cdn.learnx.com/avatars/new.png",
  "dsa_level": "Advanced"
}
```

---

## 3. Student Profile Fields Reference (Supabase PostgreSQL)

| # | Field Name | Data Type | Description |
|---|---|---|---|
| 1 | `id` | `UUID PRIMARY KEY` | Supabase Auth user ID |
| 2 | `full_name` | `TEXT NOT NULL` | Student full name |
| 3 | `email` | `TEXT NOT NULL UNIQUE` | Registered email address |
| 4 | `college` | `TEXT NOT NULL DEFAULT ''` | College / University Name |
| 5 | `course` | `TEXT NOT NULL DEFAULT ''` | Course (B.Tech, B.E, B.Sc, etc.) |
| 6 | `department` | `TEXT NOT NULL DEFAULT ''` | Department (CSE, IT, ECE, etc.) |
| 7 | `current_year` | `TEXT NOT NULL DEFAULT ''` | Current Academic Year |
| 8 | `preferred_language` | `TEXT NOT NULL DEFAULT ''` | Preferred Programming Language (C, C++, Java, Python, Go) |
| 9 | `learning_goal` | `TEXT NOT NULL DEFAULT ''` | Learning Goal (DSA, Interview Prep, etc.) |
| 10 | `dsa_level` | `TEXT NOT NULL DEFAULT ''` | Current DSA Level (Beginner, Intermediate, Advanced) |
| 11 | `avatar_url` | `TEXT NOT NULL DEFAULT ''` | Student avatar URL |
| 12 | `onboarding_completed` | `BOOLEAN NOT NULL DEFAULT FALSE` | Flag indicating onboarding status |
| 13 | `created_at` | `TIMESTAMPTZ NOT NULL DEFAULT NOW()` | Record creation timestamp |
| 14 | `updated_at` | `TIMESTAMPTZ NOT NULL DEFAULT NOW()` | Record last updated timestamp |
