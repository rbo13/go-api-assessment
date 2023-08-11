package main

import (
	"context"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rbo13/go-api-assessment/src/api"
	"github.com/rbo13/go-api-assessment/src/http/server"
)

func execute(ctx context.Context) error {
	engine := echo.New()

	server := server.New(
		server.WithHandler(engine),
	)

	api := api.New(ctx, server)

	go func() {
		if err := api.Run(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Shutting down server due to: %v \n", err)
		}
	}()

	log.Println("Server running on 0.0.0.0:3000")

	<-ctx.Done()

	return api.Terminate(ctx)
}
