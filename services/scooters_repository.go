package services

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

func (s *ScootersRepository) GetScooterById(id uuid.UUID) (*types.Scooter, error) {
	rows, err := s.db.Query("SELECT * FROM scooters WHERE id = ?", id)
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

func (s *ScootersRepository) CreateScooter(scooter types.Scooter) error {
	_, err := s.db.Exec("INSERT INTO scooters (id, location, is_available, occupied_by) VALUES (?, ?, ?, ?)",
		scooter.ID, scooter.Location, scooter.IsAvailable, scooter.OccupiedBy)
	if err != nil {
		return err
	}

	return nil
}

func scanRowIntoScooter(row *sql.Rows) (*types.Scooter, error) {
	scooter := new(types.Scooter)
	if err := row.Scan(&scooter.ID, &scooter.Location, &scooter.IsAvailable, &scooter.OccupiedBy); err != nil {
		return nil, err
	}

	return scooter, nil
}
