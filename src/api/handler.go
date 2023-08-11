package api

import (
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rbo13/go-api-assessment/src/service"
)

func (a *api) Handlers() *echo.Echo {
	engine := echo.New()

	engine.Use(
		middleware.Recover(),
		middleware.Gzip(),
		middleware.RequestID(),
		middleware.LoggerWithConfig(middleware.LoggerConfig{
			Output: os.Stdout,
		}),
	)

	teacherSrvc := service.NewTeacher(a.teacherRepo)
	studentSrvc := service.NewStudent(a.studentRepo)
	registrationSrvc := service.NewRegistration(a.registrationRepo)

	engine.GET(apiCommonStudents, a.getCommonStudents(teacherSrvc))
	engine.POST(apiCreateTeacher, a.createTeacher(teacherSrvc))
	engine.POST(apiRegisterStudent, a.registerStudent(teacherSrvc, studentSrvc, registrationSrvc))
	engine.POST(apiSuspend, a.suspendStudent(studentSrvc))
	engine.POST(apiRetrieveNotifications, a.retrieveForNotifications(studentSrvc))

	return engine
}
