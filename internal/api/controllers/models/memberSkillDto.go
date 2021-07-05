package models

// MemberSkillsDto is the insert request dto model
type MemberSkillsDto []struct {
	Name string `json:"name"`
	Level string `json:"level"`
}
