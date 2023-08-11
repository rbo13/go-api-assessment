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
		t.Fatalf("Failed to create mock database: %v", err)
	}
	defer db.Close()

	repo := mysql.NewTeacherRepository(db)

	t.Run("Should save a Teacher", func(t *testing.T) {
		teacher := domain.Teacher{
			TeacherName: "Teacher Ken",
			Email:       "teacherken@gmail.com",
		}
		mock.ExpectExec("INSERT INTO teachers").WithArgs(teacher.TeacherName, teacher.Email).WillReturnResult(sqlmock.NewResult(1, 1))

		err = repo.Save(ctx, teacher)
		assert.NoError(t, err, "Error saving teacher")

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err, "Unfulfilled expectations")
	})
}
