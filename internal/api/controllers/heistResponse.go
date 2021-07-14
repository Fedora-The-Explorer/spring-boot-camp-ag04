package controllers

import (
	"context"
	"elProfessor/internal/api/controllers/models"
	"github.com/gin-gonic/gin"
)

// HeistResponse implements heist related functions
type HeistResponse interface {
	InsertHeist(heistDto models.HeistDto) (string, error)
	UpdateHeistSkills(ctx context.Context, heistSkills models.HeistSkillsDto, heistId string) error
	AddHeistMembers(members []string, id string) (string, error, []string)
	StartHeist(id string) (string, error)
	GetHeistById(ctx context.Context, id string) (models.HeistDto, bool, error)
	GetHeistMembersByHeistId(ctx context.Context, id string) ([]models.MemberDto, bool, error)
	GetHeistSkillsByHeistId(ctx *gin.Context, id string) (models.HeistSkillsDto, error)
	GetHeistStatusByHeistId(ctx *gin.Context, id string) (string, error)
	EndHeist(id string) (string, error)
	GetHeistOutcomeByHeistId(ctx *gin.Context, id string) (string, bool, error)
}
