package controllers

import (
	"context"
	"elProfessor/internal/api/controllers/models"
	"github.com/gin-gonic/gin"
)

// MemberResponse implements member related functions
type HeistResponse interface {
	InsertHeist(heistDto models.HeistDto) error
	UpdateHeistSkills(ctx context.Context, heistSkills models.HeistSkillsDto, heistId string) error
	AddHeistMembers(members []string, id string) (string,error)
	StartHeist(id string) (string,error)
	GetHeistById(ctx context.Context, id string) (models.HeistDto, bool, error)
	GetHeistMembersByHeistId(ctx *gin.Context, id string) (models.MemberDto, bool, error)
}