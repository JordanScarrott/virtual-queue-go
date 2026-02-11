# Adding Analytics

[← Back to README](../../README.md)

This guide explains how to add analytics tracking to the codebase. We have streamlined the process to a single line of code for most use cases.

## The "One-Liner"
To track an event anywhere in the code (API handlers, Activities, etc.), simply call:

```go
analytics.Track(ctx, "event.name", map[string]interface{}{
    "property": "value",
})
```

**That's it.** The system automatically:
1. Extracts `UserID` from the `context.Context` (if authenticated).
2. Extracts `BusinessID` from the `context.Context` (if set).
3. Adds a timestamp.
4. Asynchronously sends the event to NATS.
5. Persists it to Postgres.

---

## Automatic API Tracking
We have an `AnalyticsMiddleware` that automatically tracks every API request:

```go
// In main.go or router setup
http.Handle("/", analytics.AnalyticsMiddleware(myHandler))
```

This produces an event `api.request` with properties:
- `path`: URL path
- `method`: HTTP method
- `status`: HTTP status code
- `duration_ms`: Request duration

## Advanced Usage

### Setting Context IDs
If you are in a flow where IDs are not automatically in the context (e.g., a background job), you can manually properties:

```go
// In a manual flow
props := map[string]interface{}{
    "item_id": 123,
}
// Pass context if available, otherwise UserID will be empty
analytics.Track(context.Background(), "job.completed", props)
```

To ensure `UserID` is tracked, ensure your context is derived from the Auth middleware or manually set it (though internal context setting is typically handled by middleware).

### Adding New Event Types
There is no schema registry strictly enforced yet. You can simply use a new event name string (e.g., `user.signup`, `queue.joined`).
However, for consistency, try to use `entity.action` format.

## Setup (Already Done)
The system is initialized in `main.go`:
1. `analytics.NewTracker(nc)` creates the tracker.
2. `analytics.SetGlobalTracker(tracker)` sets the singleton.
3. `analytics.StartIngest(nc, repo)` starts the consumer.
