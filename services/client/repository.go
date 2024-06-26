package client

import (
	"database/sql"

	"github.com/nerijusro/scootinAboot/types"
)

type ClientRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *ClientRepository {
	return &ClientRepository{db: db}
}

func (r *ClientRepository) CreateUser(client types.MobileClient) error {
	_, err := r.db.Exec("INSERT INTO users (id, full_name) VALUES (UUID_TO_BIN(?, false), ?)",
		client.ID.String(), client.FullName)
	if err != nil {
		return err
	}

	return nil
}

func (r *ClientRepository) GetUserById(id string) (*types.MobileClient, *int, error) {
	row := r.db.QueryRow("SELECT * FROM users WHERE id = UUID_TO_BIN(?, false)", id)

	var client types.MobileClient
	var optLockVersion int
	err := row.Scan(&client.ID, &client.FullName, &client.IsEligibleToTravel, &optLockVersion)
	if err != nil {
		return nil, nil, err
	}

	return &client, &optLockVersion, nil
}
