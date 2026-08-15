package user

import (
	"encoding/json"
	"errors"
	"net/http"
)

type Handler struct {
	authService *AuthService
}

type RegisterInput struct {
	Body RegisterRequest
}

type RegisterOutput struct {
	Body User
}

func NewHandler(
	authService *AuthService,
) *Handler {
	return &Handler{
		authService: authService,
	}
}

// ! Register Request
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest

	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if req.Name == "" || req.Email == "" || req.Password == "" {
		http.Error(w, "name, email or password are required", http.StatusBadRequest)
		return
	}

	user, err := h.authService.Register(
		r.Context(),
		req.Name,
		req.Email,
		req.Password,
	)

	if err != nil {
		if errors.Is(err, ErrEmailAlreadyUsed) {
			http.Error(
				w, "email already exists",
				http.StatusConflict,
			)
			return
		}
		http.Error(w, "failed to create user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(user)
}

// ! Login
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest

	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	result, err := h.authService.Login(
		r.Context(),
		req.Email,
		req.Password,
	)

	if err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    result.RefreshToken,
		Path:     "/auth",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		Expires:  result.ExpiresAt,
	})

	response := map[string]any{
		"accessToken": result.AccessToken,
		"user":        result.User,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

}

// ! Logout
func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")

	if err == nil {
		_ = h.authService.Logout(
			r.Context(),
			cookie.Value,
		)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/auth",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
	})

	w.WriteHeader(http.StatusNoContent)
}

// ! Me
func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	currentUser, ok := UserFromContext(r.Context())

	if !ok {
		http.Error(w, "user not authenticated", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(currentUser)
}

// ! Refresh Endpoint
func (h *Handler) Refresh(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")

	if err != nil {
		http.Error(w, "refresh token missing", http.StatusUnauthorized)
		return
	}

	result, err := h.authService.Refresh(
		r.Context(),
		cookie.Value,
	)

	if err != nil {
		http.Error(
			w,
			"invalid or expired refresh token",
			http.StatusUnauthorized,
		)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    result.RefreshToken,
		Path:     "/auth",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		Expires:  result.ExpiresAt,
	})

	response := map[string]any{
		"accessToken": result.AccessToken,
		"user":        result.User,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
