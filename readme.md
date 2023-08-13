# Go API Assessment

This is the official repository of my Go API Assessment.

## Tech stack:
1.  Go - Backend Language.
2.  MySQL - Database Layer.
3.  [Migrate](https://github.com/golang-migrate/migrate) - Database Migration tool.
4.  [sqlc](https://sqlc.dev) - Database Query/Migration generation tool.
5.  Docker/Docker Compose - Container tools.
6.  AWS Lightsail - Remote environment.

## Pre-requisites:
1.  `docker` service is running, `docker-compose` is also installed.

2.  [sqlc](https://sqlc.dev) is installed. We need `sqlc` to generate our SQL related queries to a valid Go types.

The installation instructions below assumes that Go is correctly installed, i.e `$GOPATH/bin` is in your `$PATH`:
```bash
# using go
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

# check version
sqlc version
```
For more installation instructions, please check the official sqlc installation [docs](https://docs.sqlc.dev/en/stable/overview/install.html).

3.  [Golang-Migrate](https://github.com/golang-migrate/migrate) is installed. Needed to run the database migrations.

To install the CLI:
```bash
# Go 1.16+
# Since we are using mysql as our database we need to use the mysql tag.
go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# check if correctly installed
migrate --version
```
> Note: The installed CLI will be installed on your `$GOPATH/bin`, so make sure that said path is available on your `$PATH`.

For more installation instructions, see golang-migrate installation [docs](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate).

## Cloning the project:
```bash
git clone git@github.com:rbo13/go-api-asessment.git
```

## Setup the project:
```bash
$ cd go-api-assessment
$ go mod tidy

# generate db related files via `sqlc`
sqlc generate
```

## Run the project via docker-compose:
```bash
$ docker-compose up --build

# or if you are on the latest docker version, docker compose is already available as a docker sub-command.

$ docker compose up --build
```

## Run the database migration:
```bash
$ export DATABASE_URL = "mysql://root:password@tcp(db)/api_db?parseTime=true&loc=Local"
$ migrate -path db/migrations -database $(DATABASE_URL) up
```

## Check application:
### The app server should now be running on: localhost:3000