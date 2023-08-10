package mysql_test

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/rbo13/go-api-assessment/src/domain"
	"github.com/rbo13/go-api-assessment/src/repository/mysql"
	"github.com/stretchr/testify/assert"
)

var ctx = context.Background()

func TestMySQLTeacherRepository(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		assert.NoError(t, err)
	}
	defer db.Close()

	repo := mysql.NewTeacherRepository(db)

	t.Run("Should Create New Teacher", func(t *testing.T) {
		mockTeacher := domain.Teacher{
			ID:          1,
			TeacherName: "Test Teacher One",
			Email:       "testteacher@example.com",
		}
		mock.ExpectExec("INSERT INTO teachers").WillReturnResult(sqlmock.NewResult(1, 1))
		err := repo.Save(ctx, mockTeacher)
		assert.NoError(t, err)
	})

	t.Run("Should Get Common Students", func(t *testing.T) {
		emails := []string{"teacherken@gmail.com", "teacherjoe@gmail.com"}
		rows := sqlmock.NewRows([]string{"email"}).
			AddRow("commonstudent1@gmail.com").
			AddRow("commonstudent2@gmail.com")

		mock.ExpectQuery(`
		    SELECT s.email AS email
		    FROM student AS s
		    JOIN registration AS r ON s.id = r.student_id
		    JOIN teacher AS t ON r.teacher_id = t.id
		    WHERE t.email IN (.+)
		    GROUP BY s.email
		    HAVING COUNT\(DISTINCT t.id\) = .+`).
			WithArgs("teacherken@gmail.com", "teacherjoe@gmail.com", len(emails)).
			WillReturnRows(rows)
		res, err := repo.FindCommonStudents(ctx, emails)
		assert.NoError(t, err)
		assert.NotEmpty(t, res)
	})
}
