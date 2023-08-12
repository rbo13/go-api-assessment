package mysql_test

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/rbo13/go-api-assessment/src/domain"
	"github.com/rbo13/go-api-assessment/src/repository/mysql"
	"github.com/stretchr/testify/assert"
)

func TestMySQLSRegistrationRepository(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		assert.NoError(t, err)
	}
	defer db.Close()

	repo := mysql.NewRegistrationRepository(db)

	t.Run("Should Register Student under Teacher", func(t *testing.T) {
		student := domain.Student{
			ID:           1,
			StudentEmail: "studentjon@gmail.com",
			Suspended:    false,
		}

		teacher := domain.Teacher{
			ID:          1,
			TeacherName: "Teacher Ken",
			Email:       "teacherken@gmail.com",
		}

		mockRegistration := domain.Registration{
			ID:        1,
			TeacherID: teacher.ID,
			StudentID: student.ID,
		}

		mock.ExpectExec("INSERT INTO registrations").WithArgs(teacher.ID, student.ID).WillReturnResult(sqlmock.NewResult(1, 1))
		err = repo.Save(ctx, mockRegistration)
		assert.NoError(t, err, "Error Registering Student to a Teacher")
	})
}
