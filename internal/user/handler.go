package user

import (
	"encoding/json"
	"net/http"
	"time"
)

type Handler struct {
	store      *Store
	jwtService *JWTService
}

func NewHandler(store *Store, jwtService *JWTService) *Handler {
	return &Handler{
		store:      store,
		jwtService: jwtService,
	}
}

// ! Register Request
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest

	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
	}

	if req.Name == "" || req.Email == "" || req.Password == "" {
		http.Error(w, "name, email or password are required", http.StatusBadRequest)
		return
	}

	hashedPassword, err := HashPassword(req.Password)

	if err != nil {
		http.Error(w, "failed to hash password", http.StatusInternalServerError)
		return
	}

	user := User{
		Name:     req.Name,
		Email:    req.Email,
		Password: hashedPassword,
		Role:     RoleCustomer,
	}

	user, err = h.store.Create(user)

	if err != nil {

		if err == ErrEmailAlreadyUsed {
			http.Error(w, "email already exists", http.StatusConflict)
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
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	user, err := h.store.GetByEmail(req.Email)

	if err != nil {
		http.Error(w, "invalid email or password", http.StatusUnauthorized)
		return
	}

	if !CheckPassword(req.Password, user.Password) {
		http.Error(w, "invalid email or password", http.StatusUnauthorized)
		return
	}

	accessToken, err := h.jwtService.GenerateAccessToken(user)

	if err != nil {
		http.Error(w, "failed to create access token", http.StatusInternalServerError)
		return
	}

	refreshToken, err := h.store.CreateRefreshToken(user.ID, 7*24*time.Hour)

	if err != nil {
		http.Error(w, "failed to create refresh token", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken.Token,
		Path:     "/auth",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		Expires:  refreshToken.ExpiresAt,
	})

	response := map[string]any{
		"accessToken": accessToken,
		"user":        user,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

}

// ! Logout
func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")

	if err != nil {
		_ = h.store.RevokeRefreshToken(cookie.Value)
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

	oldToken := cookie.Value

	refreshToken, err := h.store.GetRefreshToken(oldToken)

	if err != nil {
		http.Error(w, "invalid or expired refresh token", http.StatusUnauthorized)
		return
	}

	user, err := h.store.GetByID(refreshToken.UserID)

	if err != nil {
		http.Error(w, "user not found", http.StatusUnauthorized)
		return
	}

	//* Revoke old refresh token
	err = h.store.RevokeRefreshToken(oldToken)

	if err != nil {
		http.Error(w, "failed to revoke refresh token", http.StatusInternalServerError)
		return
	}

	// * Create new refresh token
	newRefreshToken, err := h.store.CreateRefreshToken(user.ID, 7*24*time.Hour)

	if err != nil {
		http.Error(w, "failed to create refresh token", http.StatusInternalServerError)
		return
	}

	//* Create new access token
	newAccessToken, err := h.jwtService.GenerateAccessToken(user)

	if err != nil {
		http.Error(w, "failed to create access token", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    newRefreshToken.Token,
		Path:     "/auth",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		Expires:  newRefreshToken.ExpiresAt,
	})

	response := map[string]any{
		"accessToken": newAccessToken,
		"user":        user,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
