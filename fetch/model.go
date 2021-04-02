package fetch

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

type PricePoint struct {
	Symbol string  `json:"symbol"`
	Rate   float64 `json:"rate"`
}

type ExchangeRatesAPIRespones struct {
	Success   bool   `json:"success"`
	Timestamp int64  `json:"timestamp"`
	Base      string `json:"base"`
	Date      string `json:"date"`
}

type Err struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type ExchangeRatesAPIError struct {
	Error Err `json:"error"`
}

func makeProxy(proxyURL string) *http.Transport {
	if len(proxyURL) != 0 {
		proxyURL, err := url.Parse(proxyURL)
		if err != nil {
			return nil
		}
		transport := &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		}
		return transport
	}
	return nil
}

// Scrap shall scrap the api and return the []byte of result
// if nil then scrap failed.
func Scrap(url string, proxyURL string, headers map[string]string) (b []byte) {
	transport := makeProxy(proxyURL)
	var httpClient http.Client

	if transport != nil {
		httpClient = http.Client{
			Timeout:   time.Second * 5, // Maximum of 5 secs
			Transport: transport,
		}
	} else {
		httpClient = http.Client{
			Timeout: time.Second * 5, // Maximum of 5 secs
		}
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	if headers != nil {
		for key, value := range headers {
			req.Header.Set(key, value)
		}
	}

	res, getErr := httpClient.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
		return nil
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
		return nil
	}

	return body
}
