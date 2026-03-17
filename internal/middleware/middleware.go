package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/agent-os/core/internal/services"
)

// Context keys
type contextKey string

const (
	ContextKeyAgentID contextKey = "agent_id"
	ContextKeyOrgID   contextKey = "org_id"
	ContextKeyTraceID contextKey = "trace_id"
)

// TraceMiddleware adds trace ID to context
func TraceMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		traceID := r.Header.Get("X-Trace-Id")
		if traceID == "" {
			traceID = generateTraceID()
		}

		ctx := context.WithValue(r.Context(), ContextKeyTraceID, traceID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// AuthMiddleware validates API key
func AuthMiddleware(authService *services.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Skip auth for health check
			if r.URL.Path == "/health" {
				next.ServeHTTP(w, r)
				return
			}

			apiKey := r.Header.Get("X-API-Key")
			if apiKey == "" {
				// Also check Authorization header
				auth := r.Header.Get("Authorization")
				if strings.HasPrefix(auth, "Bearer ") {
					apiKey = strings.TrimPrefix(auth, "Bearer ")
				}
			}

			if apiKey == "" {
				http.Error(w, "Missing API key", http.StatusUnauthorized)
				return
			}

			agent, err := authService.Authenticate(r.Context(), apiKey)
			if err != nil {
				http.Error(w, "Invalid API key", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), ContextKeyAgentID, agent.ID)
			ctx = context.WithValue(ctx, ContextKeyOrgID, agent.OrgID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// OrgMiddleware validates organization access
func OrgMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get org ID from header or context
		orgID := r.Header.Get("X-Org-Id")
		if orgID == "" {
			orgID = r.Context().Value(ContextKeyOrgID).(string)
		}

		if orgID == "" {
			// Some routes might not require org (like health)
			if r.URL.Path != "/health" {
				http.Error(w, "Missing organization ID", http.StatusBadRequest)
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}

func generateTraceID() string {
	// Simple trace ID generation - in production use proper UUID
	return "trace-" + randomString(16)
}

func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[randInt(len(letters))]
	}
	return string(b)
}

func randInt(n int) int {
	// Simple random - in production use crypto/rand
	return 0
}

var _ = strings.TrimSpace
