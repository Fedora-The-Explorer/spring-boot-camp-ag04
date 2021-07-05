package models

import "time"

type Heist struct {
	Id string
	Location string
	StartTime time.Time
	EndTime time.Time
	Status string
}