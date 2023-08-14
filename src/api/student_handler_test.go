package api_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/labstack/echo/v4"
	"github.com/rbo13/go-api-assessment/src/api"
	"github.com/rbo13/go-api-assessment/src/domain"
	"github.com/rbo13/go-api-assessment/src/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockStudentService struct {
	mock.Mock
}

// AddStudent implements service.StudentService.
func (m *mockStudentService) AddStudent(ctx context.Context, student domain.Student) (domain.Student, error) {
	args := m.Called(ctx, student)
	return args.Get(0).(domain.Student), args.Error(1)
}

// FindMentionedStudentsByTeacher implements service.StudentService.
func (*mockStudentService) FindMentionedStudentsByTeacher(context.Context, string, []string) (domain.NotificationRecipients, error) {
	panic("unimplemented")
}

// FindStudentByEmail implements service.StudentService.
func (m *mockStudentService) FindStudentByEmail(ctx context.Context, studentEmail string) (domain.Student, error) {
	args := m.Called(ctx, studentEmail)
	return args.Get(0).(domain.Student), args.Error(1)
}

// SuspendStudent implements service.StudentService.
func (*mockStudentService) SuspendStudent(context.Context, string) error {
	panic("unimplemented")
}

type mockRegistrationService struct {
	mock.Mock
}

// AddRegistration implements service.RegistrationService.
func (m *mockRegistrationService) AddRegistration(ctx context.Context, reg domain.Registration) error {
	args := m.Called(ctx, reg)
	return args.Error(0)
}

func TestStudentHandler(t *testing.T) {
	e := echo.New()
	logger := logger.New("test_api")

	db, _, err := sqlmock.New()
	if err != nil {
		assert.NoError(t, err)
	}
	defer db.Close()

	teacherService := new(mockTeacherService)
	studentService := new(mockStudentService)
	registrationService := new(mockRegistrationService)

	testAPI := api.New(context.Background(), logger, db)

	t.Run("Should register a Student to Teacher", func(t *testing.T) {
		requestPayload := `{
			"teacher": "teacher@example.com",
			"students": ["student1@example.com"]
		}`

		teacherService.
			On("RetrieveTeacherByEmail", mock.Anything, mock.Anything).
			Return(domain.Teacher{
				Email: "teacher@example.com",
			}, nil)

		studentService.On("FindStudentByEmail", mock.Anything, "student1@example.com").
			Return(domain.Student{
				StudentEmail: "student1@example.com",
			}, nil)

		studentService.On("AddStudent", mock.Anything, mock.Anything).
			Return(domain.Student{
				StudentEmail: "student1@example.com",
			}, nil)

		registrationService.On("AddRegistration", mock.Anything, mock.Anything).
			Return(nil)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/register", strings.NewReader(requestPayload))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		response := httptest.NewRecorder()
		c := e.NewContext(req, response)
		c.SetPath("/api/v1/register")
		c.Set("teacherSrvc", teacherService)
		c.Set("studentSrvc", studentService)
		c.Set("registrationSrvc", registrationService)

		handler := testAPI.RegisterStudent(teacherService, studentService, registrationService)
		assert.NoError(t, handler(c))
		assert.Equal(t, http.StatusNoContent, response.Code)
	})
}
