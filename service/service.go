package service

import (
	"github.com/adetunjii/ohlc/pkg/logging"
	"github.com/adetunjii/ohlc/store"
)

type Service struct {
	sqlstore store.Sqlstore
	logger   logging.Logger
}

func New(sqlstore store.Sqlstore, logger logging.Logger) *Service {
	return &Service{
		sqlstore: sqlstore,
		logger:   logger,
	}
}
