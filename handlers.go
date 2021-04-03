package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/laalaguer/fiat-rates/cache"
	"github.com/laalaguer/fiat-rates/fetch"
)

var EXCHANGE_NAME = "exchangeratesapi"
var BASE = "EUR"
var priceCache = cache.NewPriceDataCache(1000)

// RefreshPriceHandler refreshes the cache with latest price from upstream
func RefreshPriceHandler(w http.ResponseWriter, r *http.Request, proxyURL string, adminPassword string, remoteAPIKey string) {
	// Get parameters from url and form.
	err := r.ParseForm()
	if err != nil {
		fmt.Println(fmt.Errorf("%v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// Check admin password
	pwd := r.Form.Get("password")
	if pwd != adminPassword {
		return
	}
	ff := fetch.NewFiatFetcher()
	success, rates, response := ff.Get(BASE, remoteAPIKey, proxyURL)
	if !success {
		w.WriteHeader(http.StatusBadRequest)
	}
	for _, v := range rates {
		pdd := cache.PriceData{
			Symbol:       v.Symbol,
			Base:         response.Base,
			Price:        v.Rate,
			DatePoint:    time.Now().UTC(),
			ExchangeName: EXCHANGE_NAME}
		priceCache.AddOrUpdateData(pdd)
	}

	b, err := json.Marshal(priceCache.GetAll())
	if err != nil {
		fmt.Println(fmt.Errorf("%v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(b)
}

func SupportedSymbolsHandler(w http.ResponseWriter, r *http.Request) {
	prices := priceCache.GetAll()

	s := make([]string, 0, len(prices)+1)
	for _, v := range prices {
		s = append(s, v.Symbol)
	}
	s = append(s, BASE) // Don't forget to append the BASE!

	b, err := json.Marshal(s)
	if err != nil {
		fmt.Println(fmt.Errorf("%v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(b)
}

// FiatPriceHandler get you a pair of rates with format of /endpoint?symbol=CNY&base=USD
func FiatPriceHandler(w http.ResponseWriter, r *http.Request, proxyURL string) {
	// Get parameters from url and form.
	err := r.ParseForm()
	if err != nil {
		fmt.Println(fmt.Errorf("%v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	symbol := r.Form.Get("symbol")
	if symbol == "" {
		symbol = "CNY"
	}
	symbol = strings.ToUpper(symbol)

	base := r.Form.Get("base")
	if base == "" {
		base = "USD"
	}
	base = strings.ToUpper(base)

	var price float64 = -1
	// Find price in cache, or return -1
	if base == BASE {
		p := priceCache.GetData(symbol, BASE, EXCHANGE_NAME)
		if p != nil {
			price = p.Price
		}
	} else {
		p1 := priceCache.GetData(symbol, BASE, EXCHANGE_NAME)
		p2 := priceCache.GetData(base, BASE, EXCHANGE_NAME)
		if p1 != nil && p2 != nil {
			price = p1.Price / p2.Price
		}
	}

	b, err := json.Marshal(price)
	if err != nil {
		fmt.Println(fmt.Errorf("%v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(b)
}
