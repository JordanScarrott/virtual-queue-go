package analytics

import (
	"net/http"
	"time"
)

// AnalyticsMiddleware automatically tracks API requests
func AnalyticsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Wrap ResponseWriter to capture status code
		ww := &responseWriter{w, http.StatusOK}

		next.ServeHTTP(ww, r)

		// Auto-track the API hit using the context which might contain UserID from Auth middleware
		props := map[string]interface{}{
			"path":        r.URL.Path,
			"method":      r.Method,
			"status":      ww.status,
			"duration_ms": time.Since(start).Milliseconds(),
		}

		// Fire and forget tracking
		// Ensure StartIngest is running to pick this up
		_ = Track(r.Context(), "api.request", props)
	})
}

// responseWriter wrapper to capture status code
type responseWriter struct {
	http.ResponseWriter
	status int
}

func (w *responseWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}
