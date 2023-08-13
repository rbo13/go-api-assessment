package mysql_test

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/rbo13/go-api-assessment/src/domain"
	"github.com/rbo13/go-api-assessment/src/repository/mysql"
	"github.com/stretchr/testify/assert"
)

func TestMySQLStudentRepository(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		assert.NoError(t, err)
	}
	defer db.Close()

	repo := mysql.NewStudentRepository(db)

	t.Run("Should Create New Student", func(t *testing.T) {
		mockStudent := domain.Student{
			ID:           1,
			StudentEmail: "studentjon@gmail.com",
			Suspended:    false,
		}
		mock.ExpectExec("INSERT INTO students").WithArgs(mockStudent.StudentEmail, mockStudent.Suspended).WillReturnResult(sqlmock.NewResult(1, 1))
		res, err := repo.Save(ctx, mockStudent)
		assert.NotZero(t, res.ID)
		assert.NoError(t, err)
	})

	t.Run("Should Suspend a given Student", func(t *testing.T) {
		studentToBeSuspended := "mock_student@gmail.com"

		mock.ExpectExec("UPDATE students").WithArgs(studentToBeSuspended).WillReturnResult(sqlmock.NewResult(0, 1))
		err := repo.Suspend(ctx, studentToBeSuspended)
		assert.NoError(t, err)
	})

	t.Run("Should Get Student Notifications", func(t *testing.T) {
		teacherEmail := "teacherjoe@gmail.com"
		studentEmails := domain.NotificationRecipients{"studentjon@gmail.com"}

		expectedQuery := `-- name: GetMentionsFromTeacher :many SELECT DISTINCT s.student_email
			FROM students s
			LEFT JOIN registrations r ON s.id = r.student_id
			LEFT JOIN teachers t ON r.teacher_id = t.id
			WHERE (t.email = ? OR s.student_email IN (?)) AND s.suspended = 0`

		expectedRows := sqlmock.NewRows([]string{"student_email"}).
			AddRow("studentjon@gmail.com").
			AddRow("commonstudent1@gmail.com").
			AddRow("commonstudent2@gmail.com").
			AddRow("student_only_under_teacher_joe@gmail.com")

		mock.ExpectQuery(regexp.QuoteMeta(expectedQuery)).
			WithArgs(teacherEmail, studentEmails[0]).
			WillReturnRows(expectedRows)

		recipients, err := repo.GetStudentMentionsFromTeacher(ctx, teacherEmail, studentEmails)
		assert.NoError(t, err)
		assert.NotEmpty(t, recipients)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
