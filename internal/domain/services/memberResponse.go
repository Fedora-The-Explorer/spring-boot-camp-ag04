package services

import (
	"context"
	domainmodels "elProfessor/internal/api/controllers/models"
	"elProfessor/internal/infrastructure/sqlite"
)

type MemberResponse struct {
	heistRepository *sqlite.HeistRepository
}

func NewMemberResponse(heistRepository *sqlite.HeistRepository) *MemberResponse{
	return &MemberResponse{
		heistRepository: heistRepository,
	}
}

func (m MemberResponse) InsertMember(memberDto domainmodels.MemberDto) error {
	return m.heistRepository.InsertMember(memberDto)
}

func (m MemberResponse) UpdateMemberSkills(ctx context.Context,memberSkillsUpdate domainmodels.MemberSkillsUpdateDto, memberId string) error {
	return m.heistRepository.UpdateMemberSkills(ctx, memberSkillsUpdate, memberId)
}

func (m MemberResponse) DeleteMemberSkill(memberId string, skillName string) error {
	return m.heistRepository.DeleteMemberSkill(memberId,skillName)
}

func (m MemberResponse) GetEligibleMembers(ctx context.Context, id string) (domainmodels.EligibleMemberDto, bool, error){
	return m.heistRepository.GetEligibleMembers(ctx, id)
}

func (m MemberResponse) GetMemberById(ctx context.Context, id string) (domainmodels.MemberDto, bool, error){
	return m.heistRepository.GetMemberByID(ctx,id)
}

func (m MemberResponse) GetMemberSkillsById(ctx context.Context, id string) (domainmodels.MemberSkillsDto, bool, error){
	return m.heistRepository.GetMemberSkillsById(ctx, id)
}

