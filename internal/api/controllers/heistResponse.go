package controllers

import (
	"context"
	"elProfessor/internal/api/controllers/models"
)

// MemberResponse implements member related functions
type HeistResponse interface {
	InsertHeist(heistDto models.HeistDto) error
	UpdateHeistSkills(ctx context.Context, heistSkills models.HeistSkillsDto, heistId string) error
}