package middleware

import (
	"goProject/dockerGo/requsetTimeHandler"
	"log"
	"net"
	"net/http"
	"time"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Log the method and the requested URL
		log.Printf("Started %s %s", r.Method, r.URL.Path)

		// Call the next handler in the chain
		next.ServeHTTP(w, r)

		// Log how long it took
		log.Printf("Completed in %v", time.Since(start))
	})
}
func ErrorHandlingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				// Log the error and send a user-friendly message
				log.Printf("Error occurred: %v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
func RateLimitingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// IP'yi porttan ayır
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			// IP ayrıştırılamazsa direkt hata dön
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		limiter := requsetTimeHandler.GetClientLimiter(ip)

		if !limiter.Allow() {
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}
