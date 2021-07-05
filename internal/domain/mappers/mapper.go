package mappers

import (
	"github.com/nu7hatch/gouuid"
	"log"

	domainmodels "elProfessor/internal/api/controllers/models"
	storagemodels "elProfessor/internal/infrastructure/sqlite/models"
)

type HeistMapper struct{
}

func NewHeistMapper() *HeistMapper{
	return &HeistMapper{}
}

func(m *HeistMapper) MapDomainMemberToStorageMember(domainMember domainmodels.MemberDto) (storagemodels.Member, []storagemodels.Skill, []storagemodels.MemberSkill){
	memberId, err := uuid.NewV4()
	if err != nil {
		log.Fatalf("%s: %s", "failed to create uuid", err)
	}
	member := storagemodels.Member{
		Id: memberId.String(),
		Name: domainMember.Name,
		Sex:domainMember.Sex,
		Email:domainMember.Email,
		MainSkillId:domainMember.MainSkill,
		Status:domainMember.Status,
	}

	var skills []storagemodels.Skill
	var memberSkills []storagemodels.MemberSkill

	for _, skill := range domainMember.Skills {
		skillId, err := uuid.NewV4()
		if err != nil {
			log.Fatalf("%s: %s", "failed to create uuid", err)
		}
		currentSkill := storagemodels.Skill{
			Id: skillId.String(),
			Name: skill.Name,
		}


		currentMemberSkill := storagemodels.MemberSkill{
			MemberId: memberId.String(),
			SkillId: skillId.String(),
			Name: skill.Name,
			Level: skill.Level,
		}

		skills = append(skills, currentSkill)
		memberSkills = append(memberSkills, currentMemberSkill)

	}
	return member, skills, memberSkills
}

func(m *HeistMapper) MapDomainSkillsToStorageSkills(memberSkillsUpdateDto domainmodels.MemberSkillsUpdateDto, id string) ([]storagemodels.MemberSkill, []storagemodels.Skill, string) {
	var skills []storagemodels.Skill
	var memberSkills []storagemodels.MemberSkill
	for _, skill := range memberSkillsUpdateDto.Skills {
		skillId, err := uuid.NewV4()
		if err != nil {
			log.Fatalf("%s: %s", "failed to create uuid", err)
		}
		currentSkill := storagemodels.Skill{
			Id: skillId.String(),
			Name: skill.Name,
		}

		currentMemberSkill := storagemodels.MemberSkill{
			MemberId: id,
			SkillId: skillId.String(),
			Name: skill.Name,
			Level: skill.Level,
		}

		skills = append(skills, currentSkill)
		memberSkills = append(memberSkills, currentMemberSkill)
	}

	mainSkill := memberSkillsUpdateDto.MainSkill

	return memberSkills, skills, mainSkill
}