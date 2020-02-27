package main

import (
	"database/sql"
)

type service struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Endpoint float64 `json:"endpoint"`
	Command  string  `json:"command"`
}

func getServices(db *sql.DB, start, count int) ([]service, error) {
	rows, err := db.Query(
		"SELECT id, name, endpoint, command FROM services LIMIT $1 OFFSET $2",
		count, start)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	services := []service{}

	for rows.Next() {
		var p service
		if err := rows.Scan(&p.ID, &p.Name, &p.Endpoint, &p.Command); err != nil {
			return nil, err
		}
		services = append(services, p)
	}

	return services, nil
}

func (p *service) getService(db *sql.DB) error {
	return db.QueryRow("SELECT name, endpoint FROM services WHERE id=$1",
		p.ID).Scan(&p.Name, &p.Endpoint, &p.Command)
}

func (p *service) updateService(db *sql.DB) error {
	_, err :=
		db.Exec("UPDATE services SET name=$1, endpoint=$2, command=$3 WHERE id=$4",
			p.Name, p.Endpoint, &p.Command, p.ID)

	return err
}

func (p *service) deleteService(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM services WHERE id=$1", p.ID)

	return err
}

func (p *service) createService(db *sql.DB) error {
	err := db.QueryRow(
		"INSERT INTO services(name, endpoint, command) VALUES($1, $2, $3) RETURNING id",
		p.Name, p.Endpoint, &p.Command).Scan(&p.ID)

	if err != nil {
		return err
	}

	return nil
}
