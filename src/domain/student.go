package domain

import "errors"

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

type Students []Student
