package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/rbo13/go-api-assessment/src/api"
	"github.com/rbo13/go-api-assessment/src/db"
	"github.com/rbo13/go-api-assessment/src/logger"
)

func execute(ctx context.Context) error {
	_ = godotenv.Load()

	// start app logger
	logger := logger.New("api")

	conn, err := db.CreateNewConnection(&db.Config{
		Ctx:      ctx,
		MaxConns: 16,
		DSN:      os.Getenv("DB_URL"),
	})
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
