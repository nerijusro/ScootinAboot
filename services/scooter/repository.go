package scooter

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/nerijusro/scootinAboot/types"
)

type ScootersRepository struct {
	db *sql.DB
}

func NewScootersRepository(db *sql.DB) *ScootersRepository {
	return &ScootersRepository{db: db}
}

func (s *ScootersRepository) CreateScooter(scooter types.Scooter) error {
	_, err := s.db.Exec("INSERT INTO scooters (id, longitude, latitude, is_available) VALUES (UUID_TO_BIN(?, false), ?, ?, ?)",
		scooter.ID.String(), scooter.Location.Longitude, scooter.Location.Latitude, scooter.IsAvailable)
	if err != nil {
		return err
	}

	return nil
}

func (s *ScootersRepository) GetScooters() ([]*types.Scooter, error) {
	rows, err := s.db.Query("SELECT * FROM scooters")
	if err != nil {
		return nil, err
	}

	scooters := make([]*types.Scooter, 0)
	for rows.Next() {
		scooter, err := scanRowIntoScooter(rows)
		if err != nil {
			return nil, err
		}

		scooters = append(scooters, scooter)
	}

	return scooters, nil
}

func (s *ScootersRepository) GetScooterById(id string) (*types.Scooter, error) {
	rows, err := s.db.Query("SELECT * FROM scooters WHERE id = UUID_TO_BIN(?, false)", id)
	if err != nil {
		return nil, err
	}

	scooter := new(types.Scooter)
	for rows.Next() {
		scooter, err = scanRowIntoScooter(rows)
		if err != nil {
			return nil, err
		}
	}

	if scooter.ID == uuid.Nil {
		return nil, fmt.Errorf("scooter with id %s not found", id)
	}

	return scooter, nil
}

func scanRowIntoScooter(row *sql.Rows) (*types.Scooter, error) {
	optlockVersion := new(int)
	location := new(types.Location)
	scooter := new(types.Scooter)

	if err := row.Scan(&scooter.ID, &location.Latitude, &location.Longitude, &scooter.IsAvailable, &scooter.OccupiedBy, &optlockVersion); err != nil {
		return nil, err
	}

	scooter.Location = *location
	return scooter, nil
}
