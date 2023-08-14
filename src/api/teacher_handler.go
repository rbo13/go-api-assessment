package api

import (
	"database/sql"
	"fmt"
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

func (a *api) CreateTeacher(teacherSrvc service.TeacherService) echo.HandlerFunc {
	return func(c echo.Context) error {
		a.logger.Sugar().Info("CreateTeacher:: Handler Executed")

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

func (a *api) GetCommonStudents(teacherSrvc service.TeacherService) echo.HandlerFunc {
	return func(c echo.Context) error {
		a.logger.Sugar().Info("GetCommonStudents:: Handler Executed")

		json := map[string]interface{}{}

		params := c.QueryParams()["teacher"]
		if params == nil {
			json["message"] = "Please add a query parameters"
			return c.JSON(http.StatusBadRequest, json)
		}

		res, err := teacherSrvc.RetrieveCommonStudents(c.Request().Context(), params)
		if err != nil {
			json["message"] = "No Common Student Found"
			return c.JSON(http.StatusNotFound, json)
		}

		json["students"] = res

		return c.JSON(http.StatusOK, json)
	}
}

func (a *api) SuspendStudent(studentSrvc service.StudentService) echo.HandlerFunc {
	return func(c echo.Context) error {
		a.logger.Sugar().Info("SuspendStudent:: Handler Executed")

		json := map[string]interface{}{}

		var payload suspendStudentPayload
		if err := c.Bind(&payload); err != nil {
			json["message"] = JSONErrUnexpectedJSONFormat
			return c.JSON(http.StatusUnprocessableEntity, json)
		}

		student, err := studentSrvc.FindStudentByEmail(c.Request().Context(), payload.Student)
		if err != nil && err != sql.ErrNoRows {
			a.logger.Sugar().Errorf("Something went wrong FindStudentByEmail:: %v \n", err)
			json["message"] = err
			return c.JSON(http.StatusInternalServerError, json)
		}

		if student.StudentEmail == "" {
			json["message"] = "Student not found"
			return c.JSON(http.StatusNotFound, json)
		}

		if err := studentSrvc.SuspendStudent(c.Request().Context(), student.StudentEmail); err != nil {
			json["message"] = fmt.Sprintf("Cannot suspend student due to: %v \n", err)
			return c.JSON(http.StatusBadRequest, json)
		}

		return c.JSON(http.StatusNoContent, nil)
	}
}

func (a *api) RetrieveForNotifications(teacherSrvc service.TeacherService, studentSrvc service.StudentService) echo.HandlerFunc {
	return func(c echo.Context) error {
		a.logger.Sugar().Info("RetrieveForNotifications:: Handler Executed")

		json := map[string]interface{}{}

		var payload sendNotificationPayload
		if err := c.Bind(&payload); err != nil {
			json["message"] = JSONErrUnexpectedJSONFormat
			return c.JSON(http.StatusUnprocessableEntity, json)
		}

		// check first if the teacher exist before sending notification
		teacher, err := teacherSrvc.RetrieveTeacherByEmail(c.Request().Context(), payload.Teacher)
		if err != nil && err != sql.ErrNoRows {
			a.logger.Sugar().Errorf("Something went wrong RetrieveTeacherByEmail:: %v \n", err)
			json["message"] = err
			return c.JSON(http.StatusInternalServerError, json)
		}

		if teacher.Email == "" {
			json["message"] = "Cannot retrieve notification. Teacher not found"
			return c.JSON(http.StatusNotFound, json)
		}

		emails := extractMentionedEmails(payload.Notification)

		recipients, err := studentSrvc.FindMentionedStudentsByTeacher(c.Request().Context(), teacher.Email, emails)
		if err != nil {
			return c.JSON(http.StatusNotFound, err.Error())
		}

		json["recipients"] = recipients
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
