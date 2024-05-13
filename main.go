package main

import (
	"log"
	"log/slog"
	"lucy/db"
	"lucy/handlers"
	"lucy/middlewares"
	"lucy/providers"
	"lucy/repositories"
	"lucy/services"
	"lucy/types"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	port := os.Getenv("PORT")

	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
	}))

	r.Use(
		middleware.Recoverer,
		middleware.Logger,
	)

	postgresPool := db.PostgresDB()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
	}))
	jwtProvider := providers.NewJWTProvider()

	userRepo := repositories.NewUserRepo(postgresPool)
	categoryRepo := repositories.NewCategoryRepository(postgresPool)
	productRepo := repositories.NewProductRepo(postgresPool)

	authMiddleware := middlewares.NewAuthMiddleware(userRepo, jwtProvider, logger)

	userService := services.NewUserService(userRepo, logger)
	authService := services.NewAuthService(userRepo, logger, jwtProvider)
	categoryService := services.NewCategoryService(categoryRepo, logger)
	productService := services.NewProductService(productRepo, categoryRepo, logger)

	userHandler := handlers.NewUserHandler(userService, logger)
	authHandler := handlers.NewAuthHandler(authService, logger)
	categoryHandler := handlers.NewCategoryHandler(categoryService, logger)
	productHandler := handlers.NewProductHandler(productService, logger)

	r.Route("/users", func(r chi.Router) {
		r.With(authMiddleware.AllowAccounts(types.AdminAccount)).Post("/", userHandler.HandleCreateAdminAccount)
		r.With(authMiddleware.AllowAccounts(types.AnyAccount)).Put("/password", userHandler.HandleChangeUserPassword)
	})

	r.Route("/auth", func(r chi.Router) {
		r.Post("/login", authHandler.HandleLogin)
	})

	r.Route("/categories", func(r chi.Router) {
		r.With(authMiddleware.AllowAccounts(types.AdminAccount)).Post("/", categoryHandler.HandleCreateCategory)
		r.With(authMiddleware.AllowAccounts(types.AnyAccount)).Get("/", categoryHandler.HandleGetAllCategories)
		r.With(authMiddleware.AllowAccounts(types.AdminAccount)).Patch("/enabled", categoryHandler.HandleToggleEnabled)
	})

	r.Route("/products", func(r chi.Router) {
		r.With(authMiddleware.AllowAccounts(types.SellerAccount)).Route("/", func(r chi.Router) {
			r.Post("/", productHandler.HandleCreateProduct)
		})
	})

	server := http.Server{
		Addr:         net.JoinHostPort("0.0.0.0", port),
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Println("Starting server on port", port)
	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
