package repository

import (
	"context"

	"github.com/rbo13/go-api-assessment/src/domain"
)

type RegistrationRepository interface {
	Save(context.Context, domain.Registration) error
	FindById(context.Context, int32) (domain.Registration, error)
}
