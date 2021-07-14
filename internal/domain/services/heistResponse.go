package services

import (
	"context"
	domainmodels "elProfessor/internal/api/controllers/models"
	"elProfessor/internal/infrastructure/sqlite"
	"github.com/gin-gonic/gin"
)

// HeistResponse implements the heist related functions form the repository.
type HeistResponse struct {
	heistRepository *sqlite.HeistRepository
}

// NewHeistResponse creates a new instance of HeistResponse.
func NewHeistResponse(heistRepository *sqlite.HeistRepository) *HeistResponse {
	return &HeistResponse{
		heistRepository: heistRepository,
	}
}

// InsertHeist implements the InsertHeist function from the heist repository
func (h HeistResponse) InsertHeist(heistDto domainmodels.HeistDto) (string, error) {
	return h.heistRepository.InsertHeist(heistDto)
}

// UpdateHeistSkills implements the UpdateHeistSkills function from the heist repository
func (h HeistResponse) UpdateHeistSkills(ctx context.Context, heistSkills domainmodels.HeistSkillsDto, heistId string) error {
	return h.heistRepository.UpdateHeistSkills(ctx, heistSkills, heistId)
}

// AddHeistMembers implements the AddHeistMembers function from the heist repository
func (h HeistResponse) AddHeistMembers(members []string, id string) (string, error, []string) {
	return h.heistRepository.AddHeistMembers(members, id)
}

// StartHeist implements the StartHeist function from the heist repository
func (h HeistResponse) StartHeist(id string) (string, error) {
	return h.heistRepository.StartHeist(id)
}

// GetHeistById implements the GetHeistById function from the heist repository
func (h HeistResponse) GetHeistById(ctx context.Context, id string) (domainmodels.HeistDto, bool, error) {
	return h.heistRepository.GetHeistById(ctx, id)
}

// GetHeistMembersByHeistId implements the GetHeistMembersByHeistId function from the heist repository
func (h HeistResponse) GetHeistMembersByHeistId(ctx context.Context, id string) ([]domainmodels.MemberDto, bool, error) {
	return h.heistRepository.GetHeistMembersByHeistId(ctx, id)
}

// GetHeistSkillsByHeistId implements the GetHeistSkillsByHeistId function from the heist repository
func (h HeistResponse) GetHeistSkillsByHeistId(ctx *gin.Context, id string) (domainmodels.HeistSkillsDto, error) {
	return h.heistRepository.GetHeistSkillsByHeistId(ctx, id)
}

// GetHeistStatusByHeistId implements the GetHeistStatusByHeistId function from the heist repository
func (h HeistResponse) GetHeistStatusByHeistId(ctx *gin.Context, id string) (string, error) {
	return h.heistRepository.GetHeistStatusByHeistId(ctx, id)
}

// EndHeist implements the EndHeist function from the heist repository
func (h HeistResponse) EndHeist(id string) (string, error) {
	return h.heistRepository.EndHeist(id)
}

// GetHeistOutcomeByHeistId implements the GetHeistOutcomeByHeistId function from the heist repository
func (h HeistResponse) GetHeistOutcomeByHeistId(ctx *gin.Context, id string) (string, bool, error) {
	return h.heistRepository.GetHeistOutcomeByHeistId(ctx, id)
}
