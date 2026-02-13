package middleware

import (
	"net/http"
	"strconv"
	"strings"
	"supermarket-catalogue/internal/auth"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		publicRoutes := map[string]bool{
			"/":          true,
			"/register":  true,
			"/login":     true,
			"/health":    true,
			"/style.css": true,
			"/script.js": true,
		}

		if publicRoutes[path] || strings.HasPrefix(path, "/swagger/") {
			next.ServeHTTP(w, r)
			return
		}

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"error": "Authorization header required"}`))
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"error": "Invalid authorization format"}`))
			return
		}

		claims, err := auth.VerifyToken(parts[1])
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"error": "Invalid or expired token"}`))
			return
		}
		r.Header.Set("X-User-ID", strconv.Itoa(claims.UserID))
		r.Header.Set("X-User-Role", claims.Role)

		next.ServeHTTP(w, r)
	})
}

func AdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		role := r.Header.Get("X-User-Role")
		if role != "admin" {
			http.Error(w, `{"error": "Admin access required"}`, http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}
