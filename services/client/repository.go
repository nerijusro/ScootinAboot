package client

import (
	"database/sql"

	"github.com/nerijusro/scootinAboot/types"
)

type ClientsRepository struct {
	db *sql.DB
}

func NewClientsRepository(db *sql.DB) *ClientsRepository {
	return &ClientsRepository{db: db}
}

func (r *ClientsRepository) CreateUser(client types.MobileClient) error {
	_, err := r.db.Exec("INSERT INTO users (id, full_name, is_eligible_to_travel) VALUES (UUID_TO_BIN(?, false), ?)",
		client.ID.String(), client.FullName, true)
	if err != nil {
		return err
	}

	return nil
}
