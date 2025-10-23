package middleware

import (
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

// LogEntry represents a structured log entry
type LogEntry struct {
	Timestamp  string `json:"timestamp"`
	Method     string `json:"method"`
	Path       string `json:"path"`
	Query      string `json:"query,omitempty"`
	StatusCode int    `json:"status_code"`
	Duration   string `json:"duration"`
	UserAgent  string `json:"user_agent,omitempty"`
	IP         string `json:"ip"`
}

// Custom ResponseWriter
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rec *responseWriter) WriteHeader(code int) {
	rec.statusCode = code
	rec.ResponseWriter.WriteHeader(code)
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Create response writer that captures status code
		wrapped := &responseWriter{
			ResponseWriter: w,
			statusCode:     200, // Default status code
		}
		// Process request
		next.ServeHTTP(wrapped, r)
		duration := time.Since(start)

		// Log the request details
		logEntry := LogEntry{
			Timestamp:  time.Now().Format(time.RFC3339),
			Method:     r.Method,
			Path:       r.URL.Path,
			Query:      r.URL.RawQuery,
			StatusCode: wrapped.statusCode,
			Duration:   duration.String(),
			UserAgent:  r.UserAgent(),
			IP:         getIPAddress(r),
		}

		// Log with Logrus based on status code
		logger := log.WithFields(log.Fields{
			"method":      logEntry.Method,
			"path":        logEntry.Path,
			"query":       logEntry.Query,
			"status_code": logEntry.StatusCode,
			"duration":    logEntry.Duration,
			"user_agent":  logEntry.UserAgent,
			"ip":          logEntry.IP,
			"timestamp":   logEntry.Timestamp,
		})

		// Color-coded logging based on HTTP status
		switch {
		case logEntry.StatusCode >= 500:
			logger.Error("Server error")
		case logEntry.StatusCode >= 400:
			logger.Error("Client error")
		case logEntry.StatusCode >= 300:
			logger.Info("Redirection")
		default:
			logger.Info("Request completed")
		}
	})
}

func getIPAddress(r *http.Request) string {
	if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
		return ip
	}
	if ip := r.Header.Get("X-Real-IP"); ip != "" {
		return ip
	}
	return r.RemoteAddr

}
