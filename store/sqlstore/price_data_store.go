package sqlstore

import (
	"context"
	"fmt"
	"strings"

	"github.com/adetunjii/ohlc/db/model"
)

type PriceDataStore struct {
	*SqlStore
}

func newPriceDataStore(sqlstore *SqlStore) *PriceDataStore {
	return &PriceDataStore{sqlstore}
}

const createPriceData = `
	INSERT INTO price_data (symbol, open, close, high, low, unix) 
	VALUES ($1, $2, $3, $4, $5, $6);
`

func (p *PriceDataStore) CreatePrice(ctx context.Context, arg model.PriceData) error {

	err := p.db.QueryRowContext(ctx, createPriceData, arg.Symbol, arg.Open, arg.Close, arg.High, arg.Low, arg.Unix).Err()
	if err != nil {
		return err
	}

	return nil
}

const listPrices = "SELECT id, open, close, low, high, unix, symbol, created_at FROM price_data"
const countRows = "SELECT COUNT(*) FROM price_data"

func (p *PriceDataStore) ListPrices(ctx context.Context, arg model.ListPriceParams) ([]model.PriceData, int64, error) {

	sql := listPrices

	if arg.Symbol != "" {
		sql = fmt.Sprintf(listPrices+" WHERE LOWER(symbol) LIKE '%%%s%%' ", strings.ToLower(arg.Symbol))
	}
	sql += `
	OFFSET $1
	LIMIT $2;` // end of query

	rows, err := p.db.QueryContext(ctx, sql, arg.Offset, arg.Limit)
	if err != nil {
		p.logger.Error("error querying rows", err)
		return nil, 0, err
	}

	prices := []model.PriceData{}

	for rows.Next() {
		var price model.PriceData
		err := rows.Scan(&price.ID, &price.Open, &price.Close, &price.Low, &price.High, &price.Unix, &price.Symbol, &price.CreatedAt)
		if err != nil {
			return nil, 0, err
		}
		prices = append(prices, price)
	}

	if err := rows.Close(); err != nil {
		return nil, 0, err
	}

	var count int64
	if err := p.db.QueryRowContext(ctx, countRows).Scan(&count); err != nil {
		p.logger.Error("error fetching total element count", err)
		return nil, 0, err
	}

	return prices, count, nil
}

func (p *PriceDataStore) BatchInsertPrice(ctx context.Context, batch []model.PriceData) error {

	tx, err := p.db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(createPriceData)
	if err != nil {
		p.logger.Error("error preparing transaction statement: ", err)
		return tx.Rollback()
	}
	defer stmt.Close()

	// create an interface array with a capacity of the batch length * amount of items per row
	values := make([][]interface{}, 0, len(batch)*6)

	for _, d := range batch {
		values = append(values, []interface{}{d.Symbol, d.Open, d.Close, d.High, d.Low, d.Unix})
	}

	for _, v := range values {
		_, err := stmt.Exec(v...)
		if err != nil {
			p.logger.Error("error executing transaction", err)
			return tx.Rollback()
		}
	}

	if err := tx.Commit(); err != nil {
		p.logger.Error("error commiting transaction", err)
		return err
	}

	return nil
}
