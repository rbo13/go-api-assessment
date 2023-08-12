package main

import (
	"context"
	"net/http"
	"os"

	"github.com/rbo13/go-api-assessment/src/api"
	"github.com/rbo13/go-api-assessment/src/db"
	"github.com/rbo13/go-api-assessment/src/logger"
)

func execute(ctx context.Context) error {
	// start app logger
	logger := logger.New("api")

	conn, err := db.CreateNewConnection(&db.Config{
		Ctx:      ctx,
		MaxConns: 16,
		DSN:      os.Getenv("DB_URL"),
	})
	if err != nil {
		logger.Sugar().Fatalf("Cannot start API due to: %v \n", err)
	}
	defer conn.Close()

	api := api.New(ctx, logger, conn)
	server := api.StartServer()

	go func() {
		if err := server.Start(); err != nil && err != http.ErrServerClosed {
			logger.Sugar().Fatalf("Shutting down server due to: %v \n", err)
		}
	}()

	logger.Sugar().Info("Server running on 0.0.0.0:3000")

	<-ctx.Done()

	return server.Stop(ctx)
}
