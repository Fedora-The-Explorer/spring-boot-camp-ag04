package services

import (
	"context"
	domainmodels "elProfessor/internal/api/controllers/models"
)

type HeistHandler interface {
	InsertHeist(heistDto domainmodels.HeistDto) error
	UpdateHeistSkills(ctx context.Context, heistSkills domainmodels.HeistSkillsDto, heistId string) error
	AddHeistMembers(members []string, id string) (string,error)

}