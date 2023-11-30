package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

var (
	// prepared statement params is not working here
	// Could not execute query CREATE DATABASE IF NOT EXISTS $1; with params vw_page_test: pq: syntax error at or near "1"
	createDBQuery        = "CREATE DATABASE %s;"
	createUserQuery      = "CREATE USER %s;"
	grantPrivilegesQuery = "GRANT ALL ON DATABASE %s TO %s;"
)

func main() {
	dsn := "postgres://delivery:delivery@db:5432?sslmode=disable"
	dbName := "delivery"
	dbUser := "delivery"

	log.Println("Creating db and user")

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(errors.Wrapf(err, "Could not open sql connection to %s", dsn))
	}
	defer db.Close()

	//// Создаем пользователя
	if _, err := db.Exec(fmt.Sprintf(createUserQuery, dbUser)); err != nil {
		log.Fatal(errors.Wrapf(err, "Could not execute query %s with params %s", createUserQuery, dbUser))
	}

	// Создаем базу данных
	if _, err := db.Exec(fmt.Sprintf(createDBQuery, dbName)); err != nil {
		log.Fatal(errors.Wrapf(err, "Could not create database %s", dbName))
	}

	// Грантим привилегии
	if _, err := db.Exec(fmt.Sprintf(grantPrivilegesQuery, dbName, dbUser)); err != nil {
		log.Fatal(errors.Wrapf(err, "Could not execute query %s with params %s %s", grantPrivilegesQuery, dbName, dbUser))
	}

	log.Println("DB and user creation finished")
}
