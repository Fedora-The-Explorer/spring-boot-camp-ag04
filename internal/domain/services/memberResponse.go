package services

import (
	"context"
	domainmodels "elProfessor/internal/api/controllers/models"
	"elProfessor/internal/infrastructure/sqlite"
)

// MemberResponse implements the member related functions form the repository.
type MemberResponse struct {
	heistRepository *sqlite.HeistRepository
}

// NewMemberResponse creates a new instance of MemberResponse.
func NewMemberResponse(heistRepository *sqlite.HeistRepository) *MemberResponse {
	return &MemberResponse{
		heistRepository: heistRepository,
	}
}

// InsertMember implements the InsertMember function from the heist repository
func (m MemberResponse) InsertMember(memberDto domainmodels.MemberDto) error {
	return m.heistRepository.InsertMember(memberDto)
}

// UpdateMemberSkills implements the UpdateMemberSkills function from the heist repository
func (m MemberResponse) UpdateMemberSkills(ctx context.Context, memberSkillsUpdate domainmodels.MemberSkillsUpdateDto, memberId string) error {
	return m.heistRepository.UpdateMemberSkills(ctx, memberSkillsUpdate, memberId)
}

// DeleteMemberSkill implements the DeleteMemberSkill function from the heist repository
func (m MemberResponse) DeleteMemberSkill(memberId string, skillName string) error {
	return m.heistRepository.DeleteMemberSkill(memberId, skillName)
}

// GetEligibleMembers implements the GetEligibleMembers function from the heist repository
func (m MemberResponse) GetEligibleMembers(ctx context.Context, id string) (domainmodels.EligibleMemberDto, bool, error) {
	return m.heistRepository.GetEligibleMembers(ctx, id)
}

// GetMemberById implements the GetMemberById function from the heist repository
func (m MemberResponse) GetMemberById(ctx context.Context, id string) (domainmodels.MemberDto, bool, error) {
	return m.heistRepository.GetMemberByID(ctx, id)
}

// GetMemberSkillsById implements the GetMemberSkillsById function from the heist repository
func (m MemberResponse) GetMemberSkillsById(ctx context.Context, id string) (domainmodels.MemberSkillsDto, bool, error) {
	return m.heistRepository.GetMemberSkillsById(ctx, id)
}
