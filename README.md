# Fiat Rates
Get fiat rates, eg. EUR/USD. Updates daily

Since `exchangeratesapi.io` limits free account to only EUR based.

And `1000` requests per month.

I create this open source project to give you a free format api.

# API

### Get a rate
```
GET /fiat?symbol=CNY&base=USD

{
    "symbol": "CNY",
    "base": "USD",
    "rate": 0.123 // If not supported return -1
}
```

### Get supported symbols

```
GET /supported
```

### Refresh database (admin)

```
GET /refresh?password={...}

{
    "success": true // failed then false
}
```