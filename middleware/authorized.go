package middleware_auth

import (
	"context"
	"ecommerce/helper"
	"ecommerce/service"
	"net/http"

	"go.uber.org/zap"
)

type AuthMiddleware struct {
	AuthService service.AuthService
	Log         *zap.Logger
}

func NewAuthMiddleware(authService service.AuthService, logger *zap.Logger) *AuthMiddleware {
	return &AuthMiddleware{
		AuthService: authService,
		Log:         logger,
	}
}

func (m *AuthMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m.Log.Debug("Middleware: Processing request", zap.String("path", r.URL.Path))

		token := r.Header.Get("Authorization")
		if len(token) > 7 && token[:7] == "Bearer " {
			token = token[7:]
		}

		if token == "" {
			m.Log.Warn("Middleware: Token missing in request")
			helper.SendJSONResponse(w, http.StatusUnauthorized, "Token required", nil)
			return
		}

		userID, err := m.AuthService.VerifyToken(token)
		if err != nil {
			m.Log.Warn("Middleware: Invalid or expired token", zap.Error(err))
			helper.SendJSONResponse(w, http.StatusUnauthorized, "Invalid or expired token", nil)
			return
		}

		m.Log.Info("Middleware: Token validated successfully", zap.Int("userID", userID))

		ctx := context.WithValue(r.Context(), "userID", userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
