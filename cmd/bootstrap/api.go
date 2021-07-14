package bootstrap

import (
	"elProfessor/cmd/config"
	"elProfessor/internal/api"
	"elProfessor/internal/api/controllers"
	"elProfessor/internal/api/controllers/smtp"
	"elProfessor/internal/api/controllers/validators"
	"elProfessor/internal/domain/mappers"
	"elProfessor/internal/domain/services"
	"elProfessor/internal/infrastructure/sqlite"
)

func newHeistMapper() *mappers.HeistMapper{
	return mappers.NewHeistMapper()
}

func newHeistRepository(dbExecutor sqlite.DatabaseExecutor, heistMapper sqlite.HeistMapper) *sqlite.HeistRepository{
	return sqlite.NewHeistRepository(dbExecutor, heistMapper)
}

func newMemberResponse(repository *sqlite.HeistRepository) *services.MemberResponse{
	return services.NewMemberResponse(repository)
}

func newHeistResponse(repository *sqlite.HeistRepository) *services.HeistResponse{
	return services.NewHeistResponse(repository)
}

func newMemberValidator() *validators.MemberValidator{
	return validators.NewMemberValidator()
}

func newHeistValidator() *validators.HeistValidator{
	return validators.NewHeistValidator()
}

func newSmtpService() *smtp.SmtpService{
	return smtp.NewEmailService()
}

func newController(memberResponse controllers.MemberResponse, heistResponse controllers.HeistResponse, memberValidator controllers.MemberValidator, heistValidator controllers.HeistValidator, smtp controllers.SmtpService) *controllers.Controller{
	return controllers.NewController(memberResponse, heistResponse, memberValidator, heistValidator, smtp)
}

func Api(dbExecutor sqlite.DatabaseExecutor) *api.WebServer {
	mapper := newHeistMapper()
	heistRepository:= newHeistRepository(dbExecutor, mapper)
	memberService := newMemberResponse(heistRepository)
	heistService := newHeistResponse(heistRepository)
	memberValidator := newMemberValidator()
	heistValidator := newHeistValidator()
	smtpService := newSmtpService()

	controller := newController(memberService,heistService,memberValidator,heistValidator,smtpService)

	return api.NewServer(config.Cfg.Api.Port, config.Cfg.Api.ReadWriteTimeoutMs, controller)
}
