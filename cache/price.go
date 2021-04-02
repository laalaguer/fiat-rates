// Price cache
// A temp solution for not having a memcached/redis server.

package cache

import (
	"fmt"
	"strings"
	"time"

	lru "github.com/hashicorp/golang-lru"
)

// Make a key.
func makeKey(s string, b string, exchange string) string {
	mList := []string{s, b, exchange}
	temp := strings.Join(mList, ",")
	return strings.ToUpper(temp)
}

// PriceData the structure of price and according date.
type PriceData struct {
	Symbol       string    `json:"symbol"`
	Base         string    `json:"base"`
	Price        float64   `json:"price"`
	DatePoint    time.Time `json:"date_point"`
	ExchangeName string    `json:"exchange"`
}

// PriceDataCache is a cache for prices.
type PriceDataCache struct {
	Prices *lru.Cache
}

// NewPriceDataCache creates a new cache.
func NewPriceDataCache(i int) *PriceDataCache {
	pdc := PriceDataCache{}
	pdc.Prices, _ = lru.New(i)
	return &pdc
}

// AddOrUpdateData add element or update element in cache.
// Returns true on success update/add else false(capacity full).
func (pdc *PriceDataCache) AddOrUpdateData(pd PriceData) (b bool) {
	key := makeKey(pd.Symbol, pd.Base, pd.ExchangeName)
	found := pdc.Prices.Contains(key)

	// old data. Remove first.
	if found {
		fmt.Println("found!")
		pdc.Prices.Remove(key)
	}

	pdc.Prices.Add(key, pd)
	return true
}

// GetData returns the data object itself or nil.
func (pdc *PriceDataCache) GetData(s string, b string, exchange string) (pd *PriceData) {
	item, ok := pdc.Prices.Get(makeKey(s, b, exchange))
	if !ok {
		return nil
	}
	fmt.Println(item)
	a := item.(PriceData)
	return &a
}

// GetOnlyFreshData returns the data object itself or nil.
func (pdc *PriceDataCache) GetOnlyFreshData(s string, b string, exchange string, seconds time.Duration) (pd *PriceData) {
	data := pdc.GetData(s, b, exchange)
	// if not found.
	if data == nil {
		return nil
	}
	// still fresh
	threshold := data.DatePoint.Add(time.Second * seconds)

	if time.Now().UTC().After(threshold) {
		return nil
	}

	return data
}

// GetAll returns all the elements in the cache.
func (pdc *PriceDataCache) GetAll() []PriceData {
	keys := pdc.Prices.Keys()
	// total := len(keys)
	ret := make([]PriceData, 0)

	for _, key := range keys {
		value, ok := pdc.Prices.Peek(key)
		if ok {
			ret = append(ret, value.(PriceData))
		}
	}
	return ret
}

// Len returns the number of items in cache.
func (pdc *PriceDataCache) Len() (i int) {
	return pdc.Prices.Len()
}
