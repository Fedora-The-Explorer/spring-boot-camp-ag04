package bootstrap

import (
	"elProfessor/cmd/config"
	"elProfessor/internal/api"
	"elProfessor/internal/api/controllers"
	"elProfessor/internal/api/controllers/validators"
	"elProfessor/internal/domain/services"
)

func newController(memberResponse controllers.MemberResponse, memberValidator controllers.MemberValidator) *controllers.Controller {
	return controllers.NewController(memberResponse, memberValidator)
}

func newMemberResponse(memberHandler services.MemberHandler) *services.MemberResponse{
	return services.NewMemberResponse(memberHandler)
}

func newMemberValidator() *validators.MemberValidator{
	return validators.NewMemberValidator(memberHandler)
}



func Api() *api.WebServer {
	memberResponse := newMemberResponse()
	controller := newController(memberResponse)

	return api.NewServer(config.Cfg.Api.Port, config.Cfg.Api.ReadWriteTimeoutMs, controller)
}
