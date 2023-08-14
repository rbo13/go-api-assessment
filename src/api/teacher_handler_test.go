package api_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
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
	panic("unimplemented")
}

// RetrieveTeacherByEmail implements service.TeacherService.
func (m *mockTeacherService) RetrieveTeacherByEmail(ctx context.Context, teacherEmail string) (domain.Teacher, error) {
	panic("unimplemented")
}

func (m *mockTeacherService) RetrieveCommonStudents(ctx context.Context, teachersEmail []string) ([]string, error) {
	args := m.Called(ctx, teachersEmail)
	return args.Get(0).([]string), args.Error(1)
}

func TestTeacherHandler(t *testing.T) {
	e := echo.New()
	logger := logger.New("test_api")

	db, _, err := sqlmock.New()
	if err != nil {
		assert.NoError(t, err)
	}
	defer db.Close()

	teacherService := new(mockTeacherService)

	t.Run("Should GET common Students by a given Teacher", func(t *testing.T) {
		teacherService.
			On("RetrieveCommonStudents", mock.Anything, []string{"teacherken@gmail.com"}).
			Return([]string{"mock_student1@gmail.com", "mock_student2@gmail.com"}, nil)
		expectedResponse := `{"students": ["mock_student1@gmail.com", "mock_student2@gmail.com"]}`

		req := httptest.NewRequest(http.MethodGet, "/api/v1/commonstudents?teacher=teacherken@gmail.com", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/commonstudents")
		c.Set("teacherSrvc", teacherService)

		testAPI := api.New(context.Background(), logger, db)

		handler := testAPI.GetCommonStudents(teacherService)
		assert.NoError(t, handler(c))
		assert.Equal(t, http.StatusOK, rec.Code)

		// assert JSON response
		assert.JSONEq(t, expectedResponse, rec.Body.String())
	})

	t.Run("Should not be able to GET common Students by a given Teacher", func(t *testing.T) {
		teacherService.
			On("RetrieveCommonStudents", mock.Anything, []string{"teacher1"}).
			Return([]string{}, errors.New("No Common Student Found"))

		expectedResponse := `{"message": "No Common Student Found"}`

		req := httptest.NewRequest(http.MethodGet, "/api/v1/commonstudents?teacher=teacher1", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/commonstudents")
		c.Set("teacherSrvc", teacherService)

		testAPI := api.New(context.Background(), logger, db)

		handler := testAPI.GetCommonStudents(teacherService)
		assert.NoError(t, handler(c))
		assert.Equal(t, http.StatusNotFound, rec.Code)

		// assert JSON response
		assert.JSONEq(t, expectedResponse, rec.Body.String())
	})
}
