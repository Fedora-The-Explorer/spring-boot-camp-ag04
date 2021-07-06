package controllers

import (
	"elProfessor/internal/api/controllers/models"
)

// MemberResponse implements member related functions
type HeistResponse interface {
	InsertHeist(heistDto models.HeistDto) error
}