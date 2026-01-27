package middleware

import (
	"log"
	"net"
	"net/http"
)

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
