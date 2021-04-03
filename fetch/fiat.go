package fetch

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

// FiatFetcher is a fetcher that can convert different currencies via internet.
type FiatFetcher struct{}

// NewFiatFetcher creates a FiatPriceFetcher
func NewFiatFetcher() (ff *FiatFetcher) {
	return &FiatFetcher{}
}

// Get gets the prices of all XXX/base price
// It also returns the response object
func (ff *FiatFetcher) Get(base string, access_key string, proxyURL string) (success bool, rates []PricePoint, response ExchangeRatesAPIRespones) {
	url := fmt.Sprintf(
		"http://api.exchangeratesapi.io/v1/latest?base=%s&access_key=%s",
		strings.ToUpper(base),
		access_key)
	body := Scrap(url, proxyURL, nil)

	// Scrap failed. HTTP error
	if body == nil {
		return false, nil, ExchangeRatesAPIRespones{Success: false}
	}

	// Json Unmarshall Error
	var responseObj interface{}
	jsonErr := json.Unmarshal(body, &responseObj)
	data := responseObj.(map[string]interface{})
	if jsonErr != nil {
		log.Fatal(jsonErr)
		return false, nil, ExchangeRatesAPIRespones{Success: false}
	}

	// Response contains error message
	if _, ok := data["error"]; ok {
		var r ExchangeRatesAPIError
		jsonErr := json.Unmarshal(body, &r)
		if jsonErr != nil {
			log.Fatal(jsonErr)
		}
		log.Fatal(fmt.Sprintf("error message: %s %s", r.Error.Code, r.Error.Message))
		return false, nil, ExchangeRatesAPIRespones{Success: false}
	}

	// No error? Let's unmarshal the result
	var r ExchangeRatesAPIRespones
	jsonErr = json.Unmarshal(body, &r)
	if jsonErr != nil {
		log.Fatal(jsonErr)
		return false, nil, ExchangeRatesAPIRespones{Success: false}
	}

	// Get all the price points
	l := make([]PricePoint, 0, 300)

	for k, v := range data {
		if k == "rates" {
			for i, u := range v.(map[string]interface{}) {
				x := PricePoint{Symbol: strings.ToUpper(i), Rate: u.(float64)}
				l = append(l, x)
			}
		}
	}

	return r.Success, l, r
}
