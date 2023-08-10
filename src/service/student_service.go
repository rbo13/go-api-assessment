package service

import (
	"context"

	"github.com/rbo13/go-api-assessment/src/domain"
	"github.com/rbo13/go-api-assessment/src/repository"
)

type StudentService interface {
	AddStudent(context.Context, domain.Student) (domain.Student, error)
	FindStudentByEmail(context.Context, string) (domain.Student, error)
}

type studentService struct {
	studentRepository repository.StudentRepository
}

func NewStudent(repo repository.StudentRepository) StudentService {
	return &studentService{studentRepository: repo}
}

// AddStudent implements StudentService.
func (ss *studentService) AddStudent(ctx context.Context, student domain.Student) (domain.Student, error) {
	return ss.studentRepository.Save(ctx, student)
}

// FindStudentByEmail implements StudentService.
func (ss *studentService) FindStudentByEmail(ctx context.Context, email string) (domain.Student, error) {
	return ss.studentRepository.FindByEmail(ctx, email)
}
