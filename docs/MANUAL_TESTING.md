# Manual Testing Cheat Sheet

[← Back to README](../../README.md)

Use this guide to manually verify the **Authentication**, **Queue Management**, and **Service Point Routing** features using `curl`.

**Assumptions:**
- Backend is running behind Caddy at `http://localhost:2015`.
- You have `jq` installed (optional, for pretty printing JSON).

---

## 0. Start Services

Ensure the following services are running in separate terminals:

1.  **Temporal Server**:
    ```bash
    temporal server start-dev
    ```

2.  **Go Server**:
    ```bash
    (cd virtual-queue-go && go run cmd/server/main.go)
    ```

3.  **Go Worker**:
    ```bash
    (cd virtual-queue-go && go run cmd/worker/main.go)
    ```

4.  **NATS Listener (Optional)**:
    Listen to all events (`>`) to verify messages are being published.
    ```bash
    nats sub ">"
    ```

---

## Setup: Create Queue

Before testing, ensure a queue exists.

**Command:**
```bash
curl -X POST "http://localhost:2015/create_queue?business_id=barbershop-1&queue_id=barbershop-1"
```

---

## A. The Customer Flow (Guest Mode)

### 1. Join Queue
Customers join purely anonymously. The system assigns a `user_id` which acts as their session.

**Command:**
```bash
curl -X POST "http://localhost:2015/join_queue?business_id=barbershop-1&queue_id=barbershop-1" \
  -H "Content-Type: application/json" \
  -d '{}'
```

**Response:**
```json
{
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "position": 1,
  "estimated_wait_minutes": 5
}
```

> **IMPORTANT:** Copy the `user_id` from the response. You will need it to leave the queue or check status.

### 2. Leave Queue
If a customer decides to walk away, they can leave the queue.

**Command:**
```bash
# Replace <PASTE_UUID_HERE> with the user_id from step 1
curl -X POST "http://localhost:2015/leave_queue?business_id=barbershop-1&queue_id=barbershop-1" \
  -H "Content-Type: application/json" \
  -d '{"user_id": "<PASTE_UUID_HERE>"}'
```

---

## B. The Staff Flow (Authenticated)

Staff members (e.g., barbers) need to login to manage the queue and call customers.

### 1. Login (Start Magic Auth)
Initiate the login flow using an email address.

**Command:**
```bash
curl -X POST http://localhost:2015/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "owner@barbershop.com"}'
```

> **INSTRUCTION:** Check your **terminal logs** (where `go run cmd/server/main.go` is running) to find the **6-digit Magic Code**.

### 2. Verify (Get Token)
Exchange the magic code for a JWT.

**Command:**
```bash
# Replace "123456" with the code from your logs
curl -X POST http://localhost:2015/auth/verify \
  -H "Content-Type: application/json" \
  -d '{"email": "owner@barbershop.com", "code": "123456"}'
```

**Response:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": { ... }
}
```

> **IMPORTANT:** Copy the `token` string. This is your Bearer token for authorized requests.

### 3. Call Next Customer (Service Routing)
When a staff member is ready, they call the next waiting customer to their counter.

**Command:**
```bash
# Replace <PASTE_TOKEN_HERE> with your JWT
curl -X POST http://localhost:2015/queues/barbershop-1/call-next \
  -H "Authorization: Bearer <PASTE_TOKEN_HERE>" \
  -H "Content-Type: application/json" \
  -d '{"counter_id": "Counter 3"}'
```

**Effect:**
1. The system finds the next `WAITING` ticket.
2. Changes status to `READY`.
3. Sets `assigned_to` to "Counter 3".
4. **Triggers NATS:** A message is published to `events.queue.barbershop-1` with `{ "instruction": "Go to Counter 3" }`.

---

## C. Media Verification

Verify the Caddy -> MinIO pipeline is working correctly.

### 1. Upload Test
First, ensure MinIO is configured and has the default placeholder image.

```bash
# Run from project root
./scripts/reset_media_server.sh
```

### 2. Network Test
Check that Caddy serves the image with correct headers.

**Command:**
```bash
curl -I http://localhost:2015/media/default/logo.png
```

**Expected Response:**
```
HTTP/1.1 200 OK
Content-Type: image/png
Server: MinIO
Cache-Control: public, max-age=86400
Date: ...
```

### 3. API Integration Test
Verify the Go Worker returns the correct Media URLs.

**Command:**
```bash
# Replace "biz1" and "q1" with your queue details
curl -X GET "http://localhost:2015/queue_status?business_id=biz1&queue_id=q1"
```

**Expected Response:**
```json
{
  ...
  "media": {
    "logo_url": "http://localhost:2015/media/biz1/logo.png",
    "header_url": "http://localhost:2015/media/biz1/header.jpg"
  }
}
```

---

## D. Managing Test Images

To update or add test images:

1.  **Add Files**: Place them in `virtual-queue-go/assets/defaults/`.
2.  **Sync**: Run `./scripts/reset_media_server.sh`.
3.  **Verify**: Access `http://localhost:2015/media/default/<filename>`.
