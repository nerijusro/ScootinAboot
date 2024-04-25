package trip

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/nerijusro/scootinAboot/types"
)

type TripRepository struct {
	db *sql.DB
}

var publishEventQuery = "INSERT INTO events (trip_id, event_type, latitude, longitude, created_at, sequence) VALUES (UUID_TO_BIN(?, false), ?, ?, ?, ?, ?)"
var updateScooterQuery = "UPDATE scooters SET is_available = ?, opt_lock_version = ? + 1 WHERE id = UUID_TO_BIN(?, false) AND opt_lock_version = ?"
var updateUserQuery = "UPDATE users SET is_eligible_to_travel = ?, opt_lock_version = ? + 1 WHERE id = UUID_TO_BIN(?, false) AND opt_lock_version = ?"

func NewRepository(db *sql.DB) *TripRepository {
	return &TripRepository{db: db}
}

func (r *TripRepository) StartTrip(trip types.Trip, scooterOptLockVersion *int, userOptLockVersion *int, event types.TripEvent) error {
	createTripQuery := "INSERT INTO trips (id, user_id, scooter_id) VALUES (UUID_TO_BIN(?, false), UUID_TO_BIN(?, false), UUID_TO_BIN(?, false))"

	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec(createTripQuery, trip.ID.String(), trip.ClientId.String(), trip.ScooterId.String())
	if err != nil {
		tx.Rollback()
		return err
	}

	err = r.updateAvailablity(tx, updateScooterQuery, trip.ScooterId.String(), false, scooterOptLockVersion)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = r.updateAvailablity(tx, updateUserQuery, trip.ClientId.String(), false, userOptLockVersion)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec(publishEventQuery, event.TripID.String(), event.Type, event.Location.Latitude, event.Location.Longitude, event.CreatedAt, event.Sequence)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (r *TripRepository) UpdateTrip(trip *types.Trip, scooterOptLockVersion *int, event types.TripEvent) error {
	updateScootersLocationQuery := "UPDATE scooters SET latitude = ?, longitude = ?, opt_lock_version = ? + 1 WHERE id = UUID_TO_BIN(?, false) AND opt_lock_version = ?"

	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	scooterUpdateResult, err := tx.Exec(updateScootersLocationQuery, event.Location.Latitude, event.Location.Longitude, *scooterOptLockVersion, trip.ScooterId.String(), *scooterOptLockVersion)
	if err != nil {
		tx.Rollback()
		return err
	}

	rowsAffecter, err := scooterUpdateResult.RowsAffected()
	if err != nil {
		tx.Rollback()
		return err
	}

	if rowsAffecter == 0 {
		tx.Rollback()
		return errors.New("scooter was updated by another transaction")
	}

	_, err = tx.Exec(publishEventQuery, event.TripID.String(), event.Type, event.Location.Latitude, event.Location.Longitude, event.CreatedAt, event.Sequence)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (r *TripRepository) EndTrip(trip *types.Trip, scooterOptLockVersion *int, userOptLockVersion *int, event types.TripEvent) error {
	updateTripQuery := "UPDATE trips SET is_finished = true WHERE id = UUID_TO_BIN(?, false)"

	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec(updateTripQuery, trip.ID.String())
	if err != nil {
		tx.Rollback()
		return err
	}

	err = r.updateAvailablity(tx, updateScooterQuery, trip.ScooterId.String(), true, scooterOptLockVersion)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = r.updateAvailablity(tx, updateUserQuery, trip.ClientId.String(), true, userOptLockVersion)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec(publishEventQuery, event.TripID.String(), event.Type, event.Location.Latitude, event.Location.Longitude, event.CreatedAt, event.Sequence)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (r *TripRepository) GetTripById(id string) (*types.Trip, error) {
	row := r.db.QueryRow("SELECT * FROM trips WHERE id = UUID_TO_BIN(?, false)", id)

	var trip types.Trip
	err := row.Scan(&trip.ID, &trip.ClientId, &trip.ScooterId, &trip.IsFinished)
	if err != nil {
		return nil, err
	}

	if trip.ID.String() != id {
		return nil, fmt.Errorf("trip with id %s not found", id)
	}

	return &trip, nil
}

func (r *TripRepository) updateAvailablity(tx *sql.Tx, query string, id string, newValue bool, optLockVersion *int) error {
	rowUpdateResult, err := tx.Exec(query, newValue, *optLockVersion, id, *optLockVersion)
	if err != nil {
		tx.Rollback()
		return err
	}

	rowsAffecter, err := rowUpdateResult.RowsAffected()
	if err != nil {
		tx.Rollback()
		return err
	}

	if rowsAffecter == 0 {
		tx.Rollback()
		return errors.New("row was updated by another transaction")
	}

	return nil
}
