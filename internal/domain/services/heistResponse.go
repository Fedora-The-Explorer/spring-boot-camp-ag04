package services

import domainmodels "elProfessor/internal/api/controllers/models"

type HeistResponse struct {
	heistHandler HeistHandler
}

func NewHeistResponse(heistHandler HeistHandler) *HeistResponse{
	return &HeistResponse{
		heistHandler: heistHandler,
	}
}

func (m MemberResponse) InsertHeist(heistDto domainmodels.HeistDto) error {
	return m.InsertHeist(heistDto)
}
