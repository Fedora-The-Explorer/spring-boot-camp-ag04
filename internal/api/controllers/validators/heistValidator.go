package validators

import "elProfessor/internal/api/controllers/models"

// HeistValidator validates insert heist requests
type HeistValidator struct {
}
// NewHeistValidator creates a new instance of HeistValidator
func NewHeistValidator () *MemberValidator{
	return &MemberValidator{
	}
}


func (v *MemberValidator) HeistSkillUpdateValidator(heistSkills models.HeistSkillsDto) bool {
	for i:= 0; i<len(heistSkills); i++{
		for k:= 0; k<len(heistSkills); k++{
			if heistSkills[i].Name == heistSkills[k].Name && heistSkills[i].Level == heistSkills[k].Level && i!=k{
				return false
			}
		}
	}
	return true
}