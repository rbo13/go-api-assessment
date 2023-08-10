package main

import (
	"context"
	"log"
	"net/http"
	"regexp"

	"github.com/labstack/echo/v4"
	"github.com/rbo13/go-api-assessment/src/db"
	"github.com/rbo13/go-api-assessment/src/domain"
	"github.com/rbo13/go-api-assessment/src/repository/mysql"
	"github.com/rbo13/go-api-assessment/src/service"

	_ "github.com/go-sql-driver/mysql"
)

type registerPayload struct {
	Teacher  string   `json:"teacher"`
	Students []string `json:"students"`
}

type teacherPayload struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

type suspendStudentPayload struct {
	Student string `json:"student"`
}

type sendNotificationPayload struct {
	Teacher      string `json:"teacher"`
	Notification string `json:"notification"`
}

func main() {
	ctx := context.Background()

	conn, err := db.CreateNewConnection(&db.Config{
		Ctx:      ctx,
		MaxConns: 16,
		DSN:      "root:password@tcp(localhost:3306)/api_db?parseTime=true&loc=Local",
	})
	if err != nil {
		log.Fatalf("Cannot start API due to: %v \n", err)
	}
	defer conn.Close()

	e := echo.New()

	e.GET("/api/commonstudents", func(c echo.Context) error {
		json := map[string]interface{}{}

		teacherRepo := mysql.NewTeacherRepository(conn)
		teacherSrvc := service.NewTeacher(teacherRepo)

		params := c.QueryParams()["teacher"]
		if params == nil {
			json["message"] = "Please add a query parameters"
			return c.JSON(http.StatusBadRequest, json)
		}

		res, err := teacherSrvc.RetrieveCommonStudents(c.Request().Context(), params)
		if err != nil {
			return c.JSON(http.StatusNotFound, err)
		}

		json["students"] = res

		return c.JSON(http.StatusOK, json)
	})

	e.POST("/api/teachers", func(c echo.Context) error {
		var payload teacherPayload
		if err := c.Bind(&payload); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		teacherRepo := mysql.NewTeacherRepository(conn)
		teacherSrvc := service.NewTeacher(teacherRepo)

		teacher := domain.Teacher{
			TeacherName: payload.Name,
			Email:       payload.Email,
		}

		if err := teacherSrvc.AddTeacher(c.Request().Context(), teacher); err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}

		return c.JSON(http.StatusCreated, teacher)
	})

	e.POST("/api/register", func(c echo.Context) error {
		var payload registerPayload
		if err := c.Bind(&payload); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		teacherRepo := mysql.NewTeacherRepository(conn)
		teacherSrvc := service.NewTeacher(teacherRepo)

		studentRepo := mysql.NewStudentRepository(conn)
		studentSrvc := service.NewStudent(studentRepo)

		registrationRepo := mysql.NewRegistrationRepository(conn)
		registrationSrvc := service.NewRegistration(registrationRepo)

		// if the teacher, does not exist
		// they should not be able to add
		currentTeacher, err := teacherSrvc.RetrieveTeacherByEmail(c.Request().Context(), payload.Teacher)
		if err != nil {
			return c.JSON(http.StatusNotFound, "Teacher does not exist, please create!")
		}

		for _, student := range payload.Students {

			s := domain.Student{
				StudentEmail: student,
				Suspended:    false,
			}

			currStudent, _ := studentSrvc.FindStudentByEmail(c.Request().Context(), s.StudentEmail)
			if currStudent.ID == 0 {
				newStudent, err := studentSrvc.AddStudent(c.Request().Context(), s)
				if err != nil {
					continue
				}

				if err := registrationSrvc.AddRegistration(c.Request().Context(), domain.Registration{
					TeacherID: currentTeacher.ID,
					StudentID: newStudent.ID,
				}); err != nil {
					continue
				}
			} else {
				if err := registrationSrvc.AddRegistration(c.Request().Context(), domain.Registration{
					TeacherID: currentTeacher.ID,
					StudentID: currStudent.ID,
				}); err != nil {
					continue
				}
			}
		}

		return c.JSON(http.StatusNoContent, nil)
	})

	e.POST("/api/suspend", func(c echo.Context) error {
		var payload suspendStudentPayload
		if err := c.Bind(&payload); err != nil {
			return c.JSON(http.StatusUnprocessableEntity, err.Error())
		}

		studentRepo := mysql.NewStudentRepository(conn)
		studentSrvc := service.NewStudent(studentRepo)

		if err := studentSrvc.SuspendStudent(c.Request().Context(), payload.Student); err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}

		return c.JSON(http.StatusNoContent, nil)
	})

	e.POST("/api/retrievefornotifications", func(c echo.Context) error {
		json := map[string]interface{}{}

		var payload sendNotificationPayload
		if err := c.Bind(&payload); err != nil {
			json["message"] = err.Error()
			return c.JSON(http.StatusUnprocessableEntity, json)
		}

		emails := extractMentionedEmails(payload.Notification)

		studentRepo := mysql.NewStudentRepository(conn)
		studentSrvc := service.NewStudent(studentRepo)

		recipients, err := studentSrvc.FindMentionedStudentsByTeacher(c.Request().Context(), payload.Teacher, emails)
		if err != nil {
			return c.JSON(http.StatusNotFound, err.Error())
		}

		json["recepients"] = recipients
		return c.JSON(http.StatusOK, json)
	})

	e.Logger.Fatal(e.Start(":3000"))
}

func extractMentionedEmails(text string) []string {
	var emails []string

	re, err := regexp.Compile(`@([^\s]+)`)
	if err != nil {
		return emails
	}

	matches := re.FindAllStringSubmatch(text, -1)
	for _, match := range matches {
		if len(match) > 1 {
			emails = append(emails, match[1])
		}
	}
	return emails
}
