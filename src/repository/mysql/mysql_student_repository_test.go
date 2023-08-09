package mysql_test

import (
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
		mock.ExpectExec("INSERT INTO students").WillReturnResult(sqlmock.NewResult(1, 1))
		res, err := repo.Save(ctx, mockStudent)
		assert.NotZero(t, res.ID)
		assert.NoError(t, err)
	})
}
