package fetch

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEurBasedRates(t *testing.T) {
	ff := NewFiatFetcher()
	success, rates, _ := ff.Get("EUR", "", "")
	assert.Equal(t, true, success)
	assert.GreaterOrEqual(t, len(rates), 10)
	for _, v := range rates {
		fmt.Println(v.Symbol, v.Rate)
	}
}
