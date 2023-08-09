package mysql

import (
	"context"
	"database/sql"

	database "github.com/rbo13/go-api-assessment/generated/db"
	"github.com/rbo13/go-api-assessment/src/domain"
)

type MySQLStudentRepository struct {
	queries *database.Queries
}

func NewStudentRepository(db *sql.DB) *MySQLStudentRepository {
	query := database.New(db)

	return &MySQLStudentRepository{
		queries: query,
	}
}

func (repo *MySQLStudentRepository) Save(ctx context.Context, student domain.Student) error {
	_, err := repo.queries.CreateStudent(ctx, database.CreateStudentParams{
		StudentEmail: student.StudentEmail,
		Suspended:    student.Suspended,
	})

	return err
}
