package services

import (
	domainmodels "elProfessor/internal/api/controllers/models"
)

type HeistHandler interface {
	InsertHeist(heistDto domainmodels.HeistDto) error
}