package database

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/alimasry/gopad/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	db      *gorm.DB
	dbError error
	once    sync.Once
)

// creates the database tables and run migrations
func migrateTables(db *gorm.DB) {
	if err := db.AutoMigrate(&models.Document{}); err != nil {
		log.Fatalf("Can't create table documents: %v", err)
	}
	if err := db.AutoMigrate(&models.DocumentVersion{}); err != nil {
		log.Fatalf("Can't create table document_version: %v", err)
	}
}

// initialize the database
func initializeDb() {
	// Get database credentials from environment variables
	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbPortStr := os.Getenv("DB_PORT")

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		dbHost,
		dbUser,
		dbPassword,
		dbName,
		dbPortStr,
	)

	db, dbError = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if dbError != nil {
		log.Fatal(dbError)
	}
}

// get the database and initialize it if it's not
func GetDb() *gorm.DB {
	once.Do(initializeDb)
	return db
}

// creates the tables and does the migrations
func Init() error {
	db := GetDb()
	migrateTables(db)
	return dbError
}
