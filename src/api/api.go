package api

import (
	"context"
	"database/sql"

	"github.com/rbo13/go-api-assessment/src/http/server"
	"github.com/rbo13/go-api-assessment/src/logger"
	"github.com/rbo13/go-api-assessment/src/repository"
	"github.com/rbo13/go-api-assessment/src/repository/mysql"
)

type api struct {
	ctx    context.Context
	logger *logger.Log

	studentRepo      repository.StudentRepository
	teacherRepo      repository.TeacherRepository
	registrationRepo repository.RegistrationRepository
}

func New(ctx context.Context, log *logger.Log, db *sql.DB) *api {
	// initialize repositories
	studentRepo := mysql.NewStudentRepository(db)
	teacherRepo := mysql.NewTeacherRepository(db)
	registrationRepo := mysql.NewRegistrationRepository(db)

	return &api{
		ctx:    ctx,
		logger: log,

		studentRepo:      studentRepo,
		teacherRepo:      teacherRepo,
		registrationRepo: registrationRepo,
	}
}

func (a *api) StartServer(appAddress string) *server.Server {
	if appAddress == "" {
		appAddress = server.DefaultServerAddr
	}

	server := server.New(
		server.WithHandler(a.Handlers()),
		server.WithAddress(appAddress),
	)

	return server
}
