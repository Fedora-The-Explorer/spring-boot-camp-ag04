package models

type HeistDto struct {
	Name string `json:"name"`
	Location string `json:"location"`
	StartTime string `json:"startTime"`
	EndTime string `json:"endTime"`
	Skills HeistSkillsDto `json:"skills"`
}