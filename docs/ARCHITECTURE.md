# Project Architecture

## 📸 The Media Pipeline

The application employs a "Short Circuit" architecture for handling media assets (images, logos, headers). This design decouples high-bandwidth media serving from the application logic.

### Request Flow

1.  **Media Requests** (`GET /media/*`):
    *   **Caddy** receives the request.
    *   **Rewrite**: Path `/media/foo` is rewritten to `/public/foo`.
    *   **Proxy**: Request is forwarded directly to **MinIO** (Port 9000).
    *   **Response**: Served immediately with `Cache-Control: public, max-age=86400`.
    *   *Note: The Go Worker is never touched.*

2.  **API Requests** (`GET /api/*` or other endpoints):
    *   **Caddy** proxies the request to the **Go Worker** (Port 8080).
    *   The Worker processes the request (Auth, DB, Temporal).

### Why "Sovereign"?

This architecture provides several key benefits:

*   **Zero Application CPU Usage**: The Go application doesn't waste cycles reading files from disk or piping bytes to the network. MinIO handles this efficiently.
*   **Aggressive Caching**: Caddy handles cache headers at the edge, ensuring browsers only request images once per day.
*   **No External CDN Required**: This setup mimics a CDN but runs entirely on `localhost` or a single VPS. It is friendly to self-hosting and air-gapped environments.
