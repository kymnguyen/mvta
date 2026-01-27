package middleware

import (
	"log"
	"net"
	"net/http"
	"time"

	"github.com/kymnguyen/mvta/apps/backend/vehicle-svc/internal/infrastructure/resilience"
)

var (
	apiCircuitBreaker = resilience.NewCircuitBreaker(5, 2, 10*time.Second)
	apiRetryPolicy    = resilience.NewRetryPolicy(2, 50*time.Millisecond, 500*time.Millisecond)
)

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
		log.Printf("Method=%s Path=%s IP=%s", r.Method, r.URL.Path, ip)
		next.ServeHTTP(w, r)
	})
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
