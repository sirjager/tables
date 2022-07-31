package core_repo

import (
	"database/sql"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/SirJager/tables/config"
)

var testQueries *Queries
var testDb *sql.DB

func TestMain(m *testing.M) {
	var err error
	cfg, err := config.LoadConfig("../../../")
	if err != nil {
		log.Fatal("unable to load environment variables:", err)
	}
	testDb, err = sql.Open(strings.Split(cfg.DBSource, ":")[0], cfg.DBSource)
	if err != nil {
		log.Fatal("cannot connect to database:", err)
	}
	testQueries = New(testDb)
	os.Exit(m.Run())
}
