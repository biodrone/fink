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
		var s service
		if err := rows.Scan(&s.ID, &s.Name, &s.Endpoint, &s.Command); err != nil {
			return nil, err
		}
		services = append(services, s)
	}

	return services, nil
}

func (s *service) getService(db *sql.DB) error {
	return db.QueryRow("SELECT name, endpoint FROM services WHERE id=$1",
		s.ID).Scan(&s.Name, &s.Endpoint, &s.Command)
}

func (s *service) updateService(db *sql.DB) error {
	_, err :=
		db.Exec("UPDATE services SET name=$1, endpoint=$2, command=$3 WHERE id=$4",
			s.Name, s.Endpoint, &s.Command, s.ID)

	return err
}

func (s *service) deleteService(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM services WHERE id=$1", s.ID)

	return err
}

func (s *service) createService(db *sql.DB) error {
	err := db.QueryRow(
		"INSERT INTO services(name, endpoint, command) VALUES($1, $2, $3) RETURNING id",
		s.Name, s.Endpoint, &s.Command).Scan(&s.ID)

	if err != nil {
		return err
	}

	return nil
}
