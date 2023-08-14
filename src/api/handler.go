package api

import (
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rbo13/go-api-assessment/src/service"
)

const (
	JSONErrUnexpectedJSONFormat = "Unexpected JSON Payload Format. Please check"
	APIVersion                  = "v1.0.5"
)

func (a *api) Handlers() *echo.Echo {
	// initialize different services
	teacherSrvc := service.NewTeacher(a.teacherRepo)
	studentSrvc := service.NewStudent(a.studentRepo)
	registrationSrvc := service.NewRegistration(a.registrationRepo)

	engine := echo.New()
	engine.Use(
		middleware.Recover(),
		middleware.Gzip(),
		middleware.RequestID(),
		middleware.LoggerWithConfig(middleware.LoggerConfig{
			Output: os.Stdout,
		}),
	)

	engine.GET("/", func(c echo.Context) error {
		json := map[string]interface{}{
			"message": "OK",
			"version": APIVersion,
		}

		return c.JSON(http.StatusOK, json)
	})

	engine.GET("/health", a.healthCheckHandler())

	api := engine.Group(apiGroupVersion)

	api.GET(apiCommonStudents, a.GetCommonStudents(teacherSrvc))
	api.POST(apiCreateTeacher, a.CreateTeacher(teacherSrvc))
	api.POST(apiRegisterStudent, a.RegisterStudent(teacherSrvc, studentSrvc, registrationSrvc))
	api.POST(apiSuspend, a.SuspendStudent(studentSrvc))
	api.POST(apiRetrieveNotifications, a.RetrieveForNotifications(teacherSrvc, studentSrvc))

	return engine
}
