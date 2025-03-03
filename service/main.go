package main

import (
	"database/sql"
	"embed"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

	"github.com/yourusername/template-service/db"
	"github.com/yourusername/template-service/handlers"
)

//go:embed templates/*
var templateFS embed.FS

//go:embed static/*
var staticFS embed.FS

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: Error loading .env file, using defaults or environment variables")
	}

	// Wait for database to be ready
	db.WaitForDatabase()

	// Set up database connection
	db.SetupDatabase()
	defer func(DB *sql.DB) {
		err := DB.Close()
		if err != nil {
			print(err)
		}
	}(db.DB)

	// Set up template filesystem for handlers
	handlers.FS = templateFS

	// Set up HTTP router
	router := mux.NewRouter()

	// Serve static files
	router.PathPrefix("/static/").Handler(http.FileServer(http.FS(staticFS)))

	// Create API subrouter
	apiRouter := router.PathPrefix("/api").Subrouter()

	// Set up routes
	router.HandleFunc("/", handlers.HandleIndex)
	router.HandleFunc("/templates", handlers.HandleListTemplates)
	router.HandleFunc("/templates/new", handlers.HandleNewTemplateForm).Methods("GET")
	router.HandleFunc("/templates", handlers.HandleCreateTemplate).Methods("POST")
	router.HandleFunc("/templates/{id}", handlers.HandleViewTemplate).Methods("GET")
	router.HandleFunc("/templates/{id}/render", handlers.HandleRenderTemplate).Methods("POST")
	router.HandleFunc("/templates/{id}/pdf", handlers.HandleGeneratePDF).Methods("POST")

	// Web UI Health check endpoint
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")

	// API Routes
	apiRouter.HandleFunc("/health", handlers.APIHealthCheck).Methods("GET")

	// Templates API
	apiRouter.HandleFunc("/templates", handlers.APIGetTemplates).Methods("GET")
	apiRouter.HandleFunc("/templates", handlers.APICreateTemplate).Methods("POST")
	apiRouter.HandleFunc("/templates/{id}", handlers.APIGetTemplate).Methods("GET")
	apiRouter.HandleFunc("/templates/{id}", handlers.APIUpdateTemplate).Methods("PUT")
	apiRouter.HandleFunc("/templates/{id}", handlers.APIDeleteTemplate).Methods("DELETE")
	apiRouter.HandleFunc("/templates/{id}/render", handlers.APIRenderTemplate).Methods("POST")

	// Template Variables API
	apiRouter.HandleFunc("/templates/{id}/variables", handlers.APIGetTemplateVariables).Methods("GET")
	apiRouter.HandleFunc("/templates/{id}/variables", handlers.APIAddTemplateVariable).Methods("POST")

	// Categories API
	apiRouter.HandleFunc("/categories", handlers.APIGetCategories).Methods("GET")

	// Start the server
	port := getEnv("SERVER_PORT", "8080")
	log.Printf("Starting template service on port %s...", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}

// Helper to get environment variable with default
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
