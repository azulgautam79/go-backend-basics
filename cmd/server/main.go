package main

import (
	"log"
	"net/http"

	"github.com/azulgautam79/ecommerce-basic/internal/database"
	"github.com/azulgautam79/ecommerce-basic/internal/product"
	"github.com/azulgautam79/ecommerce-basic/internal/user"
)

func main() {

	//! Database
	db, dberr := database.NewPostgres(
		database.Config{
			Host:     "localhost",
			Port:     "5432",
			User:     "ecommerce",
			Password: "ecommerce",
			Name:     "ecommerce",
		},
	)

	if dberr != nil {
		log.Fatal(dberr)
	}

	defer db.Close()

	//* Product Dependencies
	store := product.NewStore()
	handler := product.NewHandler(store)

	//* User Dependencies
	userRepo := user.NewPostgresUserRepository(db)
	refreshTokenRepo := user.NewPostgresRefreshTokenRepository(db)

	jwtService := user.NewJWTService(
		"super-secret-development-key",
	)
	authService := user.NewAuthService(
		userRepo,
		refreshTokenRepo,
		jwtService,
	)

	userHandler := user.NewHandler(authService)

	authMiddleware := user.AuthMiddleware(jwtService, authService)

	//? Routes
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
