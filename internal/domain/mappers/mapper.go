package mappers

import (
	"fmt"
	"github.com/nu7hatch/gouuid"
	"log"
	"time"

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

func(m *HeistMapper) MapDomainHeistToStorageHeist(heistDto domainmodels.HeistDto) (storagemodels.Heist,[]storagemodels.Skill,[]storagemodels.HeistSkill){
	heistId, err := uuid.NewV4()
	if err != nil {
		log.Fatalf("%s: %s", "failed to create uuid", err)
	}


	layout := "2006-01-02T15:04:05.000Z"
	startTime, err := time.Parse(layout, heistDto.StartTime)
	if err != nil {
		fmt.Println(err)
	}
	endTime, err := time.Parse(layout, heistDto.EndTime)
	if err != nil {
		fmt.Println(err)
	}
	
	heist := storagemodels.Heist{
		Id:        heistId.String(),
		Name:      heistDto.Name,
		Location:  heistDto.Location,
		StartTime: startTime,
		EndTime:   endTime,
	}

	var skills []storagemodels.Skill
	var heistSkills []storagemodels.HeistSkill

	for _, skill := range heistDto.Skills {
		skillId, err := uuid.NewV4()
		if err != nil {
			log.Fatalf("%s: %s", "failed to create uuid", err)
		}
		currentSkill := storagemodels.Skill{
			Id: skillId.String(),
			Name: skill.Name,
		}


		currentMemberSkill := storagemodels.HeistSkill{
			SkillId: skillId.String(),
			HeistId: heistId.String(),
			Level: skill.Level,
			Members: skill.Members,
		}

		skills = append(skills, currentSkill)
		heistSkills = append(heistSkills, currentMemberSkill)

	}

	return heist, skills, heistSkills
	
}