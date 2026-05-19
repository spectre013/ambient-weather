package main

import (
	"encoding/json"
	"net/http"
	"time"
)

func writeJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")

	b, err := json.Marshal(payload)
	if err != nil {
		logger.WithError(err).Error("json marshal failed")
		// We may already have been about to write a 200 -- only call
		// WriteHeader once.
		http.Error(w, `{"error":"internal error"}`, http.StatusInternalServerError)
		return
	}

	if status != http.StatusOK {
		w.WriteHeader(status)
	}

	if _, err := w.Write(b); err != nil {
		// We can't change the status code at this point (headers are
		// already flushed). Just log.
		logger.WithError(err).Error("response write failed")
	}
}

// writeJSONError writes a JSON error body with the given status.
func writeJSONError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

// loggingMiddleware logs each request.
func loggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		logger.WithFields(map[string]interface{}{
			"method":   r.Method,
			"uri":      r.RequestURI,
			"duration": time.Since(start).String(),
		}).Info("request")
	}
}
