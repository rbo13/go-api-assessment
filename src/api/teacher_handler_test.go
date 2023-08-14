package api_test

import (
	"context"
	"errors"
	"fmt"
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

type mockTeacherService struct {
	mock.Mock
}

// AddTeacher implements service.TeacherService.
func (m *mockTeacherService) AddTeacher(ctx context.Context, teacher domain.Teacher) error {
	args := m.Called(ctx, teacher)
	return args.Error(0)
}

// RetrieveTeacherByEmail implements service.TeacherService.
func (m *mockTeacherService) RetrieveTeacherByEmail(ctx context.Context, teacherEmail string) (domain.Teacher, error) {
	args := m.Called(ctx, teacherEmail)
	return args.Get(0).(domain.Teacher), args.Error(1)
}

func (m *mockTeacherService) RetrieveCommonStudents(ctx context.Context, teachersEmail []string) ([]string, error) {
	args := m.Called(ctx, teachersEmail)
	return args.Get(0).([]string), args.Error(1)
}

func TestTeacherHandler(t *testing.T) {
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

	testAPI := api.New(ctx, logger, db)

	t.Run("Should Create a new Teacher", func(t *testing.T) {
		testNewTeacher := domain.Teacher{
			Email: "testnewteacher@gmail.com",
		}

		mockPayload := `{
			"email": "%s"
		}`

		requestPayload := fmt.Sprintf(mockPayload, testNewTeacher.Email)

		teacherService.On("AddTeacher", ctx, testNewTeacher).Return(nil)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/teachers", strings.NewReader(requestPayload))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		response := httptest.NewRecorder()
		c := e.NewContext(req, response)
		c.SetPath("/api/v1/teachers")
		c.Set("teacherSrvc", teacherService)

		handler := testAPI.CreateTeacher(teacherService)
		assert.NoError(t, handler(c))
		assert.Equal(t, http.StatusCreated, response.Code)
	})

	t.Run("Should GET common Students by a given Teacher", func(t *testing.T) {
		teacherService.
			On("RetrieveCommonStudents", ctx, []string{"teacherken@gmail.com"}).
			Return([]string{"mock_student1@gmail.com", "mock_student2@gmail.com"}, nil)
		expectedResponse := `{"students": ["mock_student1@gmail.com", "mock_student2@gmail.com"]}`

		req := httptest.NewRequest(http.MethodGet, "/api/v1/commonstudents?teacher=teacherken@gmail.com", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/commonstudents")
		c.Set("teacherSrvc", teacherService)

		handler := testAPI.GetCommonStudents(teacherService)
		assert.NoError(t, handler(c))
		assert.Equal(t, http.StatusOK, rec.Code)

		// assert JSON response
		assert.JSONEq(t, expectedResponse, rec.Body.String())
	})

	t.Run("Should not be able to GET common Students by a given Teacher", func(t *testing.T) {
		teacherService.
			On("RetrieveCommonStudents", ctx, []string{"teacher1"}).
			Return([]string{}, errors.New("No Common Student Found"))

		expectedResponse := `{"message": "No Common Student Found"}`

		req := httptest.NewRequest(http.MethodGet, "/api/v1/commonstudents?teacher=teacher1", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/commonstudents")
		c.Set("teacherSrvc", teacherService)

		handler := testAPI.GetCommonStudents(teacherService)
		assert.NoError(t, handler(c))
		assert.Equal(t, http.StatusNotFound, rec.Code)

		// assert JSON response
		assert.JSONEq(t, expectedResponse, rec.Body.String())
	})

	t.Run("Should Retrieve Notification For Students", func(t *testing.T) {
		requestPayload := `{
			"teacher":  "teacherjoe@gmail.com",
			"notification": "Hey! @mock_student1@gmail.com"
		}`

		teacher := "teacherjoe@gmail.com"
		recipients := []string{"mock_student1@gmail.com"}

		teacherService.
			On("RetrieveTeacherByEmail", ctx, teacher).
			Return(domain.Teacher{
				Email: teacher,
			}, nil)

		studentService.On("FindMentionedStudentsByTeacher", ctx, teacher, recipients).
			Return(domain.NotificationRecipients{
				"mock_student1@gmail.com",
			}, nil)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/retrievefornotifications", strings.NewReader(requestPayload))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		response := httptest.NewRecorder()
		c := e.NewContext(req, response)
		c.SetPath("/api/v1/retrievefornotifications")
		c.Set("teacherSrvc", teacherService)
		c.Set("studentSrvc", studentService)

		handler := testAPI.RetrieveForNotifications(teacherService, studentService)
		assert.NoError(t, handler(c))
		assert.Equal(t, http.StatusOK, response.Code)
	})
}
