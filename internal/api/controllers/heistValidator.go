package controllers

import "elProfessor/internal/api/controllers/models"


// MemberValidator validates member insert requests
type HeistValidator interface{
	HeistSkillUpdateValidator(heistSkills models.HeistSkillsDto) bool
}
