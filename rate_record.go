package main

import (
	"time"

	"github.com/shopspring/decimal"
)

type RateRecord struct {
	Id        uint64    `gorm:"primary_key, AUTO_INCREMENT"`
	Beginning time.Time `gorm:"column:rate_beginning"`
	End       time.Time `gorm:"column:rate_end"`
	From      uint64    `gorm:"column:from_id"`
	To        uint64    `gorm:"column:to_id"`
	FromCode  string    `gorm:"-"`
	ToCode    string    `gorm:"-"`
	Rate      decimal.Decimal
}
