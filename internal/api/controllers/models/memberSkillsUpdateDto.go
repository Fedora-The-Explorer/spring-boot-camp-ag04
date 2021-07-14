package models

// MemberSkillsUpdateDto is the update member skills dto model
type MemberSkillsUpdateDto struct {
	Skills    MemberSkillsDto `json:"skills"`
	MainSkill string          `json:"mainSkill"`
}
