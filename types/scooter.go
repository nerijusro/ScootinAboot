package types

type Scooter struct {
	ID         string   `json:"id"`
	Location   Location `json:"location"`
	Status     string   `json:"status"`
	OccupiedBy string   `json:"occupied_by"`
}
