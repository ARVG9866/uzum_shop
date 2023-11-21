package main

import (
	"log"

	"github.com/ARVG9866/uzum_shop/cmd/conf"

	"database/sql"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func main() {
	log.Println("starting migrate")
	cnf, err := conf.NewConfig()
	if err != nil {
		log.Fatal("failed to get config", err.Error())
	}

	sqlConnectionString := conf.GetSqlConnectionString(cnf)

	db, err := sql.Open("postgres", sqlConnectionString)
	if err != nil {
		log.Fatal("failed to opening connection to db: ", err.Error())
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatal("failed to get WithInstance", err.Error())
	}

	m, err := migrate.NewWithDatabaseInstance("file://migrations", "postgres", driver)
	if err != nil {
		log.Fatal("failed to get NewWithDatabaseInstance", err.Error())
	}

	err = m.Up()
	if err != nil {
		log.Fatal("failed to migrate up", err.Error())
	}

	log.Println("success migrated")
}
