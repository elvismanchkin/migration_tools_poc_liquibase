package main

import (
	"database/sql"
	"embed"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

	"github.com/elvismanchkin/migration_tools_poc_liquibase/db"
	"github.com/elvismanchkin/migration_tools_poc_liquibase/handlers"
)

//go:embed templates/*
var templateFS embed.FS

//go:embed static/*
var staticFS embed.FS

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: Error loading .env file, using defaults or environment variables")
	}

	db.WaitForDatabase()

	db.SetupDatabase()
	defer func(DB *sql.DB) {
		err := DB.Close()
		if err != nil {
			print(err)
		}
	}(db.DB)

	handlers.FS = templateFS
	router := mux.NewRouter()

	router.PathPrefix("/static/").Handler(http.FileServer(http.FS(staticFS)))
	apiRouter := router.PathPrefix("/api").Subrouter()

	router.HandleFunc("/", handlers.HandleIndex)
	router.HandleFunc("/templates", handlers.HandleListTemplates)
	router.HandleFunc("/templates/new", handlers.HandleNewTemplateForm).Methods("GET")
	router.HandleFunc("/templates", handlers.HandleCreateTemplate).Methods("POST")
	router.HandleFunc("/templates/{id}", handlers.HandleViewTemplate).Methods("GET")
	router.HandleFunc("/templates/{id}/render", handlers.HandleRenderTemplate).Methods("POST")
	router.HandleFunc("/templates/{id}/pdf", handlers.HandleGeneratePDF).Methods("POST")

	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte("OK")); err != nil {
			log.Printf("Error writing the response: %v", err)
		}
	}).Methods("GET")

	apiRouter.HandleFunc("/health", handlers.APIHealthCheck).Methods("GET")

	apiRouter.HandleFunc("/templates", handlers.APIGetTemplates).Methods("GET")
	apiRouter.HandleFunc("/templates", handlers.APICreateTemplate).Methods("POST")
	apiRouter.HandleFunc("/templates/{id}", handlers.APIGetTemplate).Methods("GET")
	apiRouter.HandleFunc("/templates/{id}", handlers.APIUpdateTemplate).Methods("PUT")
	apiRouter.HandleFunc("/templates/{id}", handlers.APIDeleteTemplate).Methods("DELETE")
	apiRouter.HandleFunc("/templates/{id}/render", handlers.APIRenderTemplate).Methods("POST")
	apiRouter.HandleFunc("/templates/{id}/variables", handlers.APIGetTemplateVariables).Methods("GET")
	apiRouter.HandleFunc("/templates/{id}/variables", handlers.APIAddTemplateVariable).Methods("POST")
	apiRouter.HandleFunc("/categories", handlers.APIGetCategories).Methods("GET")

	port := getEnv("SERVER_PORT", "8080")
	log.Printf("Starting template service on port %s...", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
