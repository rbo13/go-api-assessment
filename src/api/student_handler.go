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
		var payload registerPayload
		if err := c.Bind(&payload); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

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
	}
}
