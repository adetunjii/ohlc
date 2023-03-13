package sqlstore

import (
	"github.com/adetunjii/ohlc/db"
	"github.com/adetunjii/ohlc/pkg/logging"
	"github.com/adetunjii/ohlc/store"
)

type Stores struct {
	priceData store.PriceDataStore
}

type SqlStore struct {
	db     *db.PostgresDB
	logger *logging.Logger
	stores Stores
}

var _ store.Sqlstore = (*SqlStore)(nil)

func New(db *db.PostgresDB, logger *logging.Logger) *SqlStore {
	sqlstore := &SqlStore{
		logger: logger,
		db:     db,
	}

	sqlstore.stores.priceData = newPriceDataStore(sqlstore)
	return sqlstore
}

func (s *SqlStore) PriceData() store.PriceDataStore {
	return s.stores.priceData
}
