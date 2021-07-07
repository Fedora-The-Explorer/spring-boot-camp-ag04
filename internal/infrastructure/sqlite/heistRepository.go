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
	storageMemberSkills = r.CheckAndInsertMemberSkills(storageSkills, storageMemberSkills)

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


	r.dbExecutor.Exec("INSERT INTO members VALUES ('" + storageMember.Id + "', '"+ storageMember.Name + "', '"+ storageMember.Sex + "', '"+ storageMember.Email + "', '"+ storageMember.MainSkillId + "', '"+ storageMember.Status + "');")

	for _, skill := range storageMemberSkills{
		if len(skill.Level) == 0{	// makes the default value of a skill level *
			skill.Level = "*"
		}else if len(skill.Level) > 10{
			skill.Level = "**********"
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

// CheckAndInsertMemberSkills checks if the skill is unique, if yes it inserts it into the db, if not it does not and it passes the id of the existing skill to the memberSkill
func(r *HeistRepository) CheckAndInsertMemberSkills(storageSkills []storagemodels.Skill, storageMemberSkills []storagemodels.MemberSkill) ([]storagemodels.MemberSkill){
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

// UpdateMemberSkills updates or adds new records of member's skills and a main skill
func(r *HeistRepository) UpdateMemberSkills(ctx context.Context, memberSkillsDto domainmodels.MemberSkillsUpdateDto, id string) error{
	storageMemberSkills, storageSkills, mainSkill := r.heistMapper.MapDomainSkillsToStorageSkills(memberSkillsDto, id)
	storageMemberSkills = r.CheckAndInsertMemberSkills(storageSkills, storageMemberSkills)
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

// UpsertSkills adds or updates member's skills
func(r *HeistRepository) UpsertSkills(skills []storagemodels.MemberSkill, id string){
	for _, skill := range skills{
		r.dbExecutor.Exec("IF NOT EXISTS (SELECT * FROM memberSkills WHERE memberId ='" + id + "'AND skillId = '" + skill.SkillId + "') INSERT INTO memberSkills VALUES '" + skill.MemberId + "','" + skill.SkillId + "','" + skill.Name + "','" + skill.Level + "'ELSE UPDATE memberSkills SET name = '" + skill.Name + "', level = '" + skill.Level + "'WHERE memberId = '" + id + "'AND skillId = '" + skill.SkillId + "';")
	}
}

// DeleteMemberSkill deletes the selected member's skill
func(r *HeistRepository) DeleteMemberSkill(memberId, skillName string) error{
	skillId, err := r.GetSkillIdByNameQuery(skillName)
	if err != nil {
		return err
	}
	r.dbExecutor.Exec("DELETE FROM memberSkills WHERE memberId ='" + memberId + "'AND skillId='" + skillId + "';")
	return nil
}

// InsertHeist inserts a storage model of heists, their skills and unique skills into the database
func(r *HeistRepository) InsertHeist(heistDto domainmodels.HeistDto) error {
	storageHeist, storageSkills, storageHeistSkills := r.heistMapper.MapDomainHeistToStorageHeist(heistDto)
	storageHeistSkills = r.CheckAndInsertHeistSkills(storageSkills,storageHeistSkills)

	unique, err := r.CheckUniqueName(storageHeist)
	if err != nil {
		return err
	}
	if !unique {
		return err
	}
	// if getting errors from the query check the go sdk sql drivers
	defaultStatus := "PLANNING"
	r.dbExecutor.Exec("INSERT INTO heists VALUES ('" + storageHeist.Id + "', '"+ storageHeist.Name + "', '"+ storageHeist.Location + "', '"+ storageHeist.StartTime + "', '"+ storageHeist.EndTime + "', '" + defaultStatus + "');")

	for _, skill := range storageHeistSkills{
		if len(skill.Level) == 0{	// makes the default value of a skill level *
			skill.Level = "*"
		}else if len(skill.Level) > 10{
			skill.Level = "**********"
		}
		r.dbExecutor.Exec("INSERT INTO heistSkills VALUES ('" + skill.SkillId + "', '"+ skill.HeistId + "', '"+ skill.Level + "','"+ skill.Members + "');")

	}
	return nil
}

// CheckAndInsertHeistSkills checks if the skill is unique, if yes it inserts it into the db, if not it does not and it passes the id of the existing skill to the heistSkill
func(r *HeistRepository) CheckAndInsertHeistSkills(storageSkills []storagemodels.Skill, storageHeistSkills []storagemodels.HeistSkill) ([]storagemodels.HeistSkill){
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
			storageHeistSkills[idx].SkillId = id
		} else {
			r.dbExecutor.Exec("INSERT INTO skills VALUES ('" + skill.Id + "', '"+ skill.Name + "');")
		}
		row.Close()
	}
	return storageHeistSkills
}

// CheckUniqueName checks if the heist has a unique name
func (r *HeistRepository) CheckUniqueName(heist storagemodels.Heist) (bool, error) {
	row, err := r.dbExecutor.QueryContext(context.Background(), "SELECT name FROM heists WHERE name='"+heist.Name+"';")
	if err != nil {
		panic(err)
	}
	defer row.Close()

	var name string
	err = row.Scan(&name)
	if err != nil {
		panic(err)
	}

	if name == heist.Name {
		return false, err
	} else {
		return true, nil
	}
}

func(r *HeistRepository) UpdateHeistSkills(ctx context.Context, skills domainmodels.HeistSkillsDto, id string) error{
	storageHeistSkills, storageSkills := r.heistMapper.MapDomainHeistSkillsToStorageHeistSkills(skills, id)
	storageHeistSkills = r.CheckAndInsertHeistSkills(storageSkills, storageHeistSkills)

	// if getting errors from the query check the go sdk sql drivers
	for _, skill := range storageHeistSkills{
		if len(skill.Level) == 0{	// makes the default value of a skill level *
			skill.Level = "*"
		}else if len(skill.Level) > 10{
			skill.Level = "**********"
		}
		r.dbExecutor.Exec("INSERT INTO skills VALUES ('" + skill.SkillId + "', '"+ skill.HeistId + "', '"+ skill.Level + "','"+ skill.Members + "');")

	}

	return nil
}

func(r *HeistRepository) GetEligibleMembers(ctx context.Context, id string) (domainmodels.EligibleMemberDto, bool, error){
	var eligible domainmodels.EligibleMemberDto
	heistSkills, err := r.queryGetSkillsByHeistIdEligible(ctx, id)
	if err != nil {
		return domainmodels.EligibleMemberDto{}, false, err
	}
	var memberSkills []storagemodels.MemberSkill
 	for _,skill:= range heistSkills{
		memberSkills = append(memberSkills, r.queryGetSkillsByIdEligible(ctx, skill.SkillId, skill.Level))
	}

	var members []storagemodels.Member
 	for _,skill := range memberSkills{
 		members = append(members, r.queryGetMemberByIdEligible(ctx, skill.MemberId))
	}

	for idx,skill:= range heistSkills{
		eligible.Skills[idx].Name, err = r.queryGetSkillNameById(ctx, skill.SkillId)
		if err != nil {
			return domainmodels.EligibleMemberDto{}, false, err
		}
		eligible.Skills[idx].Level =  skill.Level
		eligible.Skills[idx].Members = skill.Members
	}

	for idx, member := range members {
		eligible.Members[idx].Name = member.Name
		idy := 0
		for _, skill := range memberSkills {
			if member.Id == skill.MemberId {
				eligible.Members[idx].Skills[idy].Name = skill.Name
				eligible.Members[idx].Skills[idy].Name = skill.Level
				idy++
			}
		}
	}
	return eligible, true, nil
}

func (r *HeistRepository) queryGetSkillsByHeistIdEligible(ctx context.Context, id string) ([]storagemodels.HeistSkill, error) {
	row, err := r.dbExecutor.QueryContext(ctx, "SELECT * FROM heistSkills WHERE heistId='"+id+"';")
	if err != nil {
		return []storagemodels.HeistSkill{}, err
	}

	defer row. Close()
	var skills []storagemodels.HeistSkill

	for row.Next(){
		var skillId string
		var heistId string
		var level string
		var members int

		err = row.Scan(&skillId, &heistId, &level, &members)
		if err != nil {
			return []storagemodels.HeistSkill{}, err
		}
		skills = append(skills, storagemodels.HeistSkill{
			SkillId: skillId,
			HeistId: heistId,
			Level: level,
			Members: members,
		})

	}
	return skills, nil
}

func (r *HeistRepository) queryGetSkillsByIdEligible(ctx context.Context, id string, skillLevel string) storagemodels.MemberSkill {
	row, err := r.dbExecutor.QueryContext(ctx, "SELECT * FROM memberSkills WHERE (skillId='"+id+ "'AND level >='" + skillLevel + "');")
	if err != nil {
		panic(err)
	}
	defer row. Close()

	var memberId string
	var skillId string
	var name string
	var level string

	err = row.Scan(&memberId, &skillId, &name, &level)
	if err != nil {
		panic(err)
	}

	return storagemodels.MemberSkill{
		MemberId: memberId,
		SkillId: skillId,
		Name: name,
		Level: level,
	}
}

func (r *HeistRepository) queryGetMemberByIdEligible(ctx context.Context, id string) storagemodels.Member {
	status1, status2 := "AVAILABLE", "RETIRED"
	row, err := r.dbExecutor.QueryContext(ctx, "SELECT * FROM members WHERE id='"+id+"' AND (status='" + status1 + "'OR status='" + status2 + "');")
	if err != nil {
		panic(err)
	}
	defer row. Close()

	var memberId string
	var name string
	var sex string
	var email string
	var mainSkill string
	var status string

	err = row.Scan(&memberId, &name, &sex, &email, &mainSkill, &status)
	if err != nil {
		panic(err)
	}

	return storagemodels.Member{
		Id: memberId,
		Name: name,
		Sex: sex,
		Email: email,
		MainSkillId: mainSkill,
		Status: status,
	}
}

func (r *HeistRepository) queryGetSkillNameById(ctx context.Context, id string) (string, error) {
	row, err := r.dbExecutor.QueryContext(ctx, "SELECT * FROM members WHERE id='"+id+"';")
	if err != nil {
		return "", err
	}
	defer row.Close()

	var name string
	err = row.Scan(&name)
	if err != nil {
		return "", err
	}
	return name, nil
}

func (r *HeistRepository) AddHeistMembers(members []string, id string) (string, error){
	var memberIds []string
	code, err := r.checkHeist(id, "PLANNING")
	if err != nil {
		return code, err
	}

	for _, member := range members{
		memberIds = append(memberIds, r.queryGetMemberIdByName(member))
	}

	for _, idx := range memberIds{
		r.dbExecutor.Exec("INSERT INTO heistMembers VALUES ('" + idx + "', '"+ id + "');")
	}
	ready := "READY"
	r.dbExecutor.Exec("UPDATE heists SET status='" + ready + "'WHERE id='" + id + "';")

	return code, nil
}

func (r *HeistRepository) queryGetMemberIdByName(member string) string {
	row, err := r.dbExecutor.QueryContext(context.Background(), "SELECT id FROM members WHERE name='"+member+"';")
	if err != nil {
		panic(err)
	}
	var id string
	err = row.Scan(&id)
	if err != nil {
		panic(err)
	}
	defer row.Close()
	return id
}

func (r *HeistRepository) checkHeist(id string, heistStatus string) (string, error) {
	row, err := r.dbExecutor.QueryContext(context.Background(), "SELECT status FROM heists WHERE id='"+id+"';")
	if err != nil {
		return "404", err
	}
	defer row.Close()

	var status string
	err = row.Scan(&status)
	if err != nil {
		return "404", err
	}

	if status != heistStatus {
		return "405", err
	}

	return "", nil
}

func(r *HeistRepository) StartHeist(id string) (string, error){
	code, err := r.checkHeist(id, "READY")
	if err != nil {
		return code, err
	}
	inProgress := "IN_PROGRESS"
	r.dbExecutor.Exec("UPDATE heists SET status='" + inProgress + "'WHERE id='" + id + "';")
	return "", nil
}
