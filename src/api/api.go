package api

import (
	"context"

	"github.com/rbo13/go-api-assessment/src/http/server"
)

type api struct {
	ctx    context.Context
	server *server.Server
}

func New(ctx context.Context, server *server.Server) *api {
	return &api{
		ctx:    ctx,
		server: server,
	}
}

func (a *api) Run() error {
	return a.server.Start()
}

func (a *api) Terminate(ctx context.Context) error {
	return a.server.Stop(ctx)
}
