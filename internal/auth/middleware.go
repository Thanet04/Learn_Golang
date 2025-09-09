package auth

import (
	"net/http"
	"strings"
)

// JWTMiddleware ตรวจสอบ JWT token ใน Authorization header
func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// ดึง Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Invalid token format", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
		if tokenString == "" {
			http.Error(w, "Invalid token format", http.StatusUnauthorized)
			return
		}

		claims, err := ValidateToken(tokenString)
		if err != nil {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		if claims != nil && claims.Subject != "" {
			r = r.WithContext(ContextWithUserID(r.Context(), claims.Subject))
		}

		next.ServeHTTP(w, r)
	})
}
