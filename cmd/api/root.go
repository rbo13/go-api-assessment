package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/rbo13/go-api-assessment/src/api"
	"github.com/rbo13/go-api-assessment/src/db"
	"github.com/rbo13/go-api-assessment/src/logger"
)

const (
	maxRetryConn = 5
	maxDBConn    = 16
)

func execute(ctx context.Context) error {
	_ = godotenv.Load()

	var conn *sql.DB
	var err error

	// start app logger
	logger := logger.New("api")
	DB_URL := os.Getenv("DB_URL")

	for i := 0; i < maxRetryConn; i++ {
		conn, err = db.CreateNewConnection(&db.Config{
			Ctx:      ctx,
			MaxConns: maxDBConn,
			DSN:      DB_URL,
		})

		if err == nil {
			break
		}

		logger.Sugar().Info("Retrying database connection...")
		time.Sleep(3 * time.Second)
	}

	if err != nil {
		logger.Sugar().Fatalf("Cannot establish database connection due to: %v \n", err)
	}

	defer conn.Close()

	appAddress := fmt.Sprintf("0.0.0.0:%s", os.Getenv("PORT"))

	api := api.New(ctx, logger, conn)
	server := api.StartServer(appAddress)

	go func() {
		if err := server.Start(); err != nil && err != http.ErrServerClosed {
			logger.Sugar().Fatalf("Cannot start server, shutting down due to: %v \n", err)
		}
	}()

	logger.Sugar().Infof("Success! Server is running on %s", appAddress)

	<-ctx.Done()

	return server.Stop(ctx)
}
