package main

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/NodePrime/jsonpath"
	"github.com/shopspring/decimal"
)

func getRate(rate *RateRecord) (*RateRecord){
	var netClient = &http.Client{
		Timeout: time.Second * 10,
	}

	pair := fmt.Sprintf("%s%s=X", rate.FromCode, rate.ToCode)
	log.Debug("Loading rate for %s pair", pair)

	u, err := url.Parse("https://query1.finance.yahoo.com/v7/finance/spark?range=1d&interval=60m&indicators=close&includeTimestamps=true&includePrePost=false&corsDomain=finance.yahoo.com&.tsrc=finance")
	if err != nil {
		log.Warning("Unable to construct YahooFinance url: %v", err)
		return nil
	}
	q := u.Query()
	q.Set("symbols", pair)
	u.RawQuery = q.Encode()


	response, err := netClient.Get(u.String())
	if err != nil {
		log.Warning("Unable to retrieve ")
	}

	ratePath := "$.spark.result[0].response[0].indicators.quote[0].close[*]+"
	paths, err := jsonpath.ParsePaths(ratePath)
	if err != nil {
		log.Critical("Invalid JsonPath: %v", err)
		return nil
	}
	eval, err := jsonpath.EvalPathsInReader(response.Body, paths)
	if err != nil {
		log.Warning("Unable to parse YahooFinance json: %v", err)
		return nil
	}
	var rates []string
	for {
		if result, ok := eval.Next(); ok {
			val := string(result.Value[:])
			if val != "null" {
				rates = append(rates, val)
			}
		} else {
			break
		}
	}

	if (len(rates) == 0) {
		log.Warning("No rates returned for %s pair", pair)
		return nil
	}
	rate.Rate, err = decimal.NewFromString(rates[len(rates)-1])
	if err != nil {
		log.Warning("Unable to parse rate value: %v", err)
		return nil
	}

	log.Notice("Rate for %s%s pair is %s", rate.FromCode, rate.ToCode, rate.Rate.String())
	return rate
}

