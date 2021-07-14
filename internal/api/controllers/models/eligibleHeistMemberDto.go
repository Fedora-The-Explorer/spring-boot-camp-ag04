package models

// EligibleHeistMemberDto is the eligible heist member dto model
type EligibleHeistMemberDto []struct {
	Name   string          `json:"name"`
	Skills MemberSkillsDto `json:"skills"`
}
