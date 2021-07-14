package models

import "time"

// HeistDto is the heist dto model
type HeistDto struct {
	Name      string         `json:"name"`
	Location  string         `json:"location"`
	StartTime time.Time      `json:"startTime"`
	EndTime   time.Time      `json:"endTime"`
	Skills    HeistSkillsDto `json:"skills"`
	Status    string         `json:"status"`
}
