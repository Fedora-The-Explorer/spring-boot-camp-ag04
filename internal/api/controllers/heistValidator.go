package controllers

import "elProfessor/internal/api/controllers/models"


// HeistValidator validates heist insert requests
type HeistValidator interface{
	HeistSkillUpdateValidator(heistSkills models.HeistSkillsDto) bool
}
