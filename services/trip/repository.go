package trip

import (
	"database/sql"

	"github.com/nerijusro/scootinAboot/types"
)

type TripsRepository struct {
	db *sql.DB
}

var publishEventQuery = "INSERT INTO events (trip_id, event_type, latitude, longitude, created_at, sequence) VALUES (UUID_TO_BIN(?, false), ?, ?, ?, ?, ?)"

func NewTripsRepository(db *sql.DB) *TripsRepository {
	return &TripsRepository{db: db}
}

func (r *TripsRepository) StartTrip(trip types.Trip, scooterOptLockVersion *int, userOptLockVersion *int, event types.TripEvent) error {
	createTripQuery := "INSERT INTO trips (id, user_id, scooter_id) VALUES (UUID_TO_BIN(?, false), UUID_TO_BIN(?, false), UUID_TO_BIN(?, false))"
	updateScooterQuery := "UPDATE scooters SET is_available = false, opt_lock_version = ? + 1 WHERE id = UUID_TO_BIN(?, false) AND opt_lock_version = ?"
	updateUserQuery := "UPDATE users SET is_eligible_to_travel = false, opt_lock_version = ? + 1 WHERE id = UUID_TO_BIN(?, false) AND opt_lock_version = ?"

	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec(createTripQuery, trip.ID.String(), trip.ClientId.String(), trip.ScooterId.String())
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec(updateScooterQuery, *scooterOptLockVersion, trip.ScooterId.String(), *scooterOptLockVersion)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec(updateUserQuery, *userOptLockVersion, trip.ClientId.String(), *userOptLockVersion)
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

func (r *TripsRepository) UpdateTrip(trip *types.Trip, scooterOptLockVersion *int, event types.TripEvent) error {
	updateScootersLocationQuery := "UPDATE scooters SET latitude = ?, longitude = ?, opt_lock_version = ? + 1 WHERE id = UUID_TO_BIN(?, false) AND opt_lock_version = ?"

	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec(updateScootersLocationQuery, event.Location.Latitude, event.Location.Longitude, *scooterOptLockVersion, trip.ScooterId.String(), *scooterOptLockVersion)
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

func (r *TripsRepository) EndTrip(trip *types.Trip, scooterOptLockVersion *int, userOptLockVersion *int, event types.TripEvent) error {
	updateTripQuery := "UPDATE trips SET is_finished = true WHERE id = UUID_TO_BIN(?, false)"
	updateScooterQuery := "UPDATE scooters SET is_available = true, opt_lock_version = ? + 1 WHERE id = UUID_TO_BIN(?, false) AND opt_lock_version = ?"
	updateUserQuery := "UPDATE users SET is_eligible_to_travel = true, opt_lock_version = ? + 1 WHERE id = UUID_TO_BIN(?, false) AND opt_lock_version = ?"

	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec(updateTripQuery, trip.ID.String())
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec(updateScooterQuery, *scooterOptLockVersion, trip.ScooterId.String(), *scooterOptLockVersion)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec(updateUserQuery, *userOptLockVersion, trip.ClientId.String(), *userOptLockVersion)
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

func (r *TripsRepository) GetTripById(id string) (*types.Trip, error) {
	row := r.db.QueryRow("SELECT * FROM trips WHERE id = UUID_TO_BIN(?, false)", id)

	var trip types.Trip
	err := row.Scan(&trip.ID, &trip.ClientId, &trip.ScooterId, &trip.IsFinished)
	if err != nil {
		return nil, err
	}

	return &trip, nil
}
