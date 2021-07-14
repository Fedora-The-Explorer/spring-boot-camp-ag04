package services

import (
	"context"
	domainmodels "elProfessor/internal/api/controllers/models"
	"github.com/gin-gonic/gin"
)

type HeistResponse struct {
	heistHandler HeistHandler
}

func NewHeistResponse(heistHandler HeistHandler) *HeistResponse{
	return &HeistResponse{
		heistHandler: heistHandler,
	}
}

func (h HeistResponse) InsertHeist(heistDto domainmodels.HeistDto) (string, error) {
	return h.heistHandler.InsertHeist(heistDto)
}

func (h HeistResponse) UpdateHeistSkills(ctx context.Context, heistSkills domainmodels.HeistSkillsDto, heistId string) error{
	return h.heistHandler.UpdateHeistSkills(ctx, heistSkills,heistId)
}

func(h HeistResponse) 	AddHeistMembers(members []string, id string) (string,error, []string){
	return h.heistHandler.AddHeistMembers(members, id)
}

func (h HeistResponse) 	StartHeist(id string) (string,error) {
	return h.heistHandler.StartHeist(id)
}

func (h HeistResponse) 	GetHeistById(ctx context.Context, id string) (domainmodels.HeistDto, bool, error){
	return h.heistHandler.GetHeistById(ctx, id)
}

func (h HeistResponse) GetHeistMembersByHeistId(ctx context.Context, id string) ([]domainmodels.MemberDto, bool, error) {
	return h.heistHandler.GetHeistMembersByHeistId(ctx, id)
}

func (h HeistResponse) GetHeistSkillsByHeistId(ctx *gin.Context, id string) (domainmodels.HeistSkillsDto, error) {
	return h.heistHandler.GetHeistSkillsByHeistId(ctx, id)
}

func (h HeistResponse) 	GetHeistStatusByHeistId(ctx *gin.Context, id string) (string, error) {
	return h.heistHandler.GetHeistStatusByHeistId(ctx, id)
}

func (h HeistResponse) 	EndHeist(id string) (string, error){
	return h.heistHandler.EndHeist(id)
}

func (h HeistResponse) 	GetHeistOutcomeByHeistId(ctx *gin.Context, id string) (string, bool, error){
	return h.heistHandler.GetHeistOutcomeByHeistId(ctx, id)
}


