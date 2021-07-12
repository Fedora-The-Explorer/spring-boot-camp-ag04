package services

import (
	"context"
	domainmodels "elProfessor/internal/api/controllers/models"
)

type MemberResponse struct {
	memberHandler MemberHandler
}

func NewMemberResponse(memberHandler MemberHandler) *MemberResponse{
	return &MemberResponse{
		memberHandler: memberHandler,
	}
}

func (m MemberResponse) InsertMember(memberDto domainmodels.MemberDto) error {
	return m.InsertMember(memberDto)
}

func (m MemberResponse) UpdateMemberSkills(ctx context.Context,memberSkillsUpdate domainmodels.MemberSkillsUpdateDto, memberId string) error {
	return m.UpdateMemberSkills(ctx, memberSkillsUpdate, memberId)
}

func (m MemberResponse) DeleteMemberSkill(memberId string, skillName string) error {
	return m.DeleteMemberSkill(memberId,skillName)
}

func (m MemberResponse) GetEligibleMembers(ctx context.Context, id string) (domainmodels.EligibleMemberDto, bool, error){
	return m.GetEligibleMembers(ctx, id)
}

func (m MemberResponse) GetMemberById(ctx context.Context, id string) (domainmodels.MemberDto, bool, error){
	return m.GetMemberById(ctx,id)
}

func (m MemberResponse) GetMemberSkillsById(ctx context.Context, id string) (domainmodels.MemberSkillsDto, bool, error){
	return m.GetMemberSkillsById(ctx, id)
}

