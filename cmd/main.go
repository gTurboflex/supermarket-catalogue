package main

import (
	"log"
	"net/http"
	"time"

	"supermarket-catalogue/internal/handlers"
	"supermarket-catalogue/internal/repository"

	_ "supermarket-catalogue/docs" // Swagger docs

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title           Supermarket Goods Catalogue API
// @version         1.0
// @description     API for supermarket products catalogue
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.email  support@supermarket.com

// @host      localhost:8080
// @BasePath  /
func main() {
	err := repository.Init()
	if err != nil {
		log.Fatal("Database initialization failed:", err)
	}

	r := mux.NewRouter()

	r.HandleFunc("/products", handlers.GetProducts).Methods("GET")
	r.HandleFunc("/products", handlers.CreateProduct).Methods("POST")
	r.HandleFunc("/products/{id}", handlers.GetProductByID).Methods("GET")
	r.HandleFunc("/health", handlers.HealthCheck).Methods("GET")

	// Swagger с правильным URL
	r.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("swagger-ui"),
	))

	go backgroundLogger()

	server := &http.Server{
		Addr:         ":8080",
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Println("🚀 Server starting on http://localhost:8080")
	log.Println("📚 Swagger docs available at http://localhost:8080/swagger/index.html")

	if err := server.ListenAndServe(); err != nil {
		log.Fatal("Server failed:", err)
	}
}

func backgroundLogger() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		log.Println("✅ Background goroutine: API is running normally")
	}
}
