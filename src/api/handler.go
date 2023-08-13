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
		return c.JSON(http.StatusOK, "OK")
	})

	engine.GET("/health", a.healthCheckHandler())

	api := engine.Group(apiGroupVersion)

	api.GET(apiCommonStudents, a.getCommonStudents(teacherSrvc))
	api.POST(apiCreateTeacher, a.createTeacher(teacherSrvc))
	api.POST(apiRegisterStudent, a.registerStudent(teacherSrvc, studentSrvc, registrationSrvc))
	api.POST(apiSuspend, a.suspendStudent(studentSrvc))
	api.POST(apiRetrieveNotifications, a.retrieveForNotifications(teacherSrvc, studentSrvc))

	return engine
}
