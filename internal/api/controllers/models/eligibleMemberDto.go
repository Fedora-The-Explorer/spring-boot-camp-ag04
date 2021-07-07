package models

type EligibleMemberDto struct {
	Skills HeistSkillsDto `json:"skills"`
	Members HeistMemberDto `json:"members"`
}

