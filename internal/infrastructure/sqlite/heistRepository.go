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

// InsertMember inserts a storage model of members, their skills and unique skills into the database
func(r *HeistRepository) InsertMember(memberDto domainmodels.MemberDto) error {
	storageMember, storageSkills, storageMemberSkills := r.heistMapper.MapDomainMemberToStorageMember(memberDto)
	storageMemberSkills = r.CheckAndInsertSkills(storageSkills, storageMemberSkills)

	unique, err := r.CheckUniqueEmail(storageMember)
	if err != nil {
		return err
	}
	if !unique{
		return err
	}

	storageMember.MainSkillId, err = r.GetSkillIdByNameQuery(storageMember.MainSkillId)
	if err != nil {
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

// GetSkillIdByNameQuery gets the unique skill id from skills table based on the given name
func(r *HeistRepository) GetSkillIdByNameQuery(name string) (string,error) {
	var id string
	row, err := r.dbExecutor.QueryContext(context.Background(), "SELECT id FROM skills WHERE LOWER(name)=LOWER('"+name+"');")
	if err != nil {
		return id, err
	}
	defer row.Close()
	err = row.Scan(&id)
	if err != nil {
		return id, err
	}

	return id, nil
}

// CheckAndInsertSkills checks if the skill is unique, if yes it inserts it into the db, if not it does not and it passes the id of the existing skill to the memberSkill
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

// CheckUniqueEmail checks if the email of the member is unique
func(r *HeistRepository) CheckUniqueEmail(storageMember storagemodels.Member) (bool, error) {
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
		return false, err
	} else {
		return true, nil
	}
}

// UpdateMemberSkills
func(r *HeistRepository) UpdateMemberSkills(ctx context.Context, memberSkillsDto domainmodels.MemberSkillsUpdateDto, id string) error{
	storageMemberSkills, storageSkills, mainSkill := r.heistMapper.MapDomainSkillsToStorageSkills(memberSkillsDto, id)
	storageMemberSkills = r.CheckAndInsertSkills(storageSkills, storageMemberSkills)
	var err error
	if len(storageMemberSkills[0].Name) == 0 {
		mainSkill, err = r.GetSkillIdByNameQuery(mainSkill)
		if err != nil {
			return err
		}
		r.dbExecutor.Exec("UPDATE members SET main_skill='" + mainSkill + "'WHERE id='" + id + "';")
	} else if len(mainSkill) ==0 {
		r.UpsertSkills(storageMemberSkills, id)
	} else {
		mainSkill, err = r.GetSkillIdByNameQuery(mainSkill)
		if err != nil {
			return err
		}
		r.dbExecutor.Exec("UPDATE members SET main_skill='" + mainSkill + "'WHERE id='" + id + "';")
		r.UpsertSkills(storageMemberSkills, id)
	}


	return nil
}

func(r *HeistRepository) UpsertSkills(skills []storagemodels.MemberSkill, id string){
	for _, skill := range skills{
		r.dbExecutor.Exec("IF NOT EXISTS (SELECT * FROM memberSkills WHERE memberId ='" + id + "'AND skillId = '" + skill.SkillId + "') INSERT INTO memberSkills VALUES '" + skill.MemberId + "','" + skill.SkillId + "','" + skill.Name + "','" + skill.Level + "'ELSE UPDATE memberSkills SET name = '" + skill.Name + "', level = '" + skill.Level + "'WHERE memberId = '" + id + "'AND skillId = '" + skill.SkillId + "';")
	}
}

func(r *HeistRepository) DeleteMemberSkill(memberId, skillName string) error{
	skillId, err := r.GetSkillIdByNameQuery(skillName)
	if err != nil {
		return err
	}
	r.dbExecutor.Exec("DELETE FROM memberSkills WHERE memberId ='" + memberId + "'AND skillId='" + skillId + "';")
	return nil
}