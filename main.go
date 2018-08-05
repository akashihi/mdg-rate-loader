package main

import "time"

func findPair(pairs []RateRecord, fromCode string, toCode string) *RateRecord {
	for _, v := range pairs {
		if (v.FromCode == fromCode && v.ToCode == toCode) {
			return &v
		}
	}
	return nil
}

func serve(db *DbInterface, service *RateService) {
	currencies, err := db.ListCurrencies()
	if err != nil {
		log.Errorf("Can't load currency list: %v", err)
	}

	var ratePairs []RateRecord
	var missingRates []RateRecord

	for _, currency := range currencies {
		log.Debugf("Loaded currency %d with code %s", currency.Id, currency.Code)
		for _, pair := range currencies {
			if pair == currency {
				continue
			}
			rate := RateRecord{From: currency.Id, FromCode: currency.Code, To: pair.Id, ToCode: pair.Code}		
			newRate := getRate(&rate)
			if newRate != nil {
				ratePairs = append(ratePairs, *newRate)
			} else {
				missingRates = append(missingRates, rate)
			}
		}
	}

	for _, rate := range missingRates {
		log.Warningf("Trying to find USD based crossrate for %s%s pair", rate.FromCode, rate.ToCode)
		left := findPair(ratePairs, rate.FromCode, "USD")
		if left == nil {
			log.Warningf("No conversion rate to USD from %s, skipping", rate.FromCode)
			continue
		}
		right := findPair(ratePairs, "USD", rate.ToCode)
		if right == nil {
			log.Warningf("No conversion rate from USD to %s, skipping", rate.ToCode)
			continue
		}
		rate.Rate = left.Rate.Mul(right.Rate)
		ratePairs = append(ratePairs, rate)
		log.Noticef("Calculated %s%s crossrate: %s", rate.FromCode, rate.ToCode, rate.Rate)
	}
	
	for _, rate := range ratePairs {
		service.store(&rate)
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
		log.Infof("Sleeping for %d minutes", configuration.Period)
		time.Sleep(time.Minute*time.Duration(configuration.Period))
	}
}
