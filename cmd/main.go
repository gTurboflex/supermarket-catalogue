package main

import (
	"log"
	"net/http"
	"time"

	"supermarket-catalogue/internal/handlers"
	"supermarket-catalogue/internal/middleware"
	"supermarket-catalogue/internal/repository"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

func main() {
	err := repository.Init()
	if err != nil {
		log.Fatal("Database initialization failed:", err)
	}
	r := mux.NewRouter()

	r.Use(middleware.CORSMiddleware)
	r.Use(middleware.LoggingMiddleware)

	r.HandleFunc("/register", handlers.RegisterHandler).Methods("POST")
	r.HandleFunc("/login", handlers.LoginHandler).Methods("POST")
	r.HandleFunc("/health", handlers.HealthCheck).Methods("GET")
	r.HandleFunc("/products/compare/{barcode}", handlers.CompareByBarcode).Methods("GET")
	r.HandleFunc("/products/{id}", handlers.GetProductByID).Methods("GET")
	r.HandleFunc("/products", handlers.GetProducts).Methods("GET")
	r.HandleFunc("/admin", handlers.AdminPage).Methods("GET")

	authRouter := r.PathPrefix("").Subrouter()
	authRouter.Use(middleware.AuthMiddleware)

	authRouter.HandleFunc("/basket/compare", handlers.CompareBasket).Methods("POST")
	authRouter.HandleFunc("/users", handlers.GetUsersHandler).Methods("GET")
	authRouter.HandleFunc("/me", handlers.GetCurrentUserHandler).Methods("GET")

	authRouter.HandleFunc("/supermarkets/stats", handlers.GetSupermarketStats).Methods("GET")

	authRouter.HandleFunc("/products", handlers.CreateProduct).Methods("POST")
	authRouter.HandleFunc("/products/{id}", handlers.UpdateProduct).Methods("PUT")
	authRouter.HandleFunc("/products/{id}", handlers.DeleteProduct).Methods("DELETE")

	r.HandleFunc("/supermarkets", handlers.GetSupermarkets).Methods("GET")
	r.HandleFunc("/supermarkets/{id}", handlers.GetSupermarketByID).Methods("GET")

	adminRouter := r.PathPrefix("").Subrouter()
	adminRouter.Use(middleware.AuthMiddleware, middleware.AdminMiddleware)
	adminRouter.HandleFunc("/admin/supermarkets", handlers.CreateSupermarket).Methods("POST")
	adminRouter.HandleFunc("/admin/supermarkets/{id}", handlers.UpdateSupermarket).Methods("PUT")
	adminRouter.HandleFunc("/admin/supermarkets/{id}", handlers.DeleteSupermarket).Methods("DELETE")

	r.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("swagger-ui"),
	))

	fs := http.FileServer(http.Dir("./ui/html"))
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./ui/html/index.html")
	})
	r.HandleFunc("/style.css", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./ui/html/style.css")
	})
	r.HandleFunc("/script.js", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./ui/html/script.js")
	})
	r.PathPrefix("/").Handler(fs)

	server := &http.Server{
		Addr:         ":8080",
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Println("Server starting on http://localhost:8080")
	log.Println("Default admin: admin@example.com / admin123")
	log.Println("Swagger docs available at http://localhost:8080/swagger/index.html")
	log.Println("Frontend available at http://localhost:8080")

	if err := server.ListenAndServe(); err != nil {
		log.Fatal("Server failed:", err)
	}
}
