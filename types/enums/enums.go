package enums

type Availability string

const (
	Available   Availability = "available"
	Unavailable Availability = "unavailable"
	All         Availability = "all"
)

type TripEventType string

const (
	StartTrip  TripEventType = "start_trip_event"
	UpdateTrip TripEventType = "update_trip_event"
	EndTrip    TripEventType = "end_trip_event"
)
