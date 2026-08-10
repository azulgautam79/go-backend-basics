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
	userHandler := *user.NewHandler(userStore)

	mux := http.NewServeMux()

	mux.HandleFunc("/auth/register", userHandler.Register)
	mux.HandleFunc("/auth/login", userHandler.Login)
	mux.HandleFunc("/auth/me", userHandler.Me)
	mux.HandleFunc("/products", handler.Products)
	mux.HandleFunc("/products/", handler.Product)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	log.Println("Server listening on :8080")

	err := server.ListenAndServe()

	if err != nil {
		log.Fatal(err)
	}
}
