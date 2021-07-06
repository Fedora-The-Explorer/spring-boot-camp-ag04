package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"net/http"

	"elProfessor/internal/api/controllers/models"
)

// Controller implements handlers for web server requests.
type Controller struct {
	memberResponse MemberResponse
	heistResponse HeistResponse
	memberValidator MemberValidator
}

// NewController creates a new instance of Controller
func NewController(memberResponse MemberResponse, heistResponse HeistResponse ,memberValidator MemberValidator) *Controller {
	return &Controller{
		memberResponse: memberResponse,
		heistResponse: heistResponse,
		memberValidator: memberValidator,
	}
}

// PostMember handles the insert member request and the validation
func (e *Controller) PostMember() gin.HandlerFunc{
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
			ctx.String(http.StatusInternalServerError, "request could not be processed.")
			return
		}
		ctx.Status(http.StatusOK)
	}
}

func (e *Controller) UpdateSkills() gin.HandlerFunc{
	return func(ctx *gin.Context) {
		memberId := ctx.Param("id")
		var memberSkillsUpdateDto models.MemberSkillsUpdateDto
		err := ctx.ShouldBindWith(&memberSkillsUpdateDto, binding.JSON)
		if err != nil {
			ctx.String(http.StatusBadRequest, "put request is not valid")
			return
		}

		if !e.memberValidator.MemberSkillsUpdateValidator(memberSkillsUpdateDto) {
			ctx.String(http.StatusBadRequest, "given member data is not valid")
			return
		}

		err = e.memberResponse.UpdateMemberSkills(ctx, memberSkillsUpdateDto, memberId)
		if err != nil {
			ctx.String(http.StatusInternalServerError, "request could not be processed.")
			return
		}
		ctx.Status(http.StatusOK)
	}
}

func(e *Controller) DeleteSkill() gin.HandlerFunc{
	return func(ctx *gin.Context){
		memberId := ctx.Param("id")
		skillName := ctx.Param("name")

		err:= e.memberResponse.DeleteMemberSkill(memberId, skillName)
		if err != nil {
			ctx.String(http.StatusInternalServerError, "request could not be processed.")
			return
		}
	}
}

func(e *Controller) HeistAdd() gin.HandlerFunc {
	return func(ctx *gin.Context){
		var heistDto models.HeistDto
		err := ctx.ShouldBindWith(&heistDto, binding.JSON)
		if err != nil {
			ctx.String(http.StatusBadRequest, "post request is not valid")
			return
		}

		err = e.heistResponse.InsertHeist(heistDto)
		if err != nil {
			ctx.String(http.StatusInternalServerError, "request could not be processed.")
			return
		}
		ctx.Status(http.StatusOK)
	}
}