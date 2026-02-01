package middleware

import (
	"fmt"
	"net/http"
	"time"
)

// LoggingMiddleware перехватывает запросы и логирует их
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Использование Goroutine
		// Логируем детали запроса асинхронно
		go func(method, path string, startTime time.Time) {
			duration := time.Since(startTime)
			fmt.Printf("--- LOG RECORD ---\n")
			fmt.Printf("Time: %s\nMethod: %s\nPath: %s\nDuration: %v\n------------------\n",
				startTime.Format(time.RFC850), method, path, duration)
		}(r.Method, r.URL.Path, start)

		next.ServeHTTP(w, r)
	})
}
