package domain

import (
	"errors"
	"regexp"
)

var (
	ErrTeacherNotFound = errors.New("teacher not found")
	ErrTeacherExists   = errors.New("teacher already exist")
)

type Teacher struct {
	ID          int32
	TeacherName string
	Email       string
}

type Teachers []Teacher

func (t *Teacher) ValidEmail() bool {
	emailRegEx := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`

	match, err := regexp.MatchString(emailRegEx, t.Email)
	if err != nil {
		return false
	}

	return match
}
