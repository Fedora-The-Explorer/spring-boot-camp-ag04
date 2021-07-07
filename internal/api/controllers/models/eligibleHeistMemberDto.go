package models

// HeistMemberDto is the get request dto model
type EligibleHeistMemberDto []struct {
	Name string `json:"name"`
	Skills MemberSkillsDto `json:"skills"`
}