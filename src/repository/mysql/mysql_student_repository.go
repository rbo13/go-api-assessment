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

func (repo *MySQLStudentRepository) Save(ctx context.Context, student domain.Student) (domain.Student, error) {
	res, err := repo.queries.CreateStudent(ctx, database.CreateStudentParams{
		StudentEmail: student.StudentEmail,
		Suspended:    student.Suspended,
	})
	if err != nil {
		return domain.Student{}, err
	}

	lastInsertedId, err := res.LastInsertId()
	if err != nil {
		return domain.Student{}, err
	}

	return domain.Student{
		ID: int32(lastInsertedId),
	}, nil
}

func (repo *MySQLStudentRepository) FindByEmail(ctx context.Context, email string) (domain.Student, error) {
	student, err := repo.queries.GetStudentByEmail(ctx, email)
	if err != nil {
		return domain.Student{}, err
	}
	return domain.Student{
		ID:           student.ID,
		StudentName:  student.StudentName.String,
		StudentEmail: student.StudentEmail,
		Suspended:    student.Suspended,
	}, nil
}
