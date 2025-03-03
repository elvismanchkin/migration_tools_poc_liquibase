package db

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/elvismanchkin/migration_tools_poc_liquibase/ent"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

var DB *sql.DB
var EntClient *ent.Client
var Ctx = context.Background()

func WaitForDatabase() {
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "template_user")
	password := getEnv("DB_PASSWORD", "template_pass")
	dbname := getEnv("DB_NAME", "template_db")

	log.Println("Waiting for database to be ready...")

	maxRetries := 30
	for i := 0; i < maxRetries; i++ {
		connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			host, port, user, password, dbname)

		db, err := sql.Open("postgres", connStr)
		if err == nil {
			err = db.Ping()
			if err == nil {
				log.Println("Database is ready!")
				err := db.Close()
				if err != nil {
					log.Printf("Error closing database connection: %v", err)
				}
				return
			}
		}

		log.Printf("Database not ready yet, retrying in 2 seconds (attempt %d/%d)...", i+1, maxRetries)
		time.Sleep(2 * time.Second)
	}

	log.Fatal("Database not available after maximum retries")
}

func SetupDatabase() {
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "template_user")
	password := getEnv("DB_PASSWORD", "template_pass")
	dbname := getEnv("DB_NAME", "template_db")

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	// Initialize Ent client
	schema := getEnv("DB_SCHEMA", "template_service")
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable search_path=%s,public",
		host, port, user, password, dbname, schema)

	client, err := ent.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database with Ent: %v", err)
	}

	// Run database migrations (only in development)
	if getEnv("ENVIRONMENT", "") == "dev" {
		if err := client.Schema.Create(Ctx); err != nil {
			log.Printf("Warning: Schema creation error: %v", err)
		}
	}

	EntClient = client
	log.Println("Connected to database")
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
