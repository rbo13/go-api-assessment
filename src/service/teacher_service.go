package service

import (
	"context"

	"github.com/rbo13/go-api-assessment/src/domain"
	"github.com/rbo13/go-api-assessment/src/repository"
)

type TeacherService interface {
	AddTeacher(context.Context, domain.Teacher) error
	RetrieveTeacherByEmail(context.Context, string) (domain.Teacher, error)
	RetrieveCommonStudents(context.Context, []string) ([]string, error)
}

type teacherService struct {
	teacherRepository repository.TeacherRepository
}

func NewTeacher(repo repository.TeacherRepository) TeacherService {
	return &teacherService{teacherRepository: repo}
}

// AddTeacher implements TeacherService.
func (ts *teacherService) AddTeacher(ctx context.Context, teacher domain.Teacher) error {
	return ts.teacherRepository.Save(ctx, teacher)
}

// RetrieveTeacherByEmail implements TeacherService.
func (ts *teacherService) RetrieveTeacherByEmail(ctx context.Context, email string) (domain.Teacher, error) {
	return ts.teacherRepository.FindByEmail(ctx, email)
}

func (ts *teacherService) RetrieveCommonStudents(ctx context.Context, emails []string) ([]string, error) {
	return ts.teacherRepository.FindCommonStudents(ctx, emails)
}
