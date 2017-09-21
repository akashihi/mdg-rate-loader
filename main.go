package main

import "time"

func serve(db *DbInterface, service *RateService) {
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

func main() {
	InitLog()
	log.Notice("Starting mdg currency rate loader...")

	configuration := config()

	db := newDbInterface(configuration)
	defer db.Close()
	service := newRateService(db)

	for {
		serve(db, service)
		log.Info("Sleeping for %d minutes", configuration.Period)
		time.Sleep(time.Minute*time.Duration(configuration.Period))
	}
}
