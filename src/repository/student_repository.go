package repository

import (
	"context"

	"github.com/rbo13/go-api-assessment/src/domain"
)

type StudentRepository interface {
	Save(context.Context, domain.Student) (domain.Student, error)
	FindByEmail(context.Context, string) (domain.Student, error)
	Suspend(context.Context, string) error
}
