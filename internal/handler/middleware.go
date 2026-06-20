package handler

import (
	"context"
	"net/http"
	"strings"

	"github.com/google/uuid"

	"github.com/gitsang/order/internal/service"
	"github.com/gitsang/order/pkg/response"
)

func AuthMiddleware(authService *service.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				response.Error(w, http.StatusUnauthorized, "Authorization header required")
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || parts[0] != "Bearer" {
				response.Error(w, http.StatusUnauthorized, "Invalid authorization header format")
				return
			}

			userID, username, role, err := authService.ValidateToken(parts[1])
			if err != nil {
				response.Error(w, http.StatusUnauthorized, "Invalid token")
				return
			}

			ctx := context.WithValue(r.Context(), "user_id", userID)
			ctx = context.WithValue(ctx, "username", username)
			ctx = context.WithValue(ctx, "role", role)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func AdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		role, ok := r.Context().Value("role").(string)
		if !ok || role != "admin" {
			response.Error(w, http.StatusForbidden, "Admin access required")
			return
		}

		next.ServeHTTP(w, r)
	})
}

func GetUserID(r *http.Request) uuid.UUID {
	userID, _ := r.Context().Value("user_id").(uuid.UUID)
	return userID
}

func GetUsername(r *http.Request) string {
	username, _ := r.Context().Value("username").(string)
	return username
}

func GetRole(r *http.Request) string {
	role, _ := r.Context().Value("role").(string)
	return role
}
