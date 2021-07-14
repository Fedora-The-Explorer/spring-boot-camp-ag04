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

// newHeistMapper creates a new mapper
func newHeistMapper() *mappers.HeistMapper {
	return mappers.NewHeistMapper()
}

// newHeistRepository creates a new heist repository
func newHeistRepository(dbExecutor sqlite.DatabaseExecutor, heistMapper sqlite.HeistMapper) *sqlite.HeistRepository {
	return sqlite.NewHeistRepository(dbExecutor, heistMapper)
}

// newMemberResponse creates a new member response
func newMemberResponse(repository *sqlite.HeistRepository) *services.MemberResponse {
	return services.NewMemberResponse(repository)
}

// newHeistResponse creates a new heist response
func newHeistResponse(repository *sqlite.HeistRepository) *services.HeistResponse {
	return services.NewHeistResponse(repository)
}

// newMemberValidator creates a new member validator
func newMemberValidator() *validators.MemberValidator {
	return validators.NewMemberValidator()
}

// newHeistValidator creates a new heist validator
func newHeistValidator() *validators.HeistValidator {
	return validators.NewHeistValidator()
}

// newSmtpService creates the smtp service
func newSmtpService() *smtp.SmtpService {
	return smtp.NewEmailService()
}

// newController creates a new controller
func newController(memberResponse controllers.MemberResponse, heistResponse controllers.HeistResponse, memberValidator controllers.MemberValidator, heistValidator controllers.HeistValidator, smtp controllers.SmtpService) *controllers.Controller {
	return controllers.NewController(memberResponse, heistResponse, memberValidator, heistValidator, smtp)
}

// Api bootstraps the http server
func Api(dbExecutor sqlite.DatabaseExecutor) *api.WebServer {
	mapper := newHeistMapper()
	heistRepository := newHeistRepository(dbExecutor, mapper)
	memberService := newMemberResponse(heistRepository)
	heistService := newHeistResponse(heistRepository)
	memberValidator := newMemberValidator()
	heistValidator := newHeistValidator()
	smtpService := newSmtpService()

	controller := newController(memberService, heistService, memberValidator, heistValidator, smtpService)

	return api.NewServer(config.Cfg.Api.Port, config.Cfg.Api.ReadWriteTimeoutMs, controller)
}
