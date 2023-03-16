package sqlstore

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/adetunjii/ohlc/db/model"
)

type PriceDataStore struct {
	deadQueue chan []interface{} // TODO: move implementation to a message broker
	*SqlStore
}

func newPriceDataStore(sqlstore *SqlStore) *PriceDataStore {
	q := make(chan []interface{})
	return &PriceDataStore{q, sqlstore}
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
	var rows *sql.Rows
	var err error
	sql := listPrices

	// Writing the queries in this manner prevents against sql injection
	// The arguments are are sent seperately along with the prepared statement
	if arg.Symbol != "" {
		expr := "%" + arg.Symbol + "%"
		sql = fmt.Sprintf("%s WHERE LOWER(symbol) LIKE $1 OFFSET $1 LIMIT $2;", listPrices)

		rows, err = p.db.QueryContext(ctx, sql, expr, arg.Offset, arg.Limit)
		if err != nil {
			p.logger.Error("error querying rows", err)
			return nil, 0, err
		}
	} else {
		sql = fmt.Sprintf("%s OFFSET $1 LIMIT $2;", listPrices)
		rows, err = p.db.QueryContext(ctx, sql, arg.Offset, arg.Limit)
		if err != nil {
			p.logger.Error("error querying rows %v", err)
			return nil, 0, err
		}
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
		p.logger.Error("error initiating db transaction (%v)", err)
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(createPriceData)
	if err != nil {
		p.logger.Error("error preparing transaction statement: ", err)
		return err
	}
	defer stmt.Close()

	values := make([][]interface{}, 0, len(batch))
	for _, d := range batch {
		values = append(values, []interface{}{
			d.Symbol,
			d.Open,
			d.Close,
			d.High,
			d.Low,
			d.Unix,
		})
	}

	for _, v := range values {
		_, err := stmt.Exec(v...)
		if err != nil {
			p.logger.Error("error executing transaction (%v), sending data to dead queue...", err)
			// push the values into a deadqueue for retry
			p.deadQueue <- v
		}
	}

	if err := tx.Commit(); err != nil {
		p.logger.Error("error commiting transaction", err)
		return err
	}

	return nil
}

// RetryFromDeadQueue attempts to retry failed insertion of values
// into the database. It reads data from the deadQueue channel and
// inserts it into the database.
// For resiliency and fault tolerance, ideally this should read data
// from a messaging queue.
//
// If there is still a failure during insertion, the data is  written
// back into the queue so that it can be retried at
// some later time.
func (p *PriceDataStore) RetryFromDeadQueue(ctx context.Context) error {

	stmt, err := p.db.PrepareContext(ctx, createPriceData)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for item := range p.deadQueue {
		_, err := stmt.Exec(item...)
		if err != nil {
			p.logger.Error("error inserting from the dead queue (%v). Writing back to deadQueue...", err)
			// write back to the queue
			p.deadQueue <- item
			return err
		}
	}

	close(p.deadQueue)
	return nil
}
