package main

import (
	"fmt"

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
func newDbInterface(configuration Configuration) *DbInterface {
	dbUrl := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s", configuration.Host, configuration.User, configuration.Database, configuration.Password)

	db, err := gorm.Open("postgres", dbUrl)
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

// Retrieve exterior Rate
func (dbintf *DbInterface) FindExterior(rate *RateRecord) (*RateRecord) {
	var exterior []RateRecord
	dbintf.db.Where("rate_beginning <= ? and rate_end > ? and from_id = ? and to_id = ?", rate.Beginning, rate.End, rate.From, rate.To).First(&exterior)
	if len(exterior) == 0 {
		return nil
	}
	return &exterior[0]
}

// Retrieve following Rate
func (dbintf *DbInterface) FindFollowing(rate *RateRecord) (*RateRecord) {
	var following []RateRecord
	dbintf.db.Where("rate_beginning >= ? and from_id = ? and to_id = ?", rate.End, rate.From, rate.To).First(&following)
	if len(following) == 0 {
		return nil
	}
	return &following[0]
}

// Saves rate entity to the database
func (dbintf *DbInterface) SaveRate(rate *RateRecord) (error) {
	return dbintf.db.Create(rate).Error
}

// Updates rate entity in the database
func (dbintf *DbInterface) UpdateRate(rate *RateRecord) (error) {
	return dbintf.db.Save(rate).Error
}
