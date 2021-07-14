package controllers

import (
	"context"
	"elProfessor/internal/api/controllers/models"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/jasonlvhit/gocron"
	"net/http"
	"time"
)

// Controller implements handlers for web server requests.
type Controller struct {
	memberResponse  MemberResponse
	heistResponse   HeistResponse
	memberValidator MemberValidator
	heistValidator  HeistValidator
	smtpService     SmtpService
}

// NewController creates a new instance of Controller
func NewController(memberResponse MemberResponse, heistResponse HeistResponse, memberValidator MemberValidator, heistValidator HeistValidator, smtpService SmtpService) *Controller {
	return &Controller{
		memberResponse:  memberResponse,
		heistResponse:   heistResponse,
		memberValidator: memberValidator,
		heistValidator:  heistValidator,
		smtpService:     smtpService,
	}
}

// PostMember handles the insert member request and the validation, sends mail to newly added members
func (e *Controller) PostMember() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var memberDto models.MemberDto
		err := ctx.ShouldBindWith(&memberDto, binding.JSON)
		if err != nil {
			ctx.String(http.StatusBadRequest, "post request is not valid")
			return
		}

		if !e.memberValidator.MemberIsValid(memberDto) {
			ctx.String(http.StatusBadRequest, "given member data is not valid")
			return
		}

		err = e.memberResponse.InsertMember(memberDto)
		if err != nil {
			ctx.String(http.StatusBadRequest, "request could not be processed.")
			return
		}
		message := []byte("You have been added as a possible choice for a heist member")
		to := []string{
			memberDto.Email,
		}
		e.smtpService.SendEmail(to, message)
		ctx.Status(http.StatusCreated)
	}
}

// UpdateSkills handles the update member skills request
func (e *Controller) UpdateSkills() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		memberId := ctx.Param("id")
		var memberSkillsUpdateDto models.MemberSkillsUpdateDto
		err := ctx.ShouldBindWith(&memberSkillsUpdateDto, binding.JSON)
		if err != nil {
			ctx.String(http.StatusBadRequest, "put request is not valid")
			return
		}

		if !e.memberValidator.MemberSkillsUpdateValidator(memberSkillsUpdateDto) {
			ctx.String(http.StatusBadRequest, "given skill data is not valid")
			return
		}

		err = e.memberResponse.UpdateMemberSkills(ctx, memberSkillsUpdateDto, memberId)
		if err != nil {
			ctx.String(http.StatusNotFound, "request could not be processed.")
			return
		}
		ctx.Status(http.StatusNoContent)
	}
}

// DeleteSkill handles the deletion of one of the member's skills
func (e *Controller) DeleteSkill() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		memberId := ctx.Param("id")
		skillName := ctx.Param("name")

		err := e.memberResponse.DeleteMemberSkill(memberId, skillName)
		if err != nil {
			ctx.String(http.StatusNotFound, "request could not be processed.")
			return
		}
		ctx.Status(http.StatusNoContent)
	}
}

// PostHeist creates a new heist and validates it
// Also it automatically calls the start and end heist at the right moment
func (e *Controller) PostHeist() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var heistDto models.HeistDto
		err := ctx.ShouldBindWith(&heistDto, binding.JSON)
		if err != nil {
			ctx.String(http.StatusBadRequest, "post request is not valid")
			return
		}

		id, err := e.heistResponse.InsertHeist(heistDto)
		if err != nil {
			ctx.String(http.StatusBadRequest, "request could not be processed.")
			return
		}

		//We will be using go cron package for automation of starting and ending heists at specific times
		startCh := make(chan bool)
		startTime := heistDto.StartTime
		endTime := heistDto.EndTime
		e.AutomaticStart(id, startCh, startTime)
		endCh := make(chan bool)
		e.AutomaticEnd(id, endCh, endTime)

		ctx.Status(http.StatusCreated)
	}
}

// AutomaticStart automatically starts a heist depending on the heist's start time and
// email to the heist members that the heist has finished
func (e *Controller) AutomaticStart(id string, quit <-chan bool, time time.Time, ) {
	g := gocron.NewScheduler()
	members, _, err := e.heistResponse.GetHeistMembersByHeistId(context.Background(), id)
	if err != nil {
		return
	}
	err = g.Every(1).Hour().From(&time).Do(e.heistResponse.StartHeist(id))
	if err != nil {
		return
	}

	var emails []string
	for _, member := range members {
		emails = append(emails, member.Email)
	}

	message := []byte("The heist has started!")
	e.smtpService.SendEmail(emails, message)

	<-quit
	return
}

// AutomaticEnd automatically ends a heist depending on the heist's end time and sends and
// email to the heist members that the heist has finished
func (e *Controller) AutomaticEnd(id string, quit <-chan bool, time time.Time) {
	g := gocron.NewScheduler()
	members, _, err := e.heistResponse.GetHeistMembersByHeistId(context.Background(), id)
	if err != nil {
		return
	}
	err = g.Every(1).Hour().From(&time).Do(e.heistResponse.EndHeist(id))
	if err != nil {
		return
	}
	var emails []string
	for _, member := range members {
		emails = append(emails, member.Email)
	}

	message := []byte("The heist has ended!")
	e.smtpService.SendEmail(emails, message)

	<-quit
	return
}

// UpdateHeistSkills handles the update request for certain heist's skills
func (e *Controller) UpdateHeistSkills() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		heistId := ctx.Param("id")
		var heistSkills models.HeistSkillsDto
		err := ctx.ShouldBindWith(&heistSkills, binding.JSON)
		if err != nil {
			ctx.String(http.StatusBadRequest, "patch request is not valid")
		}

		if !e.heistValidator.HeistSkillUpdateValidator(heistSkills) {
			ctx.String(http.StatusBadRequest, "given skill data is not valid")
			return
		}

		err = e.heistResponse.UpdateHeistSkills(ctx, heistSkills, heistId)
		if err != nil {
			ctx.String(http.StatusMethodNotAllowed, "request could not be processed.")
			return
		}
		ctx.Status(http.StatusNoContent)

	}
}

// EligibleMembers returns the members that are eligible for the given heist
func (e *Controller) EligibleMembers() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")
		members, exists, err := e.memberResponse.GetEligibleMembers(ctx, id)
		if err != nil {
			ctx.String(http.StatusMethodNotAllowed, "request could not be processed")
			return
		}
		if !exists {
			ctx.String(http.StatusNotFound, "failed to get heist with given id")
			return
		}
		ctx.JSON(http.StatusOK, members)
	}
}

// AddMembersToHeist allows us to add members to a certain heist, when we add
//them, those members will be notified via email
func (e *Controller) AddMembersToHeist() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")
		var members []string
		err := ctx.ShouldBindWith(&members, binding.JSON)
		if err != nil {
			ctx.String(http.StatusBadRequest, "request could not be processed")
			return
		}

		code, err, mails := e.heistResponse.AddHeistMembers(members, id)
		if err != nil {
			if code == "404" {
				ctx.String(http.StatusNotFound, "heist not found")
				return
			} else {
				ctx.String(http.StatusMethodNotAllowed, "not allowed due to wrong heist status")
				return
			}
		}

		message := []byte("You have been added as a member of a heist")
		e.smtpService.SendEmail(mails, message)

		ctx.Status(http.StatusNoContent)
	}
}

// StartHeist allows us to manually start a heist, the members of the heist will also be notified via email
func (e *Controller) StartHeist() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")
		code, err := e.heistResponse.StartHeist(id)
		if err != nil {
			if code == "404" {
				ctx.String(http.StatusNotFound, "heist not found")
				return
			} else {
				ctx.String(http.StatusMethodNotAllowed, "not allowed due to wrong heist status")
				return
			}
		}
		ctx.Status(http.StatusOK)

	}
}

// GetMember show us the given member's information
func (e *Controller) GetMember() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")
		member, exists, err := e.memberResponse.GetMemberById(ctx, id)

		if err != nil {
			ctx.String(http.StatusBadRequest, "request could not be processed")
			return
		}
		if !exists {
			ctx.String(http.StatusNotFound, "failed to get member with given id")
			return
		}

		ctx.JSON(http.StatusOK, member)
	}
}

// GetMemberSkills show us the given member skills' information
func (e *Controller) GetMemberSkills() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")
		skills, exists, err := e.memberResponse.GetMemberSkillsById(ctx, id)
		if err != nil {
			ctx.String(http.StatusBadRequest, "request could not be processed")
			return
		}
		if !exists {
			ctx.String(http.StatusNotFound, "failed to get member skills with given member id")
			return
		}

		ctx.JSON(http.StatusOK, skills)
	}
}

// GetHeist show us the given heist's information
func (e *Controller) GetHeist() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")
		heist, exists, err := e.heistResponse.GetHeistById(ctx, id)
		if err != nil {
			ctx.String(http.StatusBadRequest, "request could not be processed")
			return
		}
		if !exists {
			ctx.String(http.StatusNotFound, "failed to get heist for a given heist id")
			return
		}

		ctx.JSON(http.StatusOK, heist)
	}
}

// GetHeistMembers show us the given heist's member information
func (e *Controller) GetHeistMembers() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")
		members, exists, err := e.heistResponse.GetHeistMembersByHeistId(ctx, id)
		if err != nil {
			ctx.String(http.StatusMethodNotAllowed, "request could not be processed")
			return
		}
		if !exists {
			ctx.String(http.StatusNotFound, "failed to get member skills with given member id")
			return
		}
		ctx.JSON(http.StatusOK, members)
	}
}

// GetHeistSkills show us the given heist's skill information
func (e *Controller) GetHeistSkills() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")
		skills, err := e.heistResponse.GetHeistSkillsByHeistId(ctx, id)
		if err != nil {
			ctx.String(http.StatusNotFound, "request could not be processed")
			return
		}

		ctx.JSON(http.StatusOK, skills)
	}
}

// GetHeistStatus shows us the status of a certain heist
func (e *Controller) GetHeistStatus() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")
		status, err := e.heistResponse.GetHeistStatusByHeistId(ctx, id)
		if err != nil {
			ctx.String(http.StatusNotFound, "request could not be processed")
			return
		}

		ctx.JSON(http.StatusOK, status)
	}
}

// GetHeistOutcome shows us the outcome of a certain heist
func (e *Controller) GetHeistOutcome() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")
		status, exists, err := e.heistResponse.GetHeistOutcomeByHeistId(ctx, id)
		if err != nil {
			ctx.String(http.StatusMethodNotAllowed, "request could not be processed")
			return
		}
		if !exists {
			ctx.String(http.StatusNotFound, "heist not found")
			return
		}

		ctx.JSON(http.StatusOK, status)
	}
}
