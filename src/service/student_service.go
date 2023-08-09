package service

import (
	"context"

	"github.com/rbo13/go-api-assessment/src/domain"
	"github.com/rbo13/go-api-assessment/src/repository"
)

type StudentService interface {
	AddStudent(context.Context, domain.Student) (domain.Student, error)
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
