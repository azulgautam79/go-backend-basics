package product

import (
	"encoding/json"
	"net/http"
	"strconv"
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

// Handles:
// GET  /products
// POST /products
func (h *Handler) Products(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.getProducts(w, r)

	case http.MethodPost:
		h.createProduct(w, r)

	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

// Handles:
// GET    /products/:id
// PUT    /products/:id
// DELETE /products/:id
func (h *Handler) Product(w http.ResponseWriter, r *http.Request) {
	idString := strings.TrimPrefix(r.URL.Path, "/products/")

	id, err := strconv.ParseInt(idString, 10, 64)

	if err != nil {
		http.Error(w, "invalid product ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.getProduct(w, id)

	case http.MethodPut:
		h.updateProduct(w, r, id)

	case http.MethodDelete:
		h.deleteProduct(w, r, id)

	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) getProducts(w http.ResponseWriter, r *http.Request) {

	minPriceStr := r.URL.Query().Get("min_price")

	var products []Product

	if minPriceStr == "" {
		products = h.store.GetAll()
	} else {
		minPrice, err := strconv.ParseFloat(minPriceStr, 64)

		if err != nil {
			http.Error(w, "invalid min_price", http.StatusBadRequest)
			return
		}
		products = h.store.GetByMinPrice(minPrice)
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(products)
}

func (h *Handler) getProduct(w http.ResponseWriter, id int64) {
	product, err := h.store.GetByID(id)

	if err != nil {
		http.Error(w, "product not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}

func (h *Handler) createProduct(w http.ResponseWriter, r *http.Request) {
	var product Product

	err := json.NewDecoder(r.Body).Decode(&product)

	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	product = h.store.Create(product)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(product)
}

func (h *Handler) updateProduct(
	w http.ResponseWriter,
	r *http.Request,
	id int64,
) {
	var product Product

	err := json.NewDecoder(r.Body).Decode(&product)

	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	product, err = h.store.Update(id, product)

	if err != nil {
		http.Error(w, "product not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}

func (h *Handler) deleteProduct(
	w http.ResponseWriter,
	_ *http.Request,
	id int64,
) {
	err := h.store.Delete(id)

	if err != nil {
		http.Error(w, "product not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
