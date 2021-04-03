# Fiat Rates
Get fiat rates, eg. EUR/USD, CNY/CAD. Updates daily.

Since `exchangeratesapi.io` limits free account to only `EUR` based rates and `1000` requests per month.

I create this open source project to give you a free format api.

# API

### GET a rate
```
GET /fiat?symbol=CNY&base=USD

0.123 // If not supported return -1

```

### GET supported pairs

```
GET /supported
```

### Refresh rates database (admin)

```
GET /refresh?password={...}
```