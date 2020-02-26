package main

import (
	"database/sql"
	"errors"
)

type service struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Endpoint float64 `json:"endpoint"`
	Command  string  `json:"command"`
}

func (p *service) getService(db *sql.DB) error {
	return errors.New("Not implemented")
}

func (p *service) updateService(db *sql.DB) error {
	return errors.New("Not implemented")
}

func (p *service) deleteService(db *sql.DB) error {
	return errors.New("Not implemented")
}

func (p *service) createService(db *sql.DB) error {
	return errors.New("Not implemented")
}

func getServices(db *sql.DB, start, count int) ([]service, error) {
	return nil, errors.New("Not implemented")
}
