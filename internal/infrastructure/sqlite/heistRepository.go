package sqlite

import (
	"context"
	"database/sql"
	"github.com/gin-gonic/gin"
	"math/rand"
	"strconv"
	"time"

	domainmodels "elProfessor/internal/api/controllers/models"
	storagemodels "elProfessor/internal/infrastructure/sqlite/models"
)

// HeistRepository provides methods that operate on members, heists, skills, etc. on bets SQLite database.
type HeistRepository struct {
	dbExecutor  DatabaseExecutor
	heistMapper HeistMapper
}

// NewHeistRepository creates and returns a new HeistRepository
func NewHeistRepository(dbExecutor DatabaseExecutor, heistMapper HeistMapper) *HeistRepository {
	return &HeistRepository{
		dbExecutor:  dbExecutor,
		heistMapper: heistMapper,
	}
}

// InsertMember inserts a storage model of members, their skills and unique skills into the database. If a
// problem occurs during the data insertion or if the member is not unique, an error will be returned.
func (r *HeistRepository) InsertMember(memberDto domainmodels.MemberDto) error {
	storageMember, storageSkills, storageMemberSkills := r.heistMapper.MapDomainMemberToStorageMember(memberDto)
	storageMemberSkills = r.CheckAndInsertMemberSkills(storageSkills, storageMemberSkills)

	unique, err := r.CheckUniqueEmail(storageMember)
	if err != nil {
		return err
	}
	if !unique {
		return err
	}

	storageMember.MainSkillId, err = r.GetSkillIdByNameQuery(storageMember.MainSkillId)
	if err != nil {
		return err
	}

	r.dbExecutor.Exec("INSERT INTO members VALUES ('" + storageMember.Id + "', '" + storageMember.Name + "', '" + storageMember.Sex + "', '" + storageMember.Email + "', '" + storageMember.MainSkillId + "', '" + storageMember.Status + "');")

	for _, skill := range storageMemberSkills {
		if len(skill.Level) == 0 { // makes the default value of a skill level *
			skill.Level = "*"
		} else if len(skill.Level) > 10 {
			skill.Level = "**********"
		}
		r.dbExecutor.Exec("INSERT INTO skills VALUES ('" + skill.MemberId + "', '" + skill.SkillId + "', '" + skill.Name + "', '" + skill.Level + "');")
	}
	return nil
}

// GetSkillIdByNameQuery gets the unique skill id from skills table based on the given name.
func (r *HeistRepository) GetSkillIdByNameQuery(name string) (string, error) {
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

// CheckAndInsertMemberSkills checks if the skill is unique, if yes it inserts it into the db, if not it does not
// and it passes the id of the existing skill to the memberSkill.
func (r *HeistRepository) CheckAndInsertMemberSkills(storageSkills []storagemodels.Skill, storageMemberSkills []storagemodels.MemberSkill) []storagemodels.MemberSkill {
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
			r.dbExecutor.Exec("INSERT INTO skills VALUES ('" + skill.Id + "', '" + skill.Name + "');")
		}
		row.Close()
	}
	return storageMemberSkills
}

// CheckUniqueEmail checks if the email of the member is unique, if not false and and an error will be returned.
func (r *HeistRepository) CheckUniqueEmail(storageMember storagemodels.Member) (bool, error) {
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

// UpdateMemberSkills updates or adds new records of member's skills and a main skill, also checks if the new
//skills exist in the unique skill table, if not it adds them.
func (r *HeistRepository) UpdateMemberSkills(ctx context.Context, memberSkillsDto domainmodels.MemberSkillsUpdateDto, id string) error {
	storageMemberSkills, storageSkills, mainSkill := r.heistMapper.MapDomainSkillsToStorageSkills(memberSkillsDto, id)
	storageMemberSkills = r.CheckAndInsertMemberSkills(storageSkills, storageMemberSkills)
	var err error
	if len(storageMemberSkills[0].Name) == 0 {
		mainSkill, err = r.GetSkillIdByNameQuery(mainSkill)
		if err != nil {
			return err
		}
		r.dbExecutor.Exec("UPDATE members SET main_skill='" + mainSkill + "'WHERE id='" + id + "';")
	} else if len(mainSkill) == 0 {
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
func (r *HeistRepository) UpsertSkills(skills []storagemodels.MemberSkill, id string) {
	for _, skill := range skills {
		r.dbExecutor.Exec("IF NOT EXISTS (SELECT * FROM memberSkills WHERE memberId ='" + id + "'AND skillId = '" + skill.SkillId + "') INSERT INTO memberSkills VALUES '" + skill.MemberId + "','" + skill.SkillId + "','" + skill.Name + "','" + skill.Level + "'ELSE UPDATE memberSkills SET name = '" + skill.Name + "', level = '" + skill.Level + "'WHERE memberId = '" + id + "'AND skillId = '" + skill.SkillId + "';")
	}
}

// DeleteMemberSkill deletes the selected member's skill
func (r *HeistRepository) DeleteMemberSkill(memberId, skillName string) error {
	skillId, err := r.GetSkillIdByNameQuery(skillName)
	if err != nil {
		return err
	}
	r.dbExecutor.Exec("DELETE FROM memberSkills WHERE memberId ='" + memberId + "'AND skillId='" + skillId + "';")
	return nil
}

// InsertHeist inserts a storage model of heists, their skills and unique skills into the database, also checks if
//some skills are new and adds them to the unique skill table.
func (r *HeistRepository) InsertHeist(heistDto domainmodels.HeistDto) (string, error) {
	storageHeist, storageSkills, storageHeistSkills := r.heistMapper.MapDomainHeistToStorageHeist(heistDto)
	storageHeistSkills = r.CheckAndInsertHeistSkills(storageSkills, storageHeistSkills)

	unique, err := r.CheckUniqueName(storageHeist)
	if err != nil {
		return "", err
	}
	if !unique {
		return "", err
	}
	// if getting errors from the query check the go sdk sql drivers
	defaultStatus := "PLANNING"
	r.dbExecutor.Exec("INSERT INTO heists VALUES ('" + storageHeist.Id + "', '" + storageHeist.Name + "', '" + storageHeist.Location + "', '" + storageHeist.StartTime.String() + "', '" + storageHeist.EndTime.String() + "', '" + defaultStatus + "');")

	for _, skill := range storageHeistSkills {
		if len(skill.Level) == 0 { // makes the default value of a skill level *
			skill.Level = "*"
		} else if len(skill.Level) > 10 {
			skill.Level = "**********"
		}
		r.dbExecutor.Exec("INSERT INTO heistSkills VALUES ('" + skill.SkillId + "', '" + skill.HeistId + "', '" + skill.Level + "','" + strconv.Itoa(skill.Members) + "');")

	}
	return storageHeist.Id, err
}

// CheckAndInsertHeistSkills checks if the skill is unique, if yes it inserts it into the db, if not it does not and
// it passes the id of the existing skill to the heistSkill.
func (r *HeistRepository) CheckAndInsertHeistSkills(storageSkills []storagemodels.Skill, storageHeistSkills []storagemodels.HeistSkill) []storagemodels.HeistSkill {
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
			r.dbExecutor.Exec("INSERT INTO skills VALUES ('" + skill.Id + "', '" + skill.Name + "');")
		}
		row.Close()
	}
	return storageHeistSkills
}

// CheckUniqueName checks if the heist has a unique name. If it doesn't it returns an error.
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

// UpdateHeistSkills updates the skills of a certain heist.
func (r *HeistRepository) UpdateHeistSkills(ctx context.Context, skills domainmodels.HeistSkillsDto, id string) error {
	storageHeistSkills, storageSkills := r.heistMapper.MapDomainHeistSkillsToStorageHeistSkills(skills, id)
	storageHeistSkills = r.CheckAndInsertHeistSkills(storageSkills, storageHeistSkills)

	// if getting errors from the query check the go sdk sql drivers
	for _, skill := range storageHeistSkills {
		if len(skill.Level) == 0 { // makes the default value of a skill level *
			skill.Level = "*"
		} else if len(skill.Level) > 10 {
			skill.Level = "**********"
		}
		r.dbExecutor.Exec("INSERT INTO skills VALUES ('" + skill.SkillId + "', '" + skill.HeistId + "', '" + skill.Level + "','" + strconv.Itoa(skill.Members) + "');")

	}

	return nil
}

// GetEligibleMembers fetches members from the database depending if they satisfy the heist skill conditions.
//The second value indicates whether the heist exists in DB. If it does not exist an error will be returned.
func (r *HeistRepository) GetEligibleMembers(ctx context.Context, id string) (domainmodels.EligibleMemberDto, bool, error) {
	var eligible domainmodels.EligibleMemberDto
	heistSkills, err := r.queryGetSkillsByHeistIdEligible(ctx, id)
	if err != nil {
		return domainmodels.EligibleMemberDto{}, false, err
	}
	var memberSkills []storagemodels.MemberSkill
	for _, skill := range heistSkills {
		memberSkills = append(memberSkills, r.queryGetSkillsByIdEligible(ctx, skill.SkillId, skill.Level))
	}

	var members []storagemodels.Member
	for _, skill := range memberSkills {
		members = append(members, r.queryGetMemberByIdEligible(ctx, skill.MemberId))
	}

	for idx, skill := range heistSkills {
		eligible.Skills[idx].Name, err = r.queryGetSkillNameById(ctx, skill.SkillId)
		if err != nil {
			return domainmodels.EligibleMemberDto{}, false, err
		}
		eligible.Skills[idx].Level = skill.Level
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

// queryGetSkillsByHeistIdEligible gets the heist skills for the eligible comparison.
func (r *HeistRepository) queryGetSkillsByHeistIdEligible(ctx context.Context, id string) ([]storagemodels.HeistSkill, error) {
	row, err := r.dbExecutor.QueryContext(ctx, "SELECT * FROM heistSkills WHERE heistId='"+id+"';")
	if err != nil {
		return []storagemodels.HeistSkill{}, err
	}

	defer row.Close()
	var skills []storagemodels.HeistSkill

	for row.Next() {
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
			Level:   level,
			Members: members,
		})

	}
	return skills, nil
}

// queryGetSkillsByIdEligible gets the member skills for the eligibility comparison.
func (r *HeistRepository) queryGetSkillsByIdEligible(ctx context.Context, id string, skillLevel string) storagemodels.MemberSkill {
	row, err := r.dbExecutor.QueryContext(ctx, "SELECT * FROM memberSkills WHERE (skillId='"+id+"'AND level >='"+skillLevel+"');")
	if err != nil {
		panic(err)
	}
	defer row.Close()

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
		SkillId:  skillId,
		Name:     name,
		Level:    level,
	}
}

// queryGetMemberByIdEligible gets the eligible members ids
func (r *HeistRepository) queryGetMemberByIdEligible(ctx context.Context, id string) storagemodels.Member {
	status1, status2 := "AVAILABLE", "RETIRED"
	row, err := r.dbExecutor.QueryContext(ctx, "SELECT * FROM members WHERE id='"+id+"' AND (status='"+status1+"'OR status='"+status2+"');")
	if err != nil {
		panic(err)
	}
	defer row.Close()

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
	possibleMember := storagemodels.Member{
		Id:          memberId,
		Name:        name,
		Sex:         sex,
		Email:       email,
		MainSkillId: mainSkill,
		Status:      status,
	}
	err = r.checkPossibleHeistMember(possibleMember)
	if err != nil {
		panic(err)
	}
	return possibleMember
}

//queryGetSkillNameById returns the skill name based on the given id
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

// AddHeistMembers adds members to heist, returns the error code, error and a slice of emails used to send mails
//to members that were added to a heist.
func (r *HeistRepository) AddHeistMembers(members []string, id string) (string, error, []string) {
	var memberIds []string
	code, err := r.checkHeist(id, "PLANNING")
	if err != nil {
		return code, err, nil
	}

	for _, member := range members {
		memberIds = append(memberIds, r.queryGetMemberIdByName(member))
	}

	for _, idx := range memberIds {
		r.dbExecutor.Exec("INSERT INTO heistMembers VALUES ('" + idx + "', '" + id + "');")
	}
	ready := "READY"
	r.dbExecutor.Exec("UPDATE heists SET status='" + ready + "'WHERE id='" + id + "';")

	emails := r.queryGetMailsFromMemberIds(memberIds)

	return code, nil, emails
}

// queryGetMailsFromMemberIds returns emails - slice of strings used to send emails to members.
func (r *HeistRepository) queryGetMailsFromMemberIds(ids []string) []string {
	var row *sql.Rows
	var err error
	for _, id := range ids {
		row, err = r.dbExecutor.QueryContext(context.Background(), "SELECT email FROM members WHERE id='"+id+"';")
		if err != nil {
			panic(err)
		}
	}
	var emails []string
	for row.Next() {
		var email string
		err = row.Scan(&email)
		if err != nil {
			panic(err)
		}
		emails = append(emails, email)
	}

	return emails
}

// queryGetMemberIdByName returns the member id based on the given member name
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

// checkHeist checks the eligibility of a heist based on the given heist status, returns errors and error codes.
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

// StartHeist starts a heist that is eligible to have it's status updated to IN PROGRESS
func (r *HeistRepository) StartHeist(id string) (string, error) {
	code, err := r.checkHeist(id, "READY")
	if err != nil {
		return code, err
	}
	inProgress := "IN_PROGRESS"
	r.dbExecutor.Exec("UPDATE heists SET status='" + inProgress + "'WHERE id='" + id + "';")
	return "", nil
}

// EndHeist ends a heist that is eligible to have it's status updated to FINISHED.
// After it finishes, it updates the heist's status and applies repercussions.
// Depending on how many members attended the heist outcome can be FAILED or SUCCEEDED.
// Depending on how many members attended and the outcome members status will be updated to EXPIRED or
// INCARCERATED, the only case where no repercussions will be applied is when all required members attended the heist.
// It also improves the member's skills.
func (r *HeistRepository) EndHeist(id string) (string, error) {
	code, err := r.checkHeist(id, "IN_PROGRESS")
	if err != nil {
		return code, nil
	}
	finished := "FINISHED"
	r.dbExecutor.Exec("UPDATE heists SET status='" + finished + "'WHERE id='" + id + "';")
	heist, err := r.queryGetHeistById(context.Background(), id)
	if err != nil {
		return "", err
	}

	members, err := r.queryGetMemberIdByHeistId(context.Background(), id)
	if err != nil {
		return "", err
	}
	row, err := r.dbExecutor.QueryContext(context.Background(), "SELECT members FROM heists WHERE id='"+id+"';")
	if err != nil {
		return "", err
	}
	var membersRequired int
	err = row.Scan(&membersRequired)
	membersUsed := len(members)
	percentage := float64(membersUsed/membersRequired) * 100
	var amountOfMembersToUpdate int
	var outcome string
	var factor float64
	incarcerateOnly := false
	if percentage < 50 {
		outcome = "FAILED"
		amountOfMembersToUpdate = membersUsed
	} else if percentage < 75 {
		if rand.Int()%2 == 0 {
			outcome = "FAILED"
		} else {
			outcome = "SUCCEEDED"
		}
		if outcome == "SUCCEEDED" {
			factor = 0.33
			amountOfMembersToUpdate = int(float64(membersUsed) * factor)
		} else {
			factor = 0.66
			amountOfMembersToUpdate = int(float64(membersUsed) * factor)
		}
	} else if percentage < 100 {
		outcome = "SUCCEEDED"
		factor = 0.33
		amountOfMembersToUpdate = int(float64(membersUsed) * factor)
		incarcerateOnly = true
	} else {
		outcome = "SUCCEEDED"
	}
	incarcerate := "INCARCERATED"
	expire := "EXPIRED"
	for i := 0; i < amountOfMembersToUpdate; i++ {
		if incarcerateOnly {
			r.dbExecutor.Exec("UPDATE members SET status='" + incarcerate + "'WHERE id='" + members[i] + "';")
		} else {
			if rand.Int()%2 == 0 {
				r.dbExecutor.Exec("UPDATE members SET status='" + incarcerate + "'WHERE id='" + members[i] + "';")
			} else {
				r.dbExecutor.Exec("UPDATE members SET status='" + expire + "'WHERE id='" + members[i] + "';")
			}
		}
	}

	r.skillImprovement(members, heist, id)
	r.dbExecutor.Exec("UPDATE heists SET outcome='" + outcome + "'WHERE id='" + id + "';")
	return "", nil
}

// skillImprovement will add more levels to member's skills that were used during a heist depending on the time that
// it took to finish the heist. For every 24h spent on a heist the skill will get one level up tp ten levels.
func (r *HeistRepository) skillImprovement(memberIds []string, heist storagemodels.Heist, id string) {
	row, err := r.dbExecutor.QueryContext(context.Background(), "SELECT skillId FROM heistSkills WHERE heistId='"+id+"';")
	if err != nil {
		panic(err)
	}
	var skillIds []string
	for row.Next() {
		var skillId string
		err = row.Scan(&skillId)
		skillIds = append(skillIds, skillId)
	}

	overlappingSKills := r.queryGetOverlapHeistMemberSkills(skillIds, memberIds)
	levelUpTime := 86400
	timeDiff := heist.EndTime.Sub(heist.StartTime)
	seconds := timeDiff.Seconds()
	toLevelUp := int(seconds) / levelUpTime
	for _, skill := range overlappingSKills {
		levelUp := toLevelUp
		currentLevel := skill.Level
		for len(currentLevel) < 10 && levelUp > 0 {
			currentLevel += "*"
			levelUp--
		}
		r.dbExecutor.Exec("UPDATE memberSkills SET level='" + currentLevel + "'WHERE skillId='" + skill.SkillId + "', memberId='" + skill.MemberId + "';")
	}
}

// queryGetOverlapHeistMemberSkills gets the skills that were used for a heist and that the members have.
func (r *HeistRepository) queryGetOverlapHeistMemberSkills(skillIds []string, memberIds []string) []storagemodels.MemberSkill {
	var row *sql.Rows
	var err error

	for _, skill := range skillIds {
		for _, member := range memberIds {
			row, err = r.dbExecutor.QueryContext(context.Background(), "SELECT * FROM memberSkills WHERE skillId='"+skill+"', memberId='"+member+"';")
		}
	}

	var allSkills []storagemodels.MemberSkill
	for row.Next() {
		var memberId string
		var skillId string
		var name string
		var level string

		err = row.Scan(&memberId, &skillId, &name, &level)
		if err != nil {
			panic(err)
		}
		allSkills = append(allSkills, storagemodels.MemberSkill{
			MemberId: memberId,
			SkillId:  skillId,
			Name:     name,
			Level:    level,
		})
	}

	return allSkills
}

// checkPossibleHeistMember checks the possible eligible member.
func (r *HeistRepository) checkPossibleHeistMember(member storagemodels.Member) error {
	row, err := r.dbExecutor.QueryContext(context.Background(), "SELECT memberId FROM heistMembers;")
	for row.Next() {
		var memberId string
		err = row.Scan(&memberId)
		if err != nil {
			return err
		}
		if memberId == member.Id {
			return err
		}

	}
	return nil
}

// GetMemberByID fetches a member from the database and returns it. The second returned value indicates whether the
// member exists in DB. If the member does not exist, an error will not be returned.
func (r *HeistRepository) GetMemberByID(ctx context.Context, id string) (domainmodels.MemberDto, bool, error) {
	storageMember, err := r.queryGetMemberByID(ctx, id)
	if err != nil {
		return domainmodels.MemberDto{}, false, err
	}
	storageSkills, err := r.queryGetMemberSkillsByID(ctx, id)
	if err != nil {
		return domainmodels.MemberDto{}, false, err
	}
	mainSkillName, err := r.queryGetSkillNameById(ctx, storageMember.MainSkillId)
	if err != nil {
		return domainmodels.MemberDto{}, false, err
	}
	domainMember := domainmodels.MemberDto{
		Name:      storageMember.Name,
		Sex:       storageMember.Sex,
		Email:     storageMember.Email,
		MainSkill: mainSkillName,
		Status:    storageMember.Status,
	}

	for idx, skills := range storageSkills {
		domainMember.Skills[idx].Name = skills.Name
		domainMember.Skills[idx].Level = skills.Level
	}
	return domainMember, true, nil
}

// queryGetMemberByID dbExecutor call to query the selected member by id
func (r *HeistRepository) queryGetMemberByID(ctx context.Context, id string) (storagemodels.Member, error) {
	row, err := r.dbExecutor.QueryContext(ctx, "SELECT * FROM members WHERE id='"+id+"';")
	if err != nil {
		return storagemodels.Member{}, err
	}
	defer row.Close()

	var idx string
	var name string
	var sex string
	var email string
	var mainSkillId string
	var status string

	err = row.Scan(&idx, &name, &sex, &email, &mainSkillId, &status)
	if err != nil {
		return storagemodels.Member{}, err
	}

	return storagemodels.Member{
		Id:          idx,
		Name:        name,
		Sex:         sex,
		Email:       email,
		MainSkillId: mainSkillId,
		Status:      status,
	}, nil
}

// queryGetMemberSkillsByID dbExecutor call to query the selected member's skills by id
func (r *HeistRepository) queryGetMemberSkillsByID(ctx context.Context, id string) ([]storagemodels.MemberSkill, error) {
	row, err := r.dbExecutor.QueryContext(ctx, "SELECT * FROM memberSkills WHERE memberId='"+id+"';")
	if err != nil {
		return []storagemodels.MemberSkill{}, err
	}
	defer row.Close()

	var skills []storagemodels.MemberSkill

	for row.Next() {
		var memberId string
		var skillId string
		var name string
		var level string

		err = row.Scan(&memberId, &skillId, &name, &level)
		if err != nil {
			return []storagemodels.MemberSkill{}, err
		}

		skills = append(skills, storagemodels.MemberSkill{
			MemberId: memberId,
			SkillId:  skillId,
			Name:     name,
			Level:    level,
		})
	}
	return skills, nil
}

// GetMemberSkillsById fetches skills from the database and returns them. The second returned value indicates whether the
// member exists in DB. If the member does not exist, an error will not be returned.
func (r *HeistRepository) GetMemberSkillsById(ctx context.Context, id string) (domainmodels.MemberSkillsDto, bool, error) {
	storageSkills, err := r.queryGetMemberSkillsByID(ctx, id)
	if err != nil {
		return domainmodels.MemberSkillsDto{}, false, err
	}
	var domainSkills domainmodels.MemberSkillsDto
	for idx, skill := range storageSkills {
		domainSkills[idx].Name = skill.Name
		domainSkills[idx].Level = skill.Level
	}

	return domainSkills, true, nil
}

// GetHeistById fetches the heist from the database and returns it. The second returned value indicates whether the
// heist exists in DB. If the heist does not exist, an error will not be returned.
func (r *HeistRepository) GetHeistById(ctx context.Context, id string) (domainmodels.HeistDto, bool, error) {
	storageHeist, err := r.queryGetHeistById(ctx, id)
	storageHeistSkills, err := r.queryGetHeistSkillsByHeistId(ctx, id)
	if err != nil {
		return domainmodels.HeistDto{}, false, err
	}

	domainHeist := domainmodels.HeistDto{
		Name:      storageHeist.Name,
		Location:  storageHeist.Location,
		StartTime: storageHeist.StartTime,
		EndTime:   storageHeist.EndTime,
		Skills:    storageHeistSkills,
		Status:    storageHeist.Status,
	}

	return domainHeist, true, nil
}

// queryGetHeistById dbExecutor call to query the selected heist by id
func (r *HeistRepository) queryGetHeistById(ctx context.Context, id string) (storagemodels.Heist, error) {
	row, err := r.dbExecutor.QueryContext(ctx, "SELECT * FROM members WHERE id='"+id+"';")
	if err != nil {
		return storagemodels.Heist{}, err
	}
	defer row.Close()

	var idx string
	var name string
	var location string
	var startTime time.Time
	var endTime time.Time
	var status string

	err = row.Scan(&idx, &name, &location, &startTime, &endTime, &status)
	if err != nil {
		return storagemodels.Heist{}, err
	}

	var heist = storagemodels.Heist{
		Id:        idx,
		Name:      name,
		Location:  location,
		StartTime: startTime,
		EndTime:   endTime,
		Status:    status,
	}

	return heist, nil
}

// queryGetHeistSkillsByHeistId dbExecutor call to query the selected heist skills by id
func (r *HeistRepository) queryGetHeistSkillsByHeistId(ctx context.Context, id string) (domainmodels.HeistSkillsDto, error) {
	row, err := r.dbExecutor.QueryContext(ctx, "SELECT * FROM heistSkills WHERE heistId='"+id+"';")
	if err != nil {
		return domainmodels.HeistSkillsDto{}, err
	}

	var skills domainmodels.HeistSkillsDto
	idx := 0
	for row.Next() {
		var skillId string
		var heistId string
		var level string
		var members int

		err = row.Scan(&skillId, &heistId, &level, &members)
		if err != nil {
			return domainmodels.HeistSkillsDto{}, err
		}

		name, err := r.queryGetSkillNameById(ctx, skillId)
		if err != nil {
			return domainmodels.HeistSkillsDto{}, err
		}

		skills[idx].Name = name
		skills[idx].Members = members
		skills[idx].Level = level
		idx++
	}
	return skills, nil
}

// GetHeistMembersByHeistId fetches heist members from the database and returns them. The second returned value indicates whether the
// heist exists in DB. If the heist does not exist, an error will not be returned.
func (r *HeistRepository) GetHeistMembersByHeistId(ctx context.Context, id string) ([]domainmodels.MemberDto, bool, error) {
	memberIds, err := r.queryGetMemberIdByHeistId(ctx, id)
	if err != nil {
		return []domainmodels.MemberDto{}, true, err
	}
	var members []storagemodels.Member
	var domainMembers []domainmodels.MemberDto

	for _, oneId := range memberIds {
		member, err := r.queryGetMemberByID(ctx, oneId)
		if err != nil {
			return []domainmodels.MemberDto{}, false, nil
		}
		members = append(members, member)
	}

	for _, member := range members {
		skills, _, err := r.GetMemberSkillsById(ctx, member.Id)
		if err != nil {
			return []domainmodels.MemberDto{}, true, err
		}
		domainMembers = append(domainMembers, domainmodels.MemberDto{
			Name:   member.Name,
			Skills: skills,
		})
	}
	return domainMembers, true, nil
}

// queryGetMemberIdByHeistId dbExecutor call to query the selected heist members by heist id
func (r *HeistRepository) queryGetMemberIdByHeistId(ctx context.Context, id string) ([]string, error) {
	status := "PLANNING"
	row, err := r.dbExecutor.QueryContext(ctx, "SELECT memberId FROM heistMembers WHERE heistId='"+id+"', status!='"+status+"';")
	if err != nil {
		return nil, err
	}

	var ids []string
	for row.Next() {
		var memberId string

		err = row.Scan(&memberId)
		if err != nil {
			return nil, err
		}
		ids = append(ids, memberId)
	}
	return ids, nil
}

// GetHeistSkillsByHeistId fetches heist skills from the database and returns them. The second returned value indicates whether the
// heist exists in DB. If the heist does not exist, an error will not be returned.
func (r *HeistRepository) GetHeistSkillsByHeistId(ctx *gin.Context, id string) (domainmodels.HeistSkillsDto, error) {
	skills, err := r.queryGetHeistSkillsByHeistId(ctx, id)
	if err != nil {
		return domainmodels.HeistSkillsDto{}, err
	}
	return skills, nil
}

// GetHeistStatusByHeistId fetches heist status from the database and returns it.
func (r *HeistRepository) GetHeistStatusByHeistId(ctx *gin.Context, id string) (string, error) {
	row, err := r.dbExecutor.QueryContext(ctx, "SELECT status FROM heists WHERE id='"+id+"';")
	if err != nil {
		return "", err
	}
	var status string

	err = row.Scan(&status)
	if err != nil {
		return "", err
	}
	return status, err
}

// UpdateHeistStatus updates the given heist's status.
func (r *HeistRepository) UpdateHeistStatus(ctx *gin.Context, id, status string) error {
	_, err := r.dbExecutor.Exec("UPDATE heists SET status='" + status + "'WHERE id='" + id + "';")
	if err != nil {
		return err
	}
	return nil
}

// GetHeistOutcomeByHeistId fetches heist status from the database and returns it. The second returned value indicates whether the
// heist exists in DB. If the heist does not exist, an error will not be returned.
func (r *HeistRepository) GetHeistOutcomeByHeistId(ctx *gin.Context, id string) (string, bool, error) {
	row, err := r.dbExecutor.QueryContext(ctx, "SELECT status FROM heists WHERE id='"+id+"';")
	if err != nil {
		return "", false, nil
	}
	var outcome string

	err = row.Scan(&outcome)
	if err != nil {
		return "", true, err
	}
	return outcome, true, err
}
