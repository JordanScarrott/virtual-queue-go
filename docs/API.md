# Virtual Queue API Documentation

[← Back to README](../../README.md)

This document describes the HTTP API endpoints provided by the `red-duck` virtual queue service.

The API server runs on port `8080` by default, but is exposed via Caddy on port `2015`.

## Base URL

```
http://localhost:2015
```

## Endpoints

### 1. Create Queue

Starts a new queue workflow for a specific business.

- **URL**: `/create_queue`
- **Method**: `POST`
- **Query Parameters**:
    - `business_id` (string, required): The unique identifier of the business.
    - `queue_id` (string, required): The unique identifier of the queue.

#### Response (201 Created)

Returns the Temporal Workflow ID and Run ID.

```json
{
    "workflow_id": "biz1:q1",
    "run_id": "c92d5c36-7e9b-4e89-9e8c-57271457145c"
}
```

#### Response (500 Internal Server Error)

If the workflow fails to start.

---

### 2. Join Queue

Adds a user to an existing queue. This is a synchronous operation that waits for the workflow to process the update.

- **URL**: `/join_queue`
- **Method**: `POST`
- **Query Parameters**:
    - `business_id` (string, required): The unique identifier of the business.
    - `queue_id` (string, required): The unique identifier of the queue.
- **Request Body** (JSON):

```json
{
    "userId": "user-123"
}
```

#### Response (200 OK)

Returns the user's position in the queue (1-based index).

```json
{
    "position": 5
}
```

#### Response (409 Conflict)

If the user is already in the queue or if the queue is closed/not accepting joins.

#### Response (500 Internal Server Error)

If the update fails or the result cannot be retrieved.

---

### 3. Leave Queue

Removes a user from the queue.

- **URL**: `/leave_queue`
- **Method**: `POST`
- **Query Parameters**:
    - `business_id` (string, required): The unique identifier of the business.
    - `queue_id` (string, required): The unique identifier of the queue.
- **Request Body** (JSON):

```json
{
    "userId": "user-123"
}
```

#### Response (200 OK)

Returns the number of users remaining in the queue.

```json
{
    "remaining_users": 4
}
```

#### Response (500 Internal Server Error)

If the update fails.

---

### 4. Get Queue Status

Retrieves the current state of the queue, including the list of users.

- **URL**: `/queue_status`
- **Method**: `GET`
- **Query Parameters**:
    - `business_id` (string, required): The unique identifier of the business.
    - `queue_id` (string, required): The unique identifier of the queue.

#### Response (200 OK)

Returns the full queue object.

```json
{
    "ID": "q1",
    "BusinessID": "biz1",
    "Users": [
        "user-101",
        "user-102",
        "user-123"
    ]
}
```

#### Response (500 Internal Server Error)

If the query fails (e.g., if the workflow is not running).
