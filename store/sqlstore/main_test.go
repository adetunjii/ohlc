package sqlstore

import (
	"os"
	"testing"

	"github.com/adetunjii/ohlc/db"
	"github.com/adetunjii/ohlc/pkg/logging"
)

var testDB *db.PostgresDB
var sqlStore *SqlStore

// assumes that the database would have been created prior to this.
const databaseUrl = "postgresql://root:secret@localhost:5432/testdb?sslmode=disable"

func TestMain(m *testing.M) {

	dbConfig := &db.Config{
		DatabaseUrl: databaseUrl,
	}

	sugarLogger := logging.NewZapSugarLogger()
	logger := logging.NewLogger(sugarLogger)

	testDB = db.New(dbConfig, logger)
	sqlStore = New(testDB, logger)

	os.Exit(m.Run())

}
