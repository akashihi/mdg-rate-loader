package main

// [DB] Currency dictionary.
type CurrencyRecord struct {
	Id   uint64
	Code string
	Name string
}

func (CurrencyRecord) TableName() string {
	return "currency"
}
