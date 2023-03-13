package db

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/adetunjii/ohlc/pkg/logging"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type Config struct {
	Host         string `json:"host"`
	Port         string `json:"port"`
	User         string `json:"user"`
	DatabaseName string `json:"database_name"`
	Password     string `json:"password"`
	DatabaseUrl  string `json:"url"`
}

type PostgresDB struct {
	logger *logging.Logger

	*sql.DB
}

func New(dbConfig *Config, logger *logging.Logger) *PostgresDB {
	db := &PostgresDB{
		DB:     nil,
		logger: logger,
	}

	if err := db.Connect(dbConfig); err != nil {
		logger.Fatal("connection to db failed: %v", err)
	}
	return db
}

func (p *PostgresDB) Connect(config *Config) error {

	var dsn string
	databaseUrl := config.DatabaseUrl

	if databaseUrl == "" {
		dsn = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
			config.User, config.Password, config.Host, config.Port, config.DatabaseName)
	} else {
		dsn = databaseUrl
	}

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return err
	}

	if err := db.Ping(); err != nil {
		return errors.New("could not ping the database")
	}

	p.DB = db

	p.logger.Info(fmt.Sprintf("Database Connected Successfully %v...", dsn))
	return nil
}

func (p *PostgresDB) CloseConnection() error {
	return p.DB.Close()
}
