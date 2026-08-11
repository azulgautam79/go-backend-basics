package main

import (
	"log"
	"net/http"

	"github.com/azulgautam79/ecommerce-basic/internal/product"
	"github.com/azulgautam79/ecommerce-basic/internal/user"
)

func main() {
	store := product.NewStore()
	handler := product.NewHandler(store)

	userStore := user.NewStore()
	jwtService := user.NewJWTService(
		"super-secret-development-key",
	)
	userHandler := user.NewHandler(userStore, jwtService)

	authMiddleware := user.AuthMiddleware(jwtService, userStore)

	mux := http.NewServeMux()

	mux.HandleFunc("/auth/register", userHandler.Register)
	mux.HandleFunc("/auth/login", userHandler.Login)
	mux.HandleFunc("/auth/refresh", userHandler.Refresh)
	mux.HandleFunc("/auth/logout", userHandler.Logout)
	mux.Handle("/auth/me", authMiddleware(http.HandlerFunc(userHandler.Me)))

	// mux.HandleFunc("/products", handler.Products)
	mux.Handle(
		"/products",
		authMiddleware(
			user.RequireRole(
				user.RoleSeller,
				user.RoleAdmin,
			)(
				http.HandlerFunc(handler.Products),
			),
		),
	)
	mux.HandleFunc("/products/", handler.Product)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	log.Println("Server listening on http://localhost:8080")

	err := server.ListenAndServe()

	if err != nil {
		log.Fatal(err)
	}
}
