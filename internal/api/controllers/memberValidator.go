package controllers

import (
	"elProfessor/internal/api/controllers/models"
)

// MemberValidator validates member insert requests
type MemberValidator interface{
	MemberIsValid(memberDto models.MemberDto) bool
	MemberSkillsUpdateValidator(memberSkillUpdateDto models.MemberSkillsUpdateDto) bool
}
