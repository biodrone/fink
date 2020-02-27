package main

import (
	"log"
	"os"
	"testing"
)

var a App

const tableCreationQuery = `CREATE TABLE IF NOT EXISTS services
(
	id SERIAL,
	name TEXT NOT NULL,
	endpoint TEXT NOT NULL,
	command TEXT NOT NULL,
	CONSTRAINT services_pkey PRIMARY KEY (id)
)`

func TestMain(m *testing.M) {
	a = App{}
	a.Initialize(
		os.Getenv("TEST_DB_USERNAME"),
		os.Getenv("TEST_DB_PASSWORD"),
		os.Getenv("TEST_DB_NAME"),
	)
	ensureTableExists()
	code := m.Run()
	clearTable()
	os.Exit(code)
}

func ensureTableExists() {
	if _, err := a.DB.Exec(tableCreationQuery); err != nil {
		log.Fatal(err)
	}
}

func clearTable() {
	a.DB.Exec("DELETE FROM services")
	a.DB.Exec("ALTER SEQUENCE services_id_seq RESTART WITH 1")
}
