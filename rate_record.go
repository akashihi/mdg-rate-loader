package main

import (
	"time"

	"github.com/shopspring/decimal"
)

type RateRecord struct {
	Id        uint64    `sql:"AUTO_INCREMENT" gorm:"primary_key"`
	Beginning time.Time `gorm:"column:rate_beginning"`
	End       time.Time `gorm:"column:rate_end"`
	From      uint64    `gorm:"column:from_id"`
	To        uint64    `gorm:"column:to_id"`
	FromCode  string    `sql:"-"`
	ToCode    string    `sql:"-"`
	Rate      decimal.Decimal
}

func (RateRecord) TableName() string {
	return "rates"
}
