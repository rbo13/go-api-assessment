DATABASE_URL ?= "mysql://root:password@tcp(localhost:3306)/api_db?parseTime=true&loc=Local"

migrate:
	migrate -path db/migrations -database $(DATABASE_URL) up

down:
	migrate -path db/migrations -database $(DATABASE_URL) down

.PHONY: all build migrate test