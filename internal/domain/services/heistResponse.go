package services

import (
	"context"
	domainmodels "elProfessor/internal/api/controllers/models"
)

type HeistResponse struct {
	heistHandler HeistHandler
}

func NewHeistResponse(heistHandler HeistHandler) *HeistResponse{
	return &HeistResponse{
		heistHandler: heistHandler,
	}
}

func (h HeistResponse) InsertHeist(heistDto domainmodels.HeistDto) error {
	return h.InsertHeist(heistDto)
}

func (h HeistResponse) UpdateHeistSkills(ctx context.Context, heistSkills domainmodels.HeistSkillsDto, heistId string) error{
	return h.UpdateHeistSkills(ctx, heistSkills,heistId)
}

func(h HeistResponse) 	AddHeistMembers(members []string, id string) (string,error){
	return h.AddHeistMembers(members, id)
}

func (h HeistResponse) 	StartHeist(id string) (string,error) {
	return h.StartHeist(id)
}

func (h HeistResponse) 	GetHeistById(ctx context.Context, id string) (domainmodels.HeistDto, bool, error){
	return h.GetHeistById(ctx, id)
}
