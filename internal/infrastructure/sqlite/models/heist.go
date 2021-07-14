package models

import "time"

// Heist is the heist storage model
type Heist struct {
	Id        string
	Name      string
	Location  string
	StartTime time.Time
	EndTime   time.Time
	Status    string
	Outcome   string
}
