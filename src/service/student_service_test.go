package service_test

import (
	"context"
	"testing"

	"github.com/rbo13/go-api-assessment/src/domain"
	"github.com/rbo13/go-api-assessment/src/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockStudentRepository struct {
	mock.Mock
}

// FindByEmail implements repository.StudentRepository.
func (*MockStudentRepository) FindByEmail(context.Context, string) (domain.Student, error) {
	panic("unimplemented")
}

// GetStudentMentionsFromTeacher implements repository.StudentRepository.
func (m *MockStudentRepository) GetStudentMentionsFromTeacher(ctx context.Context, teacherEmail string, studentEmails []string) (domain.NotificationRecipients, error) {
	args := m.Called(ctx, teacherEmail, studentEmails)
	return args.Get(0).(domain.NotificationRecipients), args.Error(1)
}

// Save implements repository.StudentRepository.
func (*MockStudentRepository) Save(context.Context, domain.Student) (domain.Student, error) {
	panic("unimplemented")
}

// Suspend implements repository.StudentRepository.
func (m *MockStudentRepository) Suspend(ctx context.Context, studentEmail string) error {
	args := m.Called(ctx, studentEmail)

	return args.Error(0)
}

func TestStudentService_SuspendAGivenStudent(t *testing.T) {
	ctx := context.Background()

	mockRepo := &MockStudentRepository{}
	studentService := service.NewStudent(mockRepo)

	studentToBeSuspended := "mock_student@gmail.com"

	mockRepo.On("Suspend", ctx, studentToBeSuspended).Return(nil)
	err := studentService.SuspendStudent(ctx, studentToBeSuspended)
	assert.NoError(t, err)
}

func TestStudentService_FindMentionedStudentsByTeacher(t *testing.T) {
	ctx := context.Background()

	mockRepo := &MockStudentRepository{}
	studentService := service.NewStudent(mockRepo)

	teacherEmail := "teacherjoe@gmail.com"
	studentEmails := domain.NotificationRecipients{"studentjon@gmail.com"}

	mockOutput := domain.NotificationRecipients{
		"studentjon@gmail.com",
		"commonstudent1@gmail.com",
		"commonstudent2@gmail.com",
		"student_only_under_teacher_joe@gmail.com",
	}

	mockRepo.On("GetStudentMentionsFromTeacher", ctx, teacherEmail, mock.AnythingOfType("[]string")).Return(mockOutput, nil)

	students, err := studentService.FindMentionedStudentsByTeacher(ctx, teacherEmail, studentEmails)
	assert.NoError(t, err)
	assert.NotEmpty(t, students)

	// Verify that the mock repository methods were called as expected
	mockRepo.AssertExpectations(t)
}
