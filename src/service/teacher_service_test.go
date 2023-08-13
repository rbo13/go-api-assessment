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

func (m *MockTeacherRepository) FindCommonStudents(ctx context.Context, emails []string) ([]string, error) {
	args := m.Called(ctx, emails)
	return args.Get(0).([]string), args.Error(1)
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

func TestTeacherService_GetCommonStudents(t *testing.T) {
	ctx := context.Background()

	mockRepo := &MockTeacherRepository{}
	teacherService := service.NewTeacher(mockRepo)

	emails := []string{"teacherken@gmail.com"}
	commonStudents := []string{
		"studenthon@gmail.com",
		"studentjon@gmail.com",
	}

	mockRepo.On("FindCommonStudents", ctx, mock.Anything).Return(commonStudents, nil)

	students, err := teacherService.RetrieveCommonStudents(ctx, emails)
	assert.NoError(t, err)
	assert.Equal(t, commonStudents, students)

	// Verify that the mock repository methods were called as expected
	mockRepo.AssertExpectations(t)
}
