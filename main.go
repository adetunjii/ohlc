package main

import (
	"github.com/adetunjii/ohlc/api"
	"github.com/adetunjii/ohlc/config"
	"github.com/adetunjii/ohlc/db"
	"github.com/adetunjii/ohlc/pkg/logging"
	"github.com/adetunjii/ohlc/service"
	"github.com/adetunjii/ohlc/store/sqlstore"
)

func main() {

	sugarLogger := logging.NewZapSugarLogger()
	logger := logging.NewLogger(sugarLogger)

	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Fatal("failed to load config%v", err)
	}

	dbConfig := &db.Config{
		Host:         cfg.DbHost,
		Port:         cfg.DbPort,
		User:         cfg.DbUser,
		Password:     cfg.DbPassword,
		DatabaseName: cfg.DbName,
		DatabaseUrl:  cfg.DbUrl,
	}

	database := db.New(dbConfig, logger)
	sqlStore := sqlstore.New(database, logger)

	svc := service.New(sqlStore, *logger)
	apiServer := api.NewServer(svc, *logger)

	apiServer.Start()

}
