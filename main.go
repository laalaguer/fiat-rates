// Cache: online

// Refresh:
// Refresh from upstream
// stuff a bunch of "CNY"/"EUR" into the price cache list
// remember the base is "EUR"

// When a price query hits, for example "CNY"/"GBP"
// Look for the base
// Find "CNY"/"EUR" and "GBP"/"EUR"
// Then output "CNY"/"GBP"

// When a supported hits,
// List all the caches entries
// extract the symbols
// extract the base
// return [...symbols, base]

// When a refresh hits,
// see refresh
