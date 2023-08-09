package mysql

import (
	"context"
	"database/sql"

	database "github.com/rbo13/go-api-assessment/generated/db"
	"github.com/rbo13/go-api-assessment/src/domain"
)

type MySQLTeacherRepository struct {
	queries *database.Queries
}

func NewTeacherRepository(db *sql.DB) *MySQLTeacherRepository {
	query := database.New(db)

	return &MySQLTeacherRepository{
		queries: query,
	}
}

func (repo *MySQLTeacherRepository) Save(ctx context.Context, teacher domain.Teacher) error {
	_, err := repo.queries.CreateTeacher(ctx, database.CreateTeacherParams{
		TeacherName: sql.NullString{
			String: teacher.TeacherName,
			Valid:  true,
		},
		Email: sql.NullString{
			String: teacher.Email,
			Valid:  true,
		},
	})

	return err
}

func (repo *MySQLTeacherRepository) FindById(context.Context, int32) (domain.Teacher, error) {
	return domain.Teacher{}, nil
}

func (repo *MySQLTeacherRepository) FindByEmail(context.Context, string) (domain.Teacher, error) {
	return domain.Teacher{}, nil
}

func (repo *MySQLTeacherRepository) List(context.Context) (domain.Teachers, error) {
	return domain.Teachers{}, nil
}
