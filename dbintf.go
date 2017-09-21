package main

import (

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"os"
)

// Data access layer
type DbInterface struct {
	db *gorm.DB
}

// Prepares database for access.
// Quite limited in tuning
func newDbInterface() *DbInterface {
	db, err := gorm.Open("postgres", "host=localhost user=mdg dbname=mdg sslmode=disable password=mdg")
	if err != nil {
		log.Fatalf("Failed to connect database: %v", err)
		os.Exit(1)
	}

	log.Notice("Database connection ready")

	return &DbInterface{db: db}
}

// Closes database connection
func (dbintf *DbInterface) Close() {
	log.Notice("Closing database connection")
	dbintf.db.Close()
}

// Retrieves list of currencies
func (dbintf *DbInterface) ListCurrencies() ([]CurrencyRecord, error) {
	var currencies []CurrencyRecord
	err := dbintf.db.Find(&currencies).Error
	return currencies, err
}

