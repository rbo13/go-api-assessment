package repository

import (
	"context"

	"github.com/rbo13/go-api-assessment/src/domain"
)

type TeacherRepository interface {
	Save(context.Context, domain.Teacher) error
	FindById(context.Context, int32) (domain.Teacher, error)
	FindByEmail(context.Context, string) (domain.Teacher, error)
	FindCommonStudents(context.Context, []string) ([]string, error)
	List(context.Context) (domain.Teachers, error)
}
