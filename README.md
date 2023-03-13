# OHLC

---

This application stores large amounts of historical OHLC price data that is shared in CSV file format. The whole idea of the application is to centralize and digitalize this historical data in the most efficient and performant means possible. These files can be ranging from a few GBs to a couple of TBs.

### Endpoints

**PRICE LIST - GET**

`/api/v1/price-data/list`

**CREATE PRICE - POST**

`/api/v1/price-data/create`

Request Body Sample

> {
> "unix": 1644719700000,
> "symbol": "BTCUSDT”,
> "open":"42123.29000000",
> "high": "42148.32000000",
> "low": "42120.82000000",
> "close":"42146.06000000"
> }

`/api/v1/price-data/upload-price-list`
Takes a csv file with a maximum file size of 5TB in the format UNIX,SYMBOL,OPEN,HIGH,LOW,CLOSE

**RUNNING THE PROJECT**

To run the project in a local environment, you would need to set the environment variables below in a file called `app.env`

```
DB_HOST=localhost
DB_PORT=5432
DB_USER=root
DB_PASSWORD=secret
DB_NAME=ohlc
DB_URL="postgresql://root:secret@localhost:5432/ohlc?sslmode=disable"
```

**TO RUN IN DOCKER**
` docker compose up`
