package models

type MemberSkillsUpdateDto struct {
	Skills MemberSkillsDto `json:"skills"`
	MainSkill string `json:"mainSkill"`
}

