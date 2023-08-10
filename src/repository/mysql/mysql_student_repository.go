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

func (repo *MySQLStudentRepository) Suspend(ctx context.Context, email string) error {
	_, err := repo.queries.SuspendStudent(ctx, email)
	if err != nil {
		return err
	}

	return nil
}

func (repo *MySQLStudentRepository) GetStudentMentionsFromTeacher(ctx context.Context, teacherEmail string, studentEmails []string) (domain.Students, error) {
	results, err := repo.queries.GetMentionsFromTeacher(ctx, database.GetMentionsFromTeacherParams{
		Email:  teacherEmail,
		Emails: studentEmails,
	})
	if err != nil {
		return domain.Students{}, err
	}

	var recipients domain.Students
	for _, result := range results {
		student := domain.Student{
			StudentEmail: result.StudentEmail,
		}

		recipients = append(recipients, student)
	}

	return recipients, nil
}
