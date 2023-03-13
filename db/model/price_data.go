package model

import (
	"fmt"
	"time"
)

type PriceData struct {
	ID        int64     `json:"id"`
	Symbol    string    `json:"symbol"`
	Open      string    `json:"open"`
	High      string    `json:"high"`
	Low       string    `json:"low"`
	Close     string    `json:"close"`
	Unix      int64     `json:"unix"`
	CreatedAt time.Time `json:"created_at"`
}

func (p *PriceData) Validate() error {
	if p.Symbol == "" {
		return newInvalidPriceDataErr("symbol")
	}
	if p.Open == "" {
		return newInvalidPriceDataErr("open")
	}
	if p.High == "" {
		return newInvalidPriceDataErr("high")
	}
	if p.Low == "" {
		return newInvalidPriceDataErr("low")
	}
	if p.Close == "" {
		return newInvalidPriceDataErr("close")
	}
	if p.Unix == 0 {
		return newInvalidPriceDataErr("unix")
	}

	if p.Unix != 0 {
		p.Unix = int64(p.Unix)
	}

	return nil
}

type ListPriceParams struct {
	Limit  int32  `json:"limit"`
	Offset int32  `json:"offset"`
	Symbol string `json:"symbol"`
}

type InvalidPriceDataError struct {
	message string
	field   string
}

func newInvalidPriceDataErr(field string) *InvalidPriceDataError {
	message := fmt.Sprintf("model PriceData is invalid %s field.error", field)
	return &InvalidPriceDataError{
		message: message,
		field:   field,
	}
}

func (i *InvalidPriceDataError) Error() string {
	return i.message
}
