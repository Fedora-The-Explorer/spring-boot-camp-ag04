package sqlite

import (
	"context"

	domainmodels "elProfessor/internal/api/controllers/models"
	storagemodels "elProfessor/internal/infrastructure/sqlite/models"
)

type HeistRepository struct {
	dbExecutor DatabaseExecutor
	heistMapper HeistMapper
}

func NewHeistRepository(dbExecutor DatabaseExecutor, heistMapper HeistMapper) * HeistRepository{
	return &HeistRepository{
		dbExecutor: dbExecutor,
		heistMapper: heistMapper,
	}
}

func(r *HeistRepository) InsertMember(memberDto domainmodels.MemberDto) error {
	storageMember, storageSkills, storageMemberSkills := r.heistMapper.MapDomainMemberToStorageMember(memberDto)
	storageMemberSkills = r.CheckAndInsertSkills(storageSkills, storageMemberSkills)

	row, err := r.dbExecutor.QueryContext(context.Background(), "SELECT id FROM skills WHERE LOWER(name)=LOWER('"+storageMember.MainSkillId+"');")
	if err != nil {
		return err
	}
	defer row.Close()

	var id string
	err = row.Scan(&id)
	if err != nil {
		return err
	}
	storageMember.MainSkillId = id

	if !r.CheckUniqueEmail(storageMember){
		return err
	}

	r.dbExecutor.Exec("INSERT INTO skills VALUES ('" + storageMember.Id + "', '"+ storageMember.Name + "', '"+ storageMember.Sex + "', '"+ storageMember.Email + "', '"+ storageMember.MainSkillId + "', '"+ storageMember.Status + "');")

	for _, skill := range storageMemberSkills{
		if len(skill.Level) == 0{	// makes the default value of a skill level *
			skill.Level = "*"
		}
		r.dbExecutor.Exec("INSERT INTO skills VALUES ('" + skill.MemberId + "', '"+ skill.SkillId + "', '"+ skill.Name + "', '"+ skill.Level + "');")
	}
	return nil
}


func(r *HeistRepository) CheckAndInsertSkills(storageSkills []storagemodels.Skill, storageMemberSkills []storagemodels.MemberSkill) ([]storagemodels.MemberSkill){
	for idx, skill := range storageSkills {
		row, err := r.dbExecutor.QueryContext(context.Background(), "SELECT id FROM skills WHERE LOWER(name)=LOWER('"+skill.Name+"');")
		if err != nil {
			panic(err)
		}

		var id string
		err = row.Scan(&id)
		if err != nil {
			panic(err)
		}

		if len(id) > 0 {
			storageMemberSkills[idx].SkillId = id
		} else {
			r.dbExecutor.Exec("INSERT INTO skills VALUES ('" + skill.Id + "', '"+ skill.Name + "');")
		}
		row.Close()
	}
	 return storageMemberSkills
}

func(r *HeistRepository) CheckUniqueEmail(storageMember storagemodels.Member) bool {
	row, err := r.dbExecutor.QueryContext(context.Background(), "SELECT email FROM members WHERE email='"+storageMember.Email+"';")
	if err != nil {
		panic(err)
	}
	defer row.Close()

	var email string
	err = row.Scan(&email)
	if err != nil {
		panic(err)
	}

	if email == storageMember.Email {
		return false
	} else {
		return true
	}
}

func(r *HeistRepository) UpdateMemberSkills(ctx context.Context, memberSkillsDto domainmodels.MemberSkillsUpdateDto, id string) error{
	storageMemberSkills, storageSkills, mainSkill := r.heistMapper.MapDomainSkillsToStorageSkills(memberSkillsDto, id)
	storageMemberSkills = r.CheckAndInsertSkills(storageSkills, storageMemberSkills)

	// TODO if memberSkills is empty fill only main skill if main skill is empty fill memberSkills
	// TODO or fill both. query by user id and use MERGE function in sql to upsert the data


	return nil
}
