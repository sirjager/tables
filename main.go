package main

import (
	"database/sql"
	"log"
	"strings"

	api "github.com/SirJager/tables/api"
	"github.com/SirJager/tables/config"
	repo "github.com/SirJager/tables/service/core/repo"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func main() {
	c, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("Could not load server conifg :", err)
	}
	conn, conn_err := sql.Open(strings.Split(c.DBSource, ":")[0], c.DBSource)

	if conn_err != nil {
		panic(conn_err)
	}
	runDBMigration(c.MigrationURL, c.DBSource)

	store := repo.NewStore(conn)
	httpServer, err := api.NewHttpServer(store, conn, c)
	if err != nil {
		log.Fatal("Could not start http server:", err)
	}
	err = httpServer.Start(":" + c.Port)
	if err != nil {
		log.Fatal("Could not start http server:", err)
	}
}

func runDBMigration(migrationURL string, dbSource string) {
	migration, err := migrate.New(migrationURL, dbSource)
	if err != nil {
		log.Fatal("cannot create new migrate instance:", err)
	}

	if err = migration.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal("Failed to migrate database to latest version:", err)
	}

	log.Println("Database migration was successful")
}
