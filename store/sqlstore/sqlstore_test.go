package sqlstore

import (
	"context"
	"testing"

	"github.com/adetunjii/ohlc/db/model"
	"github.com/stretchr/testify/require"
)

func TestListPriceData(t *testing.T) {

	listPriceParam := model.ListPriceParams{
		Offset: 0,
		Limit:  20,
		Symbol: "BTCUSDT",
	}

	price_data, _, err := sqlStore.PriceData().ListPrices(context.Background(), listPriceParam)
	require.NoError(t, err)
	require.NotEmpty(t, price_data)
	require.Len(t, price_data, 20)

}

func TestCreatePrice(t *testing.T) {
	priceData := model.PriceData{
		Symbol: "BTCUSDT",
		Open:   "42123.29000000",
		High:   "42148.32000000",
		Low:    "42120.82000000",
		Close:  "42146.06000000",
		Unix:   1644719700000,
	}

	err := sqlStore.PriceData().CreatePrice(context.Background(), priceData)
	require.NoError(t, err)
}

var testBatchData = []model.PriceData{
	{Symbol: "BTCUSDT", Open: "42123.29000000", High: "42148.32000000", Low: "42120.82000000", Close: "42146.06000000", Unix: 1644719700000},
	{Symbol: "BTCUSDT", Open: "42123.29000000", High: "42148.32000000", Low: "42120.82000000", Close: "42146.06000000", Unix: 1644719700000},
	{Symbol: "BTCUSDT", Open: "42123.29000000", High: "42148.32000000", Low: "42120.82000000", Close: "42146.06000000", Unix: 1644719700000},
	{Symbol: "BTCUSDT", Open: "42123.29000000", High: "42148.32000000", Low: "42120.82000000", Close: "42146.06000000", Unix: 1644719700000},
	{Symbol: "BTCUSDT", Open: "42123.29000000", High: "42148.32000000", Low: "42120.82000000", Close: "42146.06000000", Unix: 1644719700000},
	{Symbol: "BTCUSDT", Open: "42123.29000000", High: "42148.32000000", Low: "42120.82000000", Close: "42146.06000000", Unix: 1644719700000},
	{Symbol: "BTCUSDT", Open: "42123.29000000", High: "42148.32000000", Low: "42120.82000000", Close: "42146.06000000", Unix: 1644719700000},
	{Symbol: "BTCUSDT", Open: "42123.29000000", High: "42148.32000000", Low: "42120.82000000", Close: "42146.06000000", Unix: 1644719700000},
	{Symbol: "BTCUSDT", Open: "42123.29000000", High: "42148.32000000", Low: "42120.82000000", Close: "42146.06000000", Unix: 1644719700000},
	{Symbol: "BTCUSDT", Open: "42123.29000000", High: "42148.32000000", Low: "42120.82000000", Close: "42146.06000000", Unix: 1644719700000},
	{Symbol: "BTCUSDT", Open: "42123.29000000", High: "42148.32000000", Low: "42120.82000000", Close: "42146.06000000", Unix: 1644719700000},
	{Symbol: "BTCUSDT", Open: "42123.29000000", High: "42148.32000000", Low: "42120.82000000", Close: "42146.06000000", Unix: 1644719700000},
	{Symbol: "BTCUSDT", Open: "42123.29000000", High: "42148.32000000", Low: "42120.82000000", Close: "42146.06000000", Unix: 1644719700000},
	{Symbol: "BTCUSDT", Open: "42123.29000000", High: "42148.32000000", Low: "42120.82000000", Close: "42146.06000000", Unix: 1644719700000},
	{Symbol: "BTCUSDT", Open: "42123.29000000", High: "42148.32000000", Low: "42120.82000000", Close: "42146.06000000", Unix: 1644719700000},
	{Symbol: "BTCUSDT", Open: "42123.29000000", High: "42148.32000000", Low: "42120.82000000", Close: "42146.06000000", Unix: 1644719700000},
	{Symbol: "BTCUSDT", Open: "42123.29000000", High: "42148.32000000", Low: "42120.82000000", Close: "42146.06000000", Unix: 1644719700000},
	{Symbol: "BTCUSDT", Open: "42123.29000000", High: "42148.32000000", Low: "42120.82000000", Close: "42146.06000000", Unix: 1644719700000},
	{Symbol: "BTCUSDT", Open: "42123.29000000", High: "42148.32000000", Low: "42120.82000000", Close: "42146.06000000", Unix: 1644719700000},
}

func TestBatchInsertPrice(t *testing.T) {
	err := sqlStore.PriceData().BatchInsertPrice(context.Background(), testBatchData)

	require.NoError(t, err)
}
