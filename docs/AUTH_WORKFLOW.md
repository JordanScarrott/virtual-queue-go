# Authentication Workflow

[← Back to README](../../README.md)

This guide explains how to manually test the authentication system using `curl`.

**Base URL:** `http://localhost:2015`

## 1. The Admin Flow (Magic Link/OTP)

This flow mimics a business owner logging in. It uses a Temporal Workflow to manage the OTP process.

### Step 1: Start Login
Initiate the login process. This starts a `LoginWorkflow`.

```bash
curl -X POST http://localhost:2015/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "admin@redduck.app"}'
```

**Response:**
```json
{
  "runID": "...",
  "workflowID": "auth-admin@redduck.app"
}
```

### Step 2: Find the Code
In a real app, this code would be emailed. For development, check the worker logs.

```bash
# If running via Docker Compose
docker logs backend_worker | grep "MAGIC LOGIN CODE"

# If running locally
# Look at your terminal output for:
# INFO  MAGIC LOGIN CODE FOR admin@redduck.app: 123456
```

### Step 3: Verify & Get Token
Submit the code to complete the workflow.

```bash
curl -X POST http://localhost:2015/auth/verify \
  -H "Content-Type: application/json" \
  -d '{"email": "admin@redduck.app", "code": "123456"}'
```

**Response:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIs..."
}
```
*Copy this token.*

### Step 4: Use the Token
Access a protected route using the token.

```bash
curl -i http://localhost:2015/req/me \
  -H "Authorization: Bearer <YOUR_TOKEN>"
```

**Response:**
```json
{
  "message": "You are authenticated!"
}
```

---

## 2. The Guest Flow (Customer)

Customers joining a queue do not need to log in via email. They use "Guest Mode".

### Step 1: Join as Guest
Send a request with an empty `user_id`. The system will generate one for you.

```bash
curl -X POST http://localhost:2015/queues/join \
  -H "Content-Type: application/json" \
  -d '{"business_id": "barbershop-1", "user_id": ""}'
```

### Step 2: Note the ID
The response includes a `user_id`. The client app (Flutter) must save this ID to local storage.

**Response:**
```json
{
  "run_id": "...",
  "user_id": "d1e3d0a8-...", 
  "workflow_id": "..."
}
```
