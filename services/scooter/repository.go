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

func (r *ScootersRepository) CreateScooter(scooter types.Scooter) error {
	_, err := r.db.Exec("INSERT INTO scooters (id, longitude, latitude, is_available) VALUES (UUID_TO_BIN(?, false), ?, ?, ?)",
		scooter.ID.String(), scooter.Location.Longitude, scooter.Location.Latitude, scooter.IsAvailable)
	if err != nil {
		return err
	}

	return nil
}

func (r *ScootersRepository) GetScootersByArea(queryParams types.GetScootersQueryParameters) ([]*types.Scooter, error) {
	availabilityFilter := types.Availability(queryParams.Availability)
	getScootersByAreaQuery := tryAddingAvailabilityFilter(
		"SELECT * FROM scooters WHERE latitude >= ? AND latitude <= ? AND longitude >= ? AND longitude <= ?",
		availabilityFilter)

	rows, err := r.db.Query(getScootersByAreaQuery,
		queryParams.Y1, queryParams.Y2, queryParams.X1, queryParams.X2)
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

func (r *ScootersRepository) GetAllScooters() ([]*types.Scooter, error) {
	rows, err := r.db.Query("SELECT * FROM scooters")
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

func (r *ScootersRepository) GetScooterById(id string) (*types.Scooter, error) {
	rows, err := r.db.Query("SELECT * FROM scooters WHERE id = UUID_TO_BIN(?, false)", id)
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

func tryAddingAvailabilityFilter(query string, availability types.Availability) string {
	if availability == types.Available {
		return query + " AND is_available = true"
	}
	if availability == types.Unavailable {
		return query + " AND is_available = false"
	}

	return query
}
