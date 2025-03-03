package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *sql.DB
var GORMDB *gorm.DB

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

	// Initialize GORM
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	gormDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error connecting to database with GORM: %v", err)
	}

	// Set default schema for GORM
	schema := getEnv("DB_SCHEMA", "template_service")
	err = gormDB.Exec(fmt.Sprintf("SET search_path TO %s, public", schema)).Error
	if err != nil {
		log.Fatalf("Error setting search path: %v", err)
	}

	GORMDB = gormDB
	log.Println("Connected to database")
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
