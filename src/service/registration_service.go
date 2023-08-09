package service

import (
	"context"

	"github.com/rbo13/go-api-assessment/src/domain"
	"github.com/rbo13/go-api-assessment/src/repository"
)

type RegistrationService interface {
	AddRegistration(context.Context, domain.Registration) error
}

type registrationService struct {
	registrationRepository repository.RegistrationRepository
}

func NewRegistration(repo repository.RegistrationRepository) RegistrationService {
	return &registrationService{registrationRepository: repo}
}

// AddRegistration implements RegistrationService.
func (rs *registrationService) AddRegistration(ctx context.Context, reg domain.Registration) error {
	return rs.registrationRepository.Save(ctx, reg)
}
