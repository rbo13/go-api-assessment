package service_test

import (
	"context"
	"testing"

	"github.com/rbo13/go-api-assessment/src/domain"
	"github.com/rbo13/go-api-assessment/src/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockTeacherRepository struct {
	mock.Mock
}

// FindByEmail implements repository.TeacherRepository.
func (m *MockTeacherRepository) FindByEmail(ctx context.Context, email string) (domain.Teacher, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(domain.Teacher), args.Error(1)
}

// FindById implements repository.TeacherRepository.
func (m *MockTeacherRepository) FindById(context.Context, int32) (domain.Teacher, error) {
	panic("unimplemented")
}

// List implements repository.TeacherRepository.
func (m *MockTeacherRepository) List(context.Context) (domain.Teachers, error) {
	panic("unimplemented")
}

// Save implements repository.TeacherRepository.
func (m *MockTeacherRepository) Save(ctx context.Context, teacher domain.Teacher) error {
	args := m.Called(ctx, teacher)
	return args.Error(0)
}

func TestTeacherService_AddTeacher(t *testing.T) {
	ctx := context.Background()

	mockRepo := &MockTeacherRepository{}
	teacherService := service.NewTeacher(mockRepo)

	teacher := domain.Teacher{
		ID:          1,
		TeacherName: "Test Teacher One",
		Email:       "testteacher@example.com",
	}

	mockRepo.On("Save", ctx, teacher).Return(nil)

	err := teacherService.AddTeacher(ctx, teacher)
	assert.NoError(t, err)

	// Verify that the mock repository methods were called as expected
	mockRepo.AssertExpectations(t)
}
