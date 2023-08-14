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
func (m *mockStudentService) FindMentionedStudentsByTeacher(ctx context.Context, teacherEmail string, studentEmails []string) (domain.NotificationRecipients, error) {
	args := m.Called(ctx, teacherEmail, studentEmails)
	return args.Get(0).(domain.NotificationRecipients), args.Error(1)
}

// FindStudentByEmail implements service.StudentService.
func (m *mockStudentService) FindStudentByEmail(ctx context.Context, studentEmail string) (domain.Student, error) {
	args := m.Called(ctx, studentEmail)
	return args.Get(0).(domain.Student), args.Error(1)
}

// SuspendStudent implements service.StudentService.
func (m *mockStudentService) SuspendStudent(ctx context.Context, studentEmail string) error {
	args := m.Called(ctx, studentEmail)
	return args.Error(0)
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
	ctx := context.Background()

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
			On("RetrieveTeacherByEmail", ctx, mock.Anything).
			Return(domain.Teacher{
				Email: "teacher@example.com",
			}, nil)

		studentService.On("FindStudentByEmail", ctx, "student1@example.com").
			Return(domain.Student{
				StudentEmail: "student1@example.com",
			}, nil)

		studentService.On("AddStudent", ctx, mock.Anything).
			Return(domain.Student{
				StudentEmail: "student1@example.com",
			}, nil)

		registrationService.On("AddRegistration", ctx, mock.Anything).
			Return(nil)

		req := httptest.NewRequest(http.MethodPost, "/api/register", strings.NewReader(requestPayload))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		response := httptest.NewRecorder()
		c := e.NewContext(req, response)
		c.SetPath("/api/register")
		c.Set("teacherSrvc", teacherService)
		c.Set("studentSrvc", studentService)
		c.Set("registrationSrvc", registrationService)

		handler := testAPI.RegisterStudent(teacherService, studentService, registrationService)
		assert.NoError(t, handler(c))
		assert.Equal(t, http.StatusNoContent, response.Code)
	})

	t.Run("Should suspend a given Student", func(t *testing.T) {
		requestPayload := `{
			"student": "commonstudent1@gmail.com"
		}`

		studentService.On("FindStudentByEmail", ctx, "commonstudent1@gmail.com").
			Return(domain.Student{
				StudentEmail: "commonstudent1@gmail.com",
			}, nil)

		studentService.On("SuspendStudent", ctx, "commonstudent1@gmail.com").
			Return(nil)

		req := httptest.NewRequest(http.MethodPost, "/api/suspend", strings.NewReader(requestPayload))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		response := httptest.NewRecorder()
		c := e.NewContext(req, response)
		c.SetPath("/api/suspend")
		c.Set("studentSrvc", studentService)

		handler := testAPI.SuspendStudent(studentService)
		assert.NoError(t, handler(c))
		assert.Equal(t, http.StatusNoContent, response.Code)
	})
}
