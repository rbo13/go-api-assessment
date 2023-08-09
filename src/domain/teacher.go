package domain

import "errors"

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
