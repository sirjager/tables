package core_repo

import (
	"database/sql"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/SirJager/tables/config"
)

var testQueries *Queries
var testDb *sql.DB

func TestMain(m *testing.M) {
	var err error
	cfg := config.ServerConfig{
		Port:                 os.Getenv("PORT"),
		GinMode:              os.Getenv("GIN_MODE"),
		DBSource:             os.Getenv("DATABASE_URL"),
		MigrationURL:         os.Getenv("MIGRATION_URL"),
		TokenSecretKey:       os.Getenv("ACCESS_TOKEN_DURATION"),
		AccessTokenDuration:  time.Minute * 3,
		RefreshTokenDuration: time.Minute * 5,
	}
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
