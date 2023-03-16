package store

import (
	"context"

	"github.com/adetunjii/ohlc/db/model"
)

type Sqlstore interface {
	PriceData() PriceDataStore
}

type PriceDataStore interface {
	ListPrices(ctx context.Context, arg model.ListPriceParams) ([]model.PriceData, int64, error)
	CreatePrice(ctx context.Context, arg model.PriceData) error
	BatchInsertPrice(ctx context.Context, arg []model.PriceData) error
	RetryFromDeadQueue(ctx context.Context) error
}
