package controllers

import (
	"context"
	"elProfessor/internal/api/controllers/models"
)

// MemberResponse implements member related functions
type MemberResponse interface {
	InsertMember(memberDto models.MemberDto) error
	UpdateMemberSkills(ctx context.Context, memberSkillsUpdate models.MemberSkillsUpdateDto, memberId string) error
	DeleteMemberSkill(memberId string, skillName string) error
	GetEligibleMembers(ctx context.Context, id string) (models.EligibleMemberDto, bool, error)
	GetMemberById(ctx context.Context, id string) (models.MemberDto, bool, error)
	GetMemberSkillsById(ctx context.Context, id string) (models.MemberSkillsDto, bool, error)

}