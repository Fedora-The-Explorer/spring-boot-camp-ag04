package models

// HeistSkillsDto is the heist skills dto model
type HeistSkillsDto []struct {
	Name    string `json:"name"`
	Level   string `json:"level"`
	Members int    `json:"members"`
}
