package api

import (
	"net/http"
	"regexp"

	"github.com/labstack/echo/v4"
	"github.com/rbo13/go-api-assessment/src/domain"
	"github.com/rbo13/go-api-assessment/src/service"
)

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

func (a *api) createTeacher(teacherSrvc service.TeacherService) echo.HandlerFunc {
	return func(c echo.Context) error {
		var payload teacherPayload
		if err := c.Bind(&payload); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		teacher := domain.Teacher{
			TeacherName: payload.Name,
			Email:       payload.Email,
		}

		if err := teacherSrvc.AddTeacher(c.Request().Context(), teacher); err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}

		return c.JSON(http.StatusCreated, teacher)
	}
}

func (a *api) getCommonStudents(teacherSrvc service.TeacherService) echo.HandlerFunc {
	return func(c echo.Context) error {
		json := map[string]interface{}{}

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
	}
}

func (a *api) suspendStudent(studentSrvc service.StudentService) echo.HandlerFunc {
	return func(c echo.Context) error {
		var payload suspendStudentPayload
		if err := c.Bind(&payload); err != nil {
			return c.JSON(http.StatusUnprocessableEntity, err.Error())
		}

		if err := studentSrvc.SuspendStudent(c.Request().Context(), payload.Student); err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}

		return c.JSON(http.StatusNoContent, nil)
	}
}

func (a *api) retrieveForNotifications(studentSrvc service.StudentService) echo.HandlerFunc {
	return func(c echo.Context) error {
		json := map[string]interface{}{}

		var payload sendNotificationPayload
		if err := c.Bind(&payload); err != nil {
			json["message"] = err.Error()
			return c.JSON(http.StatusUnprocessableEntity, json)
		}

		emails := extractMentionedEmails(payload.Notification)

		recipients, err := studentSrvc.FindMentionedStudentsByTeacher(c.Request().Context(), payload.Teacher, emails)
		if err != nil {
			return c.JSON(http.StatusNotFound, err.Error())
		}

		json["recepients"] = recipients
		return c.JSON(http.StatusOK, json)
	}
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
