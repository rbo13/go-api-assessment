package mysql

import (
	"context"
	"database/sql"

	database "github.com/rbo13/go-api-assessment/generated/db"
	"github.com/rbo13/go-api-assessment/src/domain"
)

type MySQLRegistrationRepository struct {
	queries *database.Queries
}

func NewRegistrationRepository(db *sql.DB) *MySQLRegistrationRepository {
	query := database.New(db)

	return &MySQLRegistrationRepository{
		queries: query,
	}
}

func (repo *MySQLRegistrationRepository) Save(ctx context.Context, reg domain.Registration) error {
	_, err := repo.queries.RegisterStudent(ctx, database.RegisterStudentParams{
		TeacherID: reg.TeacherID,
		StudentID: reg.StudentID,
	})
	if err != nil {
		return err
	}

	return nil
}

func (repo *MySQLRegistrationRepository) FindById(ctx context.Context, id int32) (domain.Registration, error) {
	res, err := repo.queries.GetRegistrationByID(ctx, id)
	if err != nil {
		return domain.Registration{}, err
	}

	return domain.Registration{
		ID:        res.ID,
		TeacherID: res.TeacherID,
		StudentID: res.StudentID,
	}, nil
}
