package models

// HeistSkill is the heist skill storage model
type HeistSkill struct {
	SkillId string
	HeistId string
	Level   string
	Members int
}
