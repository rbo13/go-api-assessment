VERSION=v0.0.0
TODAY = $(shell date +"%Y%m%d_%H%M%S")
SHA = $(shell git rev-parse --short HEAD)

LDFLAGS=-ldflags='-X api.version=${VERSION}.build.${TODAY}:${SHA} -s -w -extldflags "-static"'
DATABASE_URL ?= "mysql://root:password@tcp(localhost:3306)/api_db?parseTime=true&loc=Local"

APP=api
OUT=build/teacher-api
SUB=cmd/${APP}/root.go
MAIN=cmd/${APP}/main.go

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -installsuffix -a -tags netgo ${LDFLAGS} -o ${OUT} ${SUB} ${MAIN}

migrate:
	migrate -path db/migrations -database $(DATABASE_URL) up

down:
	migrate -path db/migrations -database $(DATABASE_URL) down

.PHONY: all build migrate test