// Package main is the entry point for the Medicine Reminder API
package main

import (
	"log"
	"medicine-reminder/database"
	"medicine-reminder/handlers"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

// setupRouter configures and returns the API router with all route handlers
func setupRouter() *mux.Router {
	router := mux.NewRouter()

	// API Routes
	router.HandleFunc("/api/medicines", handlers.GetMedicines).Methods("GET")
	router.HandleFunc("/api/medicines", handlers.CreateMedicine).Methods("POST")
	router.HandleFunc("/api/medicines/{id}", handlers.GetMedicine).Methods("GET")
	router.HandleFunc("/api/medicines/{id}", handlers.UpdateMedicine).Methods("PUT")
	router.HandleFunc("/api/medicines/{id}", handlers.DeleteMedicine).Methods("DELETE")

	return router
}

// setupCORS configures and returns the CORS handler
func setupCORS(router *mux.Router) http.Handler {
	return cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
	}).Handler(router)
}

func main() {
	// Initialize database connection
	database.InitDB()
	defer database.DB.Close()

	// Setup router and CORS
	router := setupRouter()
	corsHandler := setupCORS(router)

	// Start server
	const port = ":8080"
	log.Printf("Server starting on port %s...", port)
	log.Fatal(http.ListenAndServe(port, corsHandler))
}
