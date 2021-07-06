package sqlite

import (
	domainmodels "elProfessor/internal/api/controllers/models"
	storagemodels "elProfessor/internal/infrastructure/sqlite/models"
)

type HeistMapper interface {
	MapDomainMemberToStorageMember(domainMember domainmodels.MemberDto) (storagemodels.Member, []storagemodels.Skill, []storagemodels.MemberSkill)
	MapDomainSkillsToStorageSkills(memberSkillsUpdate domainmodels.MemberSkillsUpdateDto, id string) ([]storagemodels.MemberSkill, []storagemodels.Skill, string)
	MapDomainHeistToStorageHeist(heistDto domainmodels.HeistDto) (storagemodels.Heist,[]storagemodels.Skill,[]storagemodels.HeistSkill)
	MapDomainHeistSkillsToStorageHeistSkills(heistSkills domainmodels.HeistSkillsDto, id string) ([]storagemodels.HeistSkill, []storagemodels.Skill)
}
