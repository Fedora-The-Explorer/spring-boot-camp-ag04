package services

import (
	"context"
	domainmodels "elProfessor/internal/api/controllers/models"
	"elProfessor/internal/infrastructure/sqlite"
	"github.com/gin-gonic/gin"
)

type HeistResponse struct {
	heistRepository *sqlite.HeistRepository
}

func NewHeistResponse(	heistRepository *sqlite.HeistRepository) *HeistResponse{
	return &HeistResponse{
		heistRepository: heistRepository,
	}
}

func (h HeistResponse) InsertHeist(heistDto domainmodels.HeistDto) (string, error) {
	return h.heistRepository.InsertHeist(heistDto)
}

func (h HeistResponse) UpdateHeistSkills(ctx context.Context, heistSkills domainmodels.HeistSkillsDto, heistId string) error{
	return h.heistRepository.UpdateHeistSkills(ctx, heistSkills,heistId)
}

func(h HeistResponse) 	AddHeistMembers(members []string, id string) (string,error, []string){
	return h.heistRepository.AddHeistMembers(members, id)
}

func (h HeistResponse) 	StartHeist(id string) (string,error) {
	return h.heistRepository.StartHeist(id)
}

func (h HeistResponse) 	GetHeistById(ctx context.Context, id string) (domainmodels.HeistDto, bool, error){
	return h.heistRepository.GetHeistById(ctx, id)
}

func (h HeistResponse) GetHeistMembersByHeistId(ctx context.Context, id string) ([]domainmodels.MemberDto, bool, error) {
	return h.heistRepository.GetHeistMembersByHeistId(ctx, id)
}

func (h HeistResponse) GetHeistSkillsByHeistId(ctx *gin.Context, id string) (domainmodels.HeistSkillsDto, error) {
	return h.heistRepository.GetHeistSkillsByHeistId(ctx, id)
}

func (h HeistResponse) 	GetHeistStatusByHeistId(ctx *gin.Context, id string) (string, error) {
	return h.heistRepository.GetHeistStatusByHeistId(ctx, id)
}

func (h HeistResponse) 	EndHeist(id string) (string, error){
	return h.heistRepository.EndHeist(id)
}

func (h HeistResponse) 	GetHeistOutcomeByHeistId(ctx *gin.Context, id string) (string, bool, error){
	return h.heistRepository.GetHeistOutcomeByHeistId(ctx, id)
}


