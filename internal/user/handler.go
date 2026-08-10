package user

import (
	"encoding/json"
	"net/http"
	"strings"
)

type Handler struct {
	store *Store
}

func NewHandler(store *Store) *Handler {
	return &Handler{
		store: store,
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

	token, err := h.store.CreateSession(user.ID)

	if err != nil {
		http.Error(w, "failed to create session", http.StatusInternalServerError)
		return
	}

	response := map[string]any{
		"token": token,
		"user":  user,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

}

// ! Me
func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")

	if authHeader == "" {
		http.Error(w, "missing authorization header", http.StatusUnauthorized)
		return
	}

	if !strings.HasPrefix(authHeader, "Bearer ") {
		http.Error(w, "invalid authorization header", http.StatusUnauthorized)
		return
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")

	if token == "" {
		http.Error(w, "invalid authorization token", http.StatusUnauthorized)
		return
	}

	user, err := h.store.GetUserBySession(token)

	if err != nil {
		http.Error(w, "invalid or expired session", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(user)
}
