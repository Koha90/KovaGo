package domain

// Status represent the operational state of a shop.
type Status string

const (
	StatusActive   Status = "active"
	StatusInactive Status = "inactive"
)
