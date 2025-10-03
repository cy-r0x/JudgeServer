# Judge Backend API Documentation

## Table of Contents
1. [Overview](#overview)
2. [Authentication](#authentication)
3. [User Roles](#user-roles)
4. [API Endpoints](#api-endpoints)
   - [Users](#users)
   - [Contests](#contests)
   - [Problems](#problems)
   - [Submissions](#submissions)
   - [Contest Problems](#contest-problems)
   - [Setter](#setter)
   - [Standings](#standings)
5. [Data Models](#data-models)
6. [Error Responses](#error-responses)

---

## Overview

This is the backend API for an online judge system supporting competitive programming contests. The API is built using Go and provides endpoints for managing users, contests, problems, submissions, and standings.

**Base URL:** `http://localhost:{HTTP_PORT}/api`

---

## Authentication

Most endpoints require authentication via JWT (JSON Web Token). The token should be included in the request headers:

```
Authorization: Bearer <access_token>
```

The JWT token is obtained through the login endpoint and contains user information including:
- User ID (`sub`)
- Username
- Full Name
- Role
- Allowed Contest ID
- Room Number and PC Number (optional)

Token expiration: **3 hours** from issuance

---

## User Roles

The system supports three roles with different access levels:

1. **Admin**: Full access to all resources
2. **Setter**: Can create and manage their own problems
3. **User**: Can participate in contests and submit solutions

### Special Authentication Middleware:
- **AuthEngine**: Special authentication for the judging engine to update submission results

---

## API Endpoints

### Users

#### 1. Register User (Create User)
**Endpoint:** `POST /api/user/register`  
**Authentication:** Required (Admin only)  
**Description:** Creates a new user account

**Request Body:**
```json
{
  "full_name": "John Doe",
  "username": "johndoe",
  "email": "john@example.com",
  "password": "securepassword",
  "role": "user",
  "room_no": "A101",
  "pc_no": 5,
  "allowed_contest": 1
}
```

**Response:** `201 Created`
```json
{
  "id": 1,
  "full_name": "John Doe",
  "username": "johndoe",
  "email": "john@example.com",
  "role": "user",
  "room_no": "A101",
  "pc_no": 5,
  "allowed_contest": 1,
  "created_at": "2025-10-02T10:00:00Z"
}
```

---

#### 2. Login
**Endpoint:** `POST /api/user/login`  
**Authentication:** Not required  
**Description:** Authenticates a user and returns a JWT token

**Request Body:**
```json
{
  "username": "johndoe",
  "password": "securepassword"
}
```

**Response:** `200 OK`
```json
{
  "sub": 1,
  "full_name": "John Doe",
  "username": "johndoe",
  "role": "user",
  "room_no": "A101",
  "pc_no": 5,
  "allowed_contest": 1,
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "exp": 1696252800,
  "iat": 1696241800
}
```

---

#### 3. Logout
**Endpoint:** `POST /api/user/logout`  
**Authentication:** Not required  
**Description:** Logs out the current user (TODO: Token blacklisting not implemented yet)

**Response:** `200 OK`

---

#### 4. Get Users by Contest
**Endpoint:** `GET /api/user/{contestId}`  
**Authentication:** Required (Admin only)  
**Description:** Retrieves all users registered for a specific contest

**Path Parameters:**
- `contestId` (integer): The contest ID

**Response:** `200 OK`
```json
[
  {
    "userId": 1,
    "full_name": "John Doe",
    "username": "johndoe"
  },
  {
    "userId": 2,
    "full_name": "Jane Smith",
    "username": "janesmith"
  }
]
```

---

### Contests

#### 1. List All Contests
**Endpoint:** `GET /api/contests`  
**Authentication:** Not required  
**Description:** Retrieves all contests with their status

**Response:** `200 OK`
```json
[
  {
    "id": 1,
    "title": "Spring Programming Contest 2025",
    "start_time": "2025-10-15T09:00:00Z",
    "duration_seconds": 10800,
    "status": "UPCOMING"
  },
  {
    "id": 2,
    "title": "Fall Programming Contest 2025",
    "start_time": "2025-09-20T09:00:00Z",
    "duration_seconds": 10800,
    "status": "ENDED"
  }
]
```

**Status Values:**
- `UPCOMING`: Contest hasn't started yet
- `RUNNING`: Contest is currently active
- `ENDED`: Contest has finished

---

#### 2. Get Contest Details
**Endpoint:** `GET /api/contests/{contestId}`  
**Authentication:** Not required  
**Description:** Retrieves detailed information about a specific contest, including problems (if contest has started)

**Path Parameters:**
- `contestId` (integer): The contest ID

**Response:** `200 OK`
```json
{
  "contest": {
    "id": 1,
    "title": "Spring Programming Contest 2025",
    "description": "Annual spring contest for all participants",
    "start_time": "2025-10-15T09:00:00Z",
    "duration_seconds": 10800,
    "status": "RUNNING",
    "created_at": "2025-09-01T10:00:00Z"
  },
  "problems": [
    {
      "id": 101,
      "title": "Two Sum",
      "slug": "two-sum",
      "index": 1
    },
    {
      "id": 102,
      "title": "Binary Search",
      "slug": "binary-search",
      "index": 2
    }
  ]
}
```

**Note:** Problems array is empty if contest status is `UPCOMING`

---

#### 3. Create Contest
**Endpoint:** `POST /api/contests`  
**Authentication:** Required (Admin only)  
**Description:** Creates a new contest

**Request Body:**
```json
{
  "title": "Spring Programming Contest 2025",
  "description": "Annual spring contest for all participants",
  "start_time": "2025-10-15T09:00:00Z",
  "duration_seconds": 10800
}
```

**Response:** `201 Created`
```json
{
  "id": 1,
  "title": "Spring Programming Contest 2025",
  "description": "Annual spring contest for all participants",
  "start_time": "2025-10-15T09:00:00Z",
  "duration_seconds": 10800,
  "created_at": "2025-10-02T10:00:00Z"
}
```

---

#### 4. Update Contest
**Endpoint:** `PUT /api/contests`  
**Authentication:** Required (Admin only)  
**Description:** Updates an existing contest

**Request Body:**
```json
{
  "id": 1,
  "title": "Updated Spring Programming Contest 2025",
  "description": "Updated description",
  "start_time": "2025-10-15T10:00:00Z",
  "duration_seconds": 12600
}
```

**Response:** `200 OK`
```json
{
  "id": 1,
  "title": "Updated Spring Programming Contest 2025",
  "description": "Updated description",
  "start_time": "2025-10-15T10:00:00Z",
  "duration_seconds": 12600,
  "created_at": "2025-10-02T10:00:00Z"
}
```

---

### Problems

#### 1. Get Problem
**Endpoint:** `GET /api/problems/{problemId}`  
**Authentication:** Required  
**Description:** Retrieves problem details with testcases

**Access Control:**
- **Users**: Can only view problems from their allowed contest, see only sample testcases
- **Setters**: Can only view problems they created, see all testcases
- **Admins**: Can view all problems with all testcases

**Path Parameters:**
- `problemId` (integer): The problem ID

**Response:** `200 OK`
```json
{
  "id": 101,
  "title": "Two Sum",
  "slug": "two-sum",
  "statement": "Given an array of integers nums and an integer target...",
  "input_statement": "First line contains n and target...",
  "output_statement": "Print the indices of the two numbers...",
  "time_limit": 1.0,
  "memory_limit": 256.0,
  "created_by": 5,
  "created_at": "2025-09-01T10:00:00Z",
  "test_cases": [
    {
      "id": 1,
      "problem_id": 101,
      "input": "4 9\n2 7 11 15",
      "expected_output": "0 1",
      "is_sample": true,
      "created_at": "2025-09-01T10:05:00Z"
    }
  ]
}
```

---

#### 2. Create Problem
**Endpoint:** `POST /api/problems`  
**Authentication:** Required (Setter only)  
**Description:** Creates a new problem (without testcases initially)

**Request Body:**
```json
{
  "title": "Two Sum"
}
```

**Response:** `201 Created`
```json
{
  "id": 101,
  "title": "Two Sum",
  "slug": "two-sum",
  "statement": "",
  "input_statement": "",
  "output_statement": "",
  "time_limit": 1.0,
  "memory_limit": 256.0,
  "created_by": 5,
  "created_at": "2025-10-02T10:00:00Z",
  "test_cases": null
}
```

---

#### 3. Update Problem
**Endpoint:** `PUT /api/problems`  
**Authentication:** Required (Setter only - must be problem creator)  
**Description:** Updates problem details and testcases

**Request Body:**
```json
{
  "id": 101,
  "title": "Two Sum Problem",
  "statement": "Given an array of integers nums and an integer target, return indices of the two numbers such that they add up to target.",
  "input_statement": "First line contains n (array size) and target. Second line contains n integers.",
  "output_statement": "Print two space-separated indices (0-indexed).",
  "time_limit": 2.0,
  "memory_limit": 512.0,
  "test_cases": [
    {
      "input": "4 9\n2 7 11 15",
      "expected_output": "0 1",
      "is_sample": true
    },
    {
      "input": "3 6\n3 2 4",
      "expected_output": "1 2",
      "is_sample": false
    }
  ]
}
```

**Response:** `200 OK`
Returns the updated problem details (same format as Get Problem)

---

### Submissions

#### 1. Create Submission
**Endpoint:** `POST /api/submissions`  
**Authentication:** Required (User)  
**Description:** Submits a solution to a problem in a contest

**Request Body:**
```json
{
  "problem_id": 101,
  "contest_id": 1,
  "language": "cpp",
  "source_code": "#include <iostream>\nusing namespace std;\n\nint main() {\n    // solution code\n    return 0;\n}"
}
```

**Response:** `200 OK`
```json
{
  "submission_id": 501
}
```

---

#### 2. Get Submission
**Endpoint:** `GET /api/submissions/{submissionId}`  
**Authentication:** Required (User can only view their own submissions)  
**Description:** Retrieves details of a specific submission

**Path Parameters:**
- `submissionId` (integer): The submission ID

**Response:** `200 OK`
```json
{
  "id": 501,
  "user_id": 1,
  "username": "johndoe",
  "problem_id": 101,
  "contest_id": 1,
  "language": "cpp",
  "source_code": "#include <iostream>...",
  "verdict": "ACCEPTED",
  "execution_time": 0.45,
  "memory_used": 2048.0,
  "submitted_at": "2025-10-15T10:30:00Z"
}
```

**Verdict Values:**
- `PENDING`: Submission is queued for evaluation
- `ACCEPTED`: All test cases passed
- `WRONG_ANSWER`: Solution produced incorrect output
- `TIME_LIMIT_EXCEEDED`: Solution exceeded time limit
- `MEMORY_LIMIT_EXCEEDED`: Solution exceeded memory limit
- `RUNTIME_ERROR`: Solution crashed during execution
- `COMPILATION_ERROR`: Code failed to compile

---

#### 3. List User Submissions
**Endpoint:** `GET /api/submissions`  
**Authentication:** Required (User)  
**Description:** Lists all submissions by the authenticated user for their allowed contest

**Response:** `200 OK`
```json
[
  {
    "id": 501,
    "user_id": 1,
    "username": "johndoe",
    "problem_id": 101,
    "contest_id": 1,
    "language": "cpp",
    "source_code": "#include <iostream>...",
    "verdict": "ACCEPTED",
    "execution_time": 0.45,
    "memory_used": 2048.0,
    "submitted_at": "2025-10-15T10:30:00Z"
  }
]
```

---

#### 4. List All Submissions for Contest
**Endpoint:** `GET /api/submissions/all/{contestId}`  
**Authentication:** Required  
**Description:** Lists all submissions for a specific contest

**Path Parameters:**
- `contestId` (integer): The contest ID

**Response:** `200 OK`
```json
[
  {
    "id": 501,
    "user_id": 1,
    "username": "johndoe",
    "problem_id": 101,
    "contest_id": 1,
    "language": "cpp",
    "source_code": "#include <iostream>...",
    "verdict": "ACCEPTED",
    "execution_time": 0.45,
    "memory_used": 2048.0,
    "submitted_at": "2025-10-15T10:30:00Z"
  }
]
```

---

#### 5. Update Submission (Engine Only)
**Endpoint:** `PUT /api/submissions`  
**Authentication:** Required (Engine authentication with ENGINE_KEY)  
**Description:** Updates submission verdict after evaluation by the judging engine

**Request Body:**
```json
{
  "submission_id": 501,
  "verdict": "ACCEPTED",
  "execution_time": 0.45,
  "execution_memory": 2048.0
}
```

**Response:** `200 OK`
```json
{
  "message": "Submission updated"
}
```

**Note:** This endpoint automatically updates contest standings when the verdict is ACCEPTED.

---

### Contest Problems

#### 1. Get Contest Problems
**Endpoint:** `GET /api/contests/problems/{contestId}`  
**Authentication:** Required  
**Description:** Retrieves all problems assigned to a contest with problem details

**Path Parameters:**
- `contestId` (integer): The contest ID

**Response:** `200 OK`
```json
[
  {
    "contest_id": 1,
    "problem_id": 101,
    "index": 1,
    "problem_name": "Two Sum",
    "problem_author": "Alice Johnson"
  },
  {
    "contest_id": 1,
    "problem_id": 102,
    "index": 2,
    "problem_name": "Binary Search",
    "problem_author": "Bob Smith"
  }
]
```

---

#### 2. Assign Problem to Contest
**Endpoint:** `POST /api/contests/assign`  
**Authentication:** Required (Admin only)  
**Description:** Assigns a problem to a contest

**Request Body:**
```json
{
  "contest_id": 1,
  "problem_id": 101
}
```

**Response:** `200 OK`
```json
{
  "contest_id": 1,
  "problem_id": 101,
  "index": 3,
  "problem_name": "Two Sum",
  "problem_author": "Alice Johnson"
}
```

**Note:** The index is automatically assigned as the next available position.

---

#### 3. Update Problem Index
**Endpoint:** `PUT /api/contests/index`  
**Authentication:** Required (Admin only)  
**Description:** Updates the ordering/index of problems in a contest

**Request Body:**
```json
[
  {
    "contest_id": 1,
    "problem_id": 102,
    "index": 1
  },
  {
    "contest_id": 1,
    "problem_id": 101,
    "index": 2
  }
]
```

**Response:** `200 OK`
```json
[
  {
    "contest_id": 1,
    "problem_id": 102,
    "index": 1
  },
  {
    "contest_id": 1,
    "problem_id": 101,
    "index": 2
  }
]
```

---

### Setter

#### 1. List Setter Problems
**Endpoint:** `GET /api/setter`  
**Authentication:** Required (Setter only)  
**Description:** Lists all problems created by the authenticated setter

**Response:** `200 OK`
```json
[
  {
    "id": 101,
    "title": "Two Sum",
    "created_at": "2025-09-01T10:00:00Z"
  },
  {
    "id": 105,
    "title": "Linked List Cycle",
    "created_at": "2025-09-15T14:30:00Z"
  }
]
```

---

### Standings

#### 1. Get Contest Standings
**Endpoint:** `GET /api/standings/{contestId}`  
**Authentication:** Not required  
**Description:** Retrieves the leaderboard/standings for a contest

**Path Parameters:**
- `contestId` (integer): The contest ID

**Response:** `200 OK`
```json
[
  {
    "contest_id": 1,
    "user_id": 1,
    "username": "johndoe",
    "penalty": 45,
    "solved_count": 3,
    "last_solved_at": "2025-10-15T11:45:00Z",
    "solved": [
      {
        "problem_id": 101,
        "problem_index": 1,
        "solved_at": "2025-10-15T10:30:00Z",
        "penalty": 30
      },
      {
        "problem_id": 102,
        "problem_index": 2,
        "solved_at": "2025-10-15T11:00:00Z",
        "penalty": 60
      },
      {
        "problem_id": 103,
        "problem_index": 3,
        "solved_at": "2025-10-15T11:45:00Z",
        "penalty": 105
      }
    ]
  },
  {
    "contest_id": 1,
    "user_id": 2,
    "username": "janesmith",
    "penalty": 120,
    "solved_count": 2,
    "last_solved_at": "2025-10-15T12:00:00Z",
    "solved": [
      {
        "problem_id": 101,
        "problem_index": 1,
        "solved_at": "2025-10-15T11:00:00Z",
        "penalty": 60
      },
      {
        "problem_id": 102,
        "problem_index": 2,
        "solved_at": "2025-10-15T12:00:00Z",
        "penalty": 120
      }
    ]
  }
]
```

**Ranking Logic:**
1. Primary: More problems solved (descending)
2. Secondary: Lower penalty time (ascending)
3. Tertiary: Earlier last solve time (ascending)
4. Final: Lower user ID (ascending)

---

## Data Models

### User
```json
{
  "id": "integer",
  "full_name": "string",
  "username": "string (unique)",
  "email": "string",
  "password": "string (bcrypt hashed)",
  "role": "string (user|setter|admin)",
  "room_no": "string (optional)",
  "pc_no": "integer (optional)",
  "allowed_contest": "integer (optional, foreign key to contests)",
  "created_at": "timestamp"
}
```

### Contest
```json
{
  "id": "integer",
  "title": "string",
  "description": "string",
  "start_time": "timestamp",
  "duration_seconds": "integer",
  "status": "string (UPCOMING|RUNNING|ENDED) - computed",
  "created_at": "timestamp"
}
```

### Problem
```json
{
  "id": "integer",
  "title": "string",
  "slug": "string (auto-generated from title)",
  "statement": "string (problem description)",
  "input_statement": "string (input format)",
  "output_statement": "string (output format)",
  "time_limit": "float32 (seconds)",
  "memory_limit": "float32 (MB)",
  "created_by": "integer (foreign key to users)",
  "created_at": "timestamp",
  "test_cases": "array of Testcase objects"
}
```

### Testcase
```json
{
  "id": "integer",
  "problem_id": "integer (foreign key)",
  "input": "string",
  "expected_output": "string",
  "is_sample": "boolean",
  "created_at": "timestamp"
}
```

### Submission
```json
{
  "id": "integer",
  "user_id": "integer (foreign key to users)",
  "username": "string",
  "problem_id": "integer (foreign key to problems)",
  "contest_id": "integer (foreign key to contests)",
  "language": "string (cpp|java|python|c|etc)",
  "source_code": "string",
  "verdict": "string",
  "execution_time": "float32 (seconds, nullable)",
  "memory_used": "float32 (KB, nullable)",
  "submitted_at": "timestamp"
}
```

### Contest Standing
```json
{
  "contest_id": "integer (foreign key)",
  "user_id": "integer (foreign key)",
  "username": "string",
  "penalty": "integer (total penalty time in minutes)",
  "solved_count": "integer",
  "last_solved_at": "timestamp (nullable)",
  "solved": "array of Solve objects"
}
```

### Solve
```json
{
  "problem_id": "integer",
  "problem_index": "integer",
  "solved_at": "timestamp",
  "penalty": "integer (penalty in minutes)"
}
```

---

## Error Responses

All error responses follow this general format:

### 400 Bad Request
```json
{
  "error": "Invalid request payload"
}
```
or
```json
"Error message string"
```

### 401 Unauthorized
```json
{
  "error": "User information not found"
}
```
or
```json
"Invalid Token"
```

### 403 Forbidden
```json
{
  "error": "You don't have access to this resource"
}
```

### 404 Not Found
```json
{
  "error": "Resource not found"
}
```

### 409 Conflict
```json
{
  "error": "Problem already assigned to this contest"
}
```

### 500 Internal Server Error
```json
{
  "error": "Internal Server Error"
}
```

---

## Environment Variables

The application requires the following environment variables (`.env` file):

```env
# Server Configuration
HTTP_PORT=8080

# JWT Authentication
JWT_SECRET=your-secret-key-here

# Judge Engine Authentication
ENGINE_KEY=your-engine-key-here

# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_NAME=judge_db
DB_USER=postgres
DB_PASSWORD=your-db-password
DB_SSL_MODE=disable
```

---

## Notes

1. **Password Security**: All passwords are hashed using bcrypt before storage
2. **Slug Generation**: Problem slugs are automatically generated from titles (lowercase, spaces replaced with hyphens)
3. **Contest Status**: Contest status is computed dynamically based on current time, start time, and duration
4. **Testcase Visibility**: Users only see sample testcases, while setters and admins see all testcases
5. **Queue System**: Submissions are automatically sent to a queue for asynchronous evaluation by the judging engine
6. **Standings Updates**: Contest standings are automatically updated when accepted submissions are recorded
7. **Timezone**: All timestamps are stored and returned in UTC

---

## Database Schema

The application uses PostgreSQL with the following main tables:
- `users`
- `contests`
- `problems`
- `testcases`
- `submissions`
- `contest_problems` (junction table)
- `contest_standings`
- `contest_solves`

Migration files are located in the `schema/` directory.

---

## Future Enhancements (TODO)

- Token blacklisting for logout functionality
- Additional user management endpoints
- Problem difficulty levels
- Contest registration system
- Real-time updates via WebSocket
- Plagiarism detection
- Editorial/solution viewing after contest
- Discussion forums
- User rating system

---

**Last Updated:** October 2, 2025  
**API Version:** 1.0
