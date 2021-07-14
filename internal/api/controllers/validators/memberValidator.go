package validators

import "elProfessor/internal/api/controllers/models"

// MemberValidator validates insert member requests
type MemberValidator struct {
}

// NewMemberValidator creates a new instance of MemberValidator
func NewMemberValidator() *MemberValidator {
	return &MemberValidator{

	}
}

// MemberIsValid checks if event update is valid
// Sex is F or M
// Status is one of the following values: AVAILABLE, EXPIRED, INCARCERATED, RETIRED
func (v *MemberValidator) MemberIsValid(memberDto models.MemberDto) bool {
	if memberDto.Sex != "F" && memberDto.Sex != "M" {
		return false
	} else if memberDto.Status != "AVAILABLE" && memberDto.Status != "EXPIRED" && memberDto.Status != "INCARCERATED" && memberDto.Status != "RETIRED" {
		return false
	} else {
		return true
	}
}

func (v *MemberValidator) MemberSkillsUpdateValidator(memberSkillsUpdate models.MemberSkillsUpdateDto) bool {
	if len(memberSkillsUpdate.MainSkill) == 0 && len(memberSkillsUpdate.Skills[0].Name) == 0 {
		return false
	} else {
		return true
	}
}
