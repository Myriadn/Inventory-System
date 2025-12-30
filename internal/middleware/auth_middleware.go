package middleware

import (
	"context"
	"project-app-inventory-restapi-golang-anas/internal/repository"

	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

type AuthMiddleware struct {
	userRepo *repository.UserRepository
}

func NewAuthMiddleware(userRepo *repository.UserRepository) *AuthMiddleware {
	return &AuthMiddleware{userRepo: userRepo}
}

// Middleware untuk mengecek Token (Authentication)
func (m *AuthMiddleware) VerifyToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Ambil header Authorization: Bearer <token>
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
			return
		}

		// Split "Bearer" dan tokennya
		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			http.Error(w, "Invalid token format", http.StatusUnauthorized)
			return
		}

		tokenString := headerParts[1]

		// Validasi apakah format token adalah UUID yang valid
		if err := uuid.Validate(tokenString); err != nil {
			http.Error(w, "Invalid token format", http.StatusUnauthorized)
			return
		}

		// Cek ke Database
		session, err := m.userRepo.GetSessionByToken(r.Context(), tokenString)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		if session == nil {
			http.Error(w, "Invalid or revoked token", http.StatusUnauthorized)
			return
		}

		// Cek Expired
		if time.Now().After(session.ExpiredAt) {
			http.Error(w, "Token expired", http.StatusUnauthorized)
			return
		}

		// Simpan UserID dan Role ke Context supaya bisa dibaca di Handler nanti
		ctx := context.WithValue(r.Context(), UserIDKey, session.UserID)
		ctx = context.WithValue(ctx, RoleKey, session.Role)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Middleware untuk mengecek Role (Authorization)
func (m *AuthMiddleware) RequireRoles(allowedRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Ambil role dari context (yang diset oleh VerifyToken)
			userRole, ok := r.Context().Value(RoleKey).(string)
			if !ok {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Cek apakah role user ada di daftar allowedRoles
			isAllowed := false
			for _, role := range allowedRoles {
				if role == userRole {
					isAllowed = true
					break
				}
			}

			if !isAllowed {
				http.Error(w, "Forbidden: You don't have access to this resource", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
