package product

//! Create Product Model
// Seperate model (domain model) independently of HTTP
// No HTTP, JSON Requests, databases, handlers, routers, //* For Seperation of Concerns
type Product struct {
	ID    int64   `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
	Stock int     `json:"stock"`
}
