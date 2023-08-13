package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rbo13/go-api-assessment/src/domain"
	"github.com/rbo13/go-api-assessment/src/service"
)

type registerPayload struct {
	Teacher  string   `json:"teacher"`
	Students []string `json:"students"`
}

func (a *api) registerStudent(teacherSrvc service.TeacherService, studentSrvc service.StudentService, registrationSrvc service.RegistrationService) echo.HandlerFunc {
	return func(c echo.Context) error {
		a.logger.Sugar().Info("registerStudent:: Handler Executed")

		json := map[string]interface{}{}

		var payload registerPayload
		if err := c.Bind(&payload); err != nil {
			json["message"] = JSONErrUnexpectedJSONFormat
			return c.JSON(http.StatusUnprocessableEntity, json)
		}

		// if the teacher, does not exist
		// they should not be able to add
		currentTeacher, err := teacherSrvc.RetrieveTeacherByEmail(c.Request().Context(), payload.Teacher)
		if err != nil {
			a.logger.Sugar().Errorf("Something went wrong retrieving teacher due to: %v", err)
			json["message"] = "Teacher does not exist, please create!"
			return c.JSON(http.StatusNotFound, json)
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
					a.logger.Sugar().Errorf("Something went wrong Adding Student: %s, skipping", s.StudentEmail)
					continue
				}

				if err := registrationSrvc.AddRegistration(c.Request().Context(), domain.Registration{
					TeacherID: currentTeacher.ID,
					StudentID: newStudent.ID,
				}); err != nil {
					a.logger.Sugar().Errorf("Something went wrong Registering Student: %s to Teacher: %s, skipping", s.StudentEmail, currentTeacher.Email)
					continue
				}
			} else {
				if err := registrationSrvc.AddRegistration(c.Request().Context(), domain.Registration{
					TeacherID: currentTeacher.ID,
					StudentID: currStudent.ID,
				}); err != nil {
					a.logger.Sugar().Errorf("Something went wrong Registering Current Student: %s to Teacher: %s, skipping", currStudent.StudentEmail, currentTeacher.Email)
					continue
				}
			}
		}

		return c.JSON(http.StatusNoContent, nil)
	}
}
