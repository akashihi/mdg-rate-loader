package main

func main() {
	InitLog()
	log.Notice("Starting mdg currency rate loader...")

	configuration := config()

	db := newDbInterface(configuration)
	defer db.Close()
	service := newRateService(db)

	currencies, err := db.ListCurrencies()
	if err != nil {
		log.Error("Can't load currency list: %v", err)
	}

	for _, currency := range currencies {
		log.Debug("Loaded currency %d with code %s", currency.Id, currency.Code)
		for _, pair := range currencies {
			if pair == currency {
				continue
			}
			rate := RateRecord{From: currency.Id, FromCode: currency.Code, To: pair.Id, ToCode: pair.Code}
			newRate := getRate(&rate)
			if newRate != nil {
				service.store(newRate)
			}
		}
	}
}
