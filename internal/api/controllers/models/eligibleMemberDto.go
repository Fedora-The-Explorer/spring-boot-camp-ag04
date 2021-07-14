package models

// EligibleMemberDto is the eligible member dto model
type EligibleMemberDto struct {
	Skills  HeistSkillsDto         `json:"skills"`
	Members EligibleHeistMemberDto `json:"members"`
}
