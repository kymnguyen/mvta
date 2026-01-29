package middleware

import (
	"context"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/kymnguyen/mvta/apps/backend/vehicle-svc/internal/infrastructure/resilience"
	"github.com/kymnguyen/mvta/apps/backend/vehicle-svc/internal/infrastructure/security"
)

var (
	apiCircuitBreaker = resilience.NewCircuitBreaker(5, 2, 10*time.Second)
	apiRetryPolicy    = resilience.NewRetryPolicy(2, 50*time.Millisecond, 500*time.Millisecond)
)

type ctxKey string

const ClaimsContextKey ctxKey = "claims"

type responseWriter struct {
	http.ResponseWriter
	statusCode int
	written    bool
}

func (rw *responseWriter) WriteHeader(code int) {
	if !rw.written {
		rw.statusCode = code
		rw.written = true
		rw.ResponseWriter.WriteHeader(code)
	}
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	if !rw.written {
		rw.statusCode = http.StatusOK
		rw.written = true
	}
	return rw.ResponseWriter.Write(b)
}

func ResilienceMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		var handlerErr error
		err := apiCircuitBreaker.Execute(func() error {
			return apiRetryPolicy.Execute(ctx, func() error {
				next.ServeHTTP(w, r)
				return nil
			})
		})
		if err != nil {
			handlerErr = err
			http.Error(w, "Service temporarily unavailable (resilience)", http.StatusServiceUnavailable)
			return
		}
		_ = handlerErr
	})
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := GetClientIP(r)
		startTime := time.Now()

		wrappedWriter := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		log.Printf("[REQUEST] Method=%s Path=%s IP=%s UserAgent=%s",
			r.Method, r.URL.Path, ip, r.Header.Get("User-Agent"))

		next.ServeHTTP(wrappedWriter, r)

		duration := time.Since(startTime)
		log.Printf("[RESPONSE] Method=%s Path=%s IP=%s StatusCode=%d Duration=%dms",
			r.Method, r.URL.Path, ip, wrappedWriter.statusCode, duration.Milliseconds())
	})
}

func AuthMiddleware(requiredRole string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, `{"message":"missing authorization header"}`, http.StatusUnauthorized)
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, `{"message":"invalid authorization header format"}`, http.StatusUnauthorized)
				return
			}

			token := parts[1]

			claims, err := security.ValidateToken(token)
			if err != nil {
				log.Printf("Token validation failed: %v", err)
				http.Error(w, `{"message":"invalid or expired token"}`, http.StatusUnauthorized)
				return
			}

			if err := security.VerifyRole(claims, requiredRole); err != nil {
				log.Printf("Role verification failed: %v", err)
				http.Error(w, `{"message":"insufficient permissions"}`, http.StatusForbidden)
				return
			}

			ctx := context.WithValue(r.Context(), ClaimsContextKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetClaimsFromContext(r *http.Request) *security.Claims {
	claims, ok := r.Context().Value(ClaimsContextKey).(*security.Claims)
	if !ok {
		return nil
	}
	return claims
}

func GetClientIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		return xff
	}

	if xrip := r.Header.Get("X-Real-IP"); xrip != "" {
		return xrip
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}
