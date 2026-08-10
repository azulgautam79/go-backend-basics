package product

import "errors"

var ErrProductNotFound = errors.New("product not found")

type Store struct {
	products map[int64]Product
	nextID   int64
}

// Function
// Returns a pointer to Store
func NewStore() *Store {
	return &Store{
		products: make(map[int64]Product),
		nextID:   1,
	}
}

//! Services
// Method
// (s *Store) = receiver is a pointer to Store
//* Create Product
func (s *Store) Create(product Product) Product {
	product.ID = s.nextID
	s.nextID++

	s.products[product.ID] = product
	return product
}

// Method
// Again, s is a pointer to Store
//* Get Product by id
func (s *Store) GetByID(id int64) (Product, error) {
	product, exists := s.products[id]

	if !exists {
		return Product{}, ErrProductNotFound
	}
	return product, nil
}

//* Get All Products
func (s *Store) GetAll() []Product {
	products := make([]Product, 0, len(s.products))

	for _, product := range s.products {
		products = append(products, product)
	}
	return products
}

//* Get Products greater than 500
func (s *Store) GetByMinPrice(minPrice float64) []Product {
	products := make([]Product, 0)

	for _, product := range s.products {
		if product.Price >= minPrice {
			products = append(products, product)
		}
	}
	return products
}

//* Update a product
func (s *Store) Update(id int64, product Product) (Product, error) {
	_, exists := s.products[id]

	if !exists {
		return Product{}, ErrProductNotFound
	}

	product.ID = id
	s.products[id] = product

	return product, nil
}

//* Delete a product
func (s *Store) Delete(id int64) error {
	_, exists := s.products[id]

	if !exists {
		return ErrProductNotFound
	}

	delete(s.products, id)
	return nil
}
