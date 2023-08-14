package domain

import (
	"errors"
	"regexp"
)

var (
	ErrStudentNotFound = errors.New("student not found")
	ErrStudentExists   = errors.New("student already exist")
)

type Student struct {
	ID           int32
	StudentName  string
	StudentEmail string
	Suspended    bool
}

type NotificationRecipients []string

type Students []Student

func (s *Student) ValidEmail() bool {
	emailRegEx := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`

	match, err := regexp.MatchString(emailRegEx, s.StudentEmail)
	if err != nil {
		return false
	}

	return match
}
