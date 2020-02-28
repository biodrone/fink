package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

//App is the app
type App struct {
	Router *mux.Router
	DB     *sql.DB
}

//Initialise the App
func (a *App) Initialise(user, password, dbname string) {
	connectionString :=
		fmt.Sprintf("user=%s password=%s dbname=%s", user, password, dbname)

	var err error
	a.DB, err = sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	a.Router = mux.NewRouter()
	a.initializeRoutes()
}

//Run the App
func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(":8000", a.Router))
}

func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/services", a.getServices).Methods("GET")
	a.Router.HandleFunc("/service", a.createService).Methods("POST")
	a.Router.HandleFunc("/service/{id:[0-9]+}", a.getService).Methods("GET")
	a.Router.HandleFunc("/service/{id:[0-9]+}", a.updateService).Methods("PUT")
	a.Router.HandleFunc("/service/{id:[0-9]+}", a.deleteService).Methods("DELETE")
}
