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
}
