package main

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	_ "github.com/pjankovsky/gdriver-go/migration"
	"github.com/pressly/goose"
	"log"
)

func setupSQL() {
	db := getDB()

	err :=goose.SetDialect("sqlite3")
	if err != nil {
		log.Fatalf("Unable to set goose dialect: %v", err)
	}

	err = goose.Up(db.DB, "migration")
	if err != nil {
		log.Fatalf("Unable to run goose migrations: %v", err)
	}
}

var cachedDB *sqlx.DB

func getDB() *sqlx.DB {
	if cachedDB == nil {
		db, err := sqlx.Connect("sqlite3", settings.SQLitePath)
		if err != nil {
			log.Fatalf("Unable to connect to sqlite3 db: %v", err)
		}
		cachedDB = db
	}
	return cachedDB
}
