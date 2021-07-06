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