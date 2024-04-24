package scooter

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/nerijusro/scootinAboot/types"
	"github.com/nerijusro/scootinAboot/types/enums"
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
	availabilityFilter := enums.Availability(queryParams.Availability)
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
		scooter, err := scanRowIntoScooter(rows, nil)
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
		scooter, err := scanRowIntoScooter(rows, nil)
		if err != nil {
			return nil, err
		}

		scooters = append(scooters, scooter)
	}

	return scooters, nil
}

func (r *ScootersRepository) GetScooterById(id string) (*types.Scooter, *int, error) {
	rows, err := r.db.Query("SELECT * FROM scooters WHERE id = UUID_TO_BIN(?, false)", id)
	if err != nil {
		return nil, nil, err
	}

	scooter := new(types.Scooter)
	optLockVersion := new(int)
	for rows.Next() {
		scooter, err = scanRowIntoScooter(rows, optLockVersion)
		if err != nil {
			return nil, nil, err
		}
	}

	if scooter.ID == uuid.Nil {
		return nil, nil, fmt.Errorf("scooter with id %s not found", id)
	}

	return scooter, optLockVersion, nil
}

func scanRowIntoScooter(row *sql.Rows, optLockVersion *int) (*types.Scooter, error) {
	location := new(types.Location)
	scooter := new(types.Scooter)

	if err := row.Scan(&scooter.ID, &location.Latitude, &location.Longitude, &scooter.IsAvailable, &optLockVersion); err != nil {
		return nil, err
	}

	scooter.Location = *location
	return scooter, nil
}

func tryAddingAvailabilityFilter(query string, availability enums.Availability) string {
	if availability == enums.Available {
		return query + " AND is_available = true"
	}
	if availability == enums.Unavailable {
		return query + " AND is_available = false"
	}

	return query
}
