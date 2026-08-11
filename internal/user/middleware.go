package user

import (
	"context"
	"net/http"
	"strings"
)

type contextKey string

const userContextKey contextKey = "user"

func AuthMiddleware(
	jwtService *JWTService,
	store *Store,
) func(http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			authHeader := r.Header.Get("Authorization")

			if authHeader == "" {
				http.Error(w, "missing authorization header", http.StatusUnauthorized)
				return
			}

			if !strings.HasPrefix(authHeader, "Bearer ") {
				http.Error(w, "invalid authorization header", http.StatusUnauthorized)
				return
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")

			claims, err := jwtService.ValidateAccessToken(tokenString)

			if err != nil {
				http.Error(w, "invalid or expired session", http.StatusUnauthorized)
				return
			}

			currentUser, err := store.GetByID(claims.UserID)

			if err != nil {
				http.Error(w, "user not found", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), userContextKey, currentUser)

			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}

func UserFromContext(ctx context.Context) (User, bool) {
	user, ok := ctx.Value(userContextKey).(User)

	return user, ok
}


