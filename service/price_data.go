package service

import (
	"context"
	"encoding/csv"
	"io"
	"strconv"

	"github.com/adetunjii/ohlc/db/model"
)

const (
	maxNumOfGoroutines = 10
	batchSize          = 100 // maximum number of rows
)

func (s *Service) ListPrices(ctx context.Context, arg model.ListPriceParams) ([]model.PriceData, int64, error) {
	return s.sqlstore.PriceData().ListPrices(ctx, arg)
}

func (s *Service) CreatePrice(ctx context.Context, arg model.PriceData) error {
	return s.sqlstore.PriceData().CreatePrice(ctx, arg)
}

// BulkInsertPrice reads the file part from r which is a stream
// from the original file. A csv reader reads the file part and
// splits the data into batches of size `batchSize` which is sent
// to a channel (batchChan) to be picked up by n (numOfGoroutines)
// go routines for insertion. The Goroutines bulk Insert each batch
// into the database

func (s *Service) BulkInsertPrice(ctx context.Context, r io.Reader) error {
	reader := csv.NewReader(r)
	errChan := make(chan error)
	batchChan := make(chan []model.PriceData)

	// read the file parts and split them into batches of `batchSize`
	go func() {
		for {
			batch, err := createBatch(reader)
			if err != nil {
				errChan <- err
			}
			batchChan <- batch
		}
	}()
	done := make(chan struct{})

	// ctx, cancel := context.WithDeadline(ctx, time.Now().Add(1*time.Minute))
	// defer cancel()

	// worker goroutines pick up the data in batches and batch insert it into the database
	for i := 0; i < maxNumOfGoroutines; i++ {
		go func() {

			for batch := range batchChan {
				if batch == nil {
					break
				}

				err := s.sqlstore.PriceData().BatchInsertPrice(ctx, batch)
				if err != nil {
					errChan <- err
				}
			}
			done <- struct{}{}
		}()
	}

	// wait for all the goroutines to finish
	for {
		select {
		case err := <-errChan:
			if err != nil {
				s.logger.Error("error uploading data (%v)", err)
				return err
			}
		case <-done:
			return nil
			// case <-ctx.Done():
			// 	return errors.New("file did not upload completely")
		}
	}
}

func (s *Service) RetryFromDeadQueue(ctx context.Context) error {
	return s.sqlstore.PriceData().RetryFromDeadQueue(ctx)
}

// converts each row in the csv file to type model.PriceData
func processRow(row []string) (price model.PriceData, err error) {
	unix_val, err := strconv.Atoi(row[0])
	if err != nil {
		return
	}

	price.Unix = int64(unix_val)
	price.Symbol = row[1]
	price.Open = row[2]
	price.High = row[3]
	price.Low = row[4]
	price.Close = row[5]

	if err := price.Validate(); err != nil {
		return model.PriceData{}, err
	}
	return
}

// createBatch reads data from reader and returns a batch of 100 rows
func createBatch(reader *csv.Reader) ([]model.PriceData, error) {
	var batch []model.PriceData
	for i := 0; i < batchSize; i++ {
		row, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		priceData, err := processRow(row)
		if err != nil {
			return nil, err
		}
		batch = append(batch, priceData)
	}
	return batch, nil
}
