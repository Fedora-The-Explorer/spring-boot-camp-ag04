package models

import "time"

type Heist struct {
	Id string
	Name string
	Location string
	StartTime time.Time
	EndTime time.Time
}