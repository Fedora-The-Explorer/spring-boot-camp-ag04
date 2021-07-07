package models

type EligibleMemberDto struct {
	Skills HeistSkillsDto `json:"skills"`
	Members EligibleHeistMemberDto `json:"members"`
}

