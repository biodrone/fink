package main

import (
	"database/sql"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

//App is the app
type App struct {
	Router *mux.Router
	DB     *sql.DB
}

//Initialize the App
func (a *App) Initialize(user, password, dbname string) {}

//Run the App
func (a *App) Run(addr string) {}
