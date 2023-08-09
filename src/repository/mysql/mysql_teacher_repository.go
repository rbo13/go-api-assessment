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
		Email: teacher.Email,
	})

	return err
}

func (repo *MySQLTeacherRepository) FindById(ctx context.Context, id int32) (domain.Teacher, error) {
	res, err := repo.queries.GetTeacher(ctx, id)
	if err != nil {
		return domain.Teacher{}, err
	}

	return domain.Teacher{
		ID:          res.ID,
		TeacherName: res.TeacherName.String,
		Email:       res.Email,
	}, nil
}

func (repo *MySQLTeacherRepository) FindByEmail(ctx context.Context, email string) (domain.Teacher, error) {
	res, err := repo.queries.GetTeacherByEmail(ctx, email)
	if err != nil {
		return domain.Teacher{}, err
	}

	return domain.Teacher{
		ID:          res.ID,
		TeacherName: res.TeacherName.String,
		Email:       res.Email,
	}, nil
}

func (repo *MySQLTeacherRepository) FindCommonStudents(ctx context.Context, emails []string) ([]string, error) {
	res, err := repo.queries.GetStudentsByCommonTeacher(ctx, database.GetStudentsByCommonTeacherParams{
		Emails: emails,
		ID:     int32(len(emails)),
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (repo *MySQLTeacherRepository) List(ctx context.Context) (domain.Teachers, error) {
	list, err := repo.queries.ListTeachers(ctx)
	if err != nil {
		return domain.Teachers{}, err
	}

	teachers := domain.Teachers{}

	for _, l := range list {
		t := domain.Teacher{
			ID:          l.ID,
			TeacherName: l.TeacherName.String,
			Email:       l.Email,
		}

		teachers = append(teachers, t)
	}

	return teachers, nil
}
