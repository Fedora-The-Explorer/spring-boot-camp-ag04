package models

type HeistSkillsDto []struct {
	Name string `json:"name"`
	Level string `json:"level"`
	Members int `json:"members"`
}

