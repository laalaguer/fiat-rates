# Fiat Rates
This repo is a project of fiat exchanges, eg. EUR/USD. Updates daily

Since `exchangeratesapi.io` limits free account to only EUR based.

And `1000` requests per month.

I create this open source project to give you a free format api.

# API

### Get a rate
Request
```
/latest?symbol=CNY&base=USD
```

Response
```json
{
    "symbol": "CNY",
    "base": "USD",
    "rate": 0.123
}
// If not supported return -1
```

### Get supported symbols
Request
```
/supported
```

Response
```json
{
    "symbols": ["CNY", "USD", "GBP", "HKD", ...]
}
```

### Refresh database (admin)
Request
```
/refresh?password={...}
```

Respones
```json
{
    "success": true
}
```