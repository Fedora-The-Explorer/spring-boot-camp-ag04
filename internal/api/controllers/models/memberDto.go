package models

// MemberDto is the member insert request dto model
type MemberDto struct {
	Name      string          `json:"name"`
	Sex       string          `json:"sex"`
	Email     string          `json:"email"`
	Skills    MemberSkillsDto `json:"skills"`
	MainSkill string          `json:"mainSkill"`
	Status    string          `json:"status"`
}
