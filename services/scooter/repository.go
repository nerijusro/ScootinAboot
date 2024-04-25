package scooter

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/nerijusro/scootinAboot/types"
	"github.com/nerijusro/scootinAboot/types/enums"
)

type ScooterRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *ScooterRepository {
	return &ScooterRepository{db: db}
}

func (r *ScooterRepository) CreateScooter(scooter types.Scooter) error {
	_, err := r.db.Exec("INSERT INTO scooters (id, longitude, latitude, is_available) VALUES (UUID_TO_BIN(?, false), ?, ?, ?)",
		scooter.ID.String(), scooter.Location.Longitude, scooter.Location.Latitude, scooter.IsAvailable)
	if err != nil {
		return err
	}

	return nil
}

func (r *ScooterRepository) GetScootersByArea(queryParams types.GetScootersQueryParameters) ([]*types.Scooter, error) {
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
		scooter, _, err := scanRowIntoScooter(rows)
		if err != nil {
			return nil, err
		}

		scooters = append(scooters, scooter)
	}

	return scooters, nil
}

func (r *ScooterRepository) GetAllScooters() ([]*types.Scooter, error) {
	rows, err := r.db.Query("SELECT * FROM scooters")
	if err != nil {
		return nil, err
	}

	scooters := make([]*types.Scooter, 0)
	for rows.Next() {
		scooter, _, err := scanRowIntoScooter(rows)
		if err != nil {
			return nil, err
		}

		scooters = append(scooters, scooter)
	}

	return scooters, nil
}

func (r *ScooterRepository) GetScooterById(id string) (*types.Scooter, *int, error) {
	rows, err := r.db.Query("SELECT * FROM scooters WHERE id = UUID_TO_BIN(?, false)", id)
	if err != nil {
		return nil, nil, err
	}

	var scooter *types.Scooter
	var optLockVersion *int
	for rows.Next() {
		scooter, optLockVersion, err = scanRowIntoScooter(rows)
		if err != nil {
			return nil, nil, err
		}
	}

	if scooter.ID == uuid.Nil {
		return nil, nil, fmt.Errorf("scooter with id %s not found", id)
	}

	return scooter, optLockVersion, nil
}

func scanRowIntoScooter(row *sql.Rows) (*types.Scooter, *int, error) {
	var location types.Location
	var scooter types.Scooter
	var optLockVersion int

	if err := row.Scan(&scooter.ID, &location.Latitude, &location.Longitude, &scooter.IsAvailable, &optLockVersion); err != nil {
		return nil, nil, err
	}

	scooter.Location = location
	return &scooter, &optLockVersion, nil
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
