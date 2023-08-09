-- name: CreateStudent :execresult
INSERT INTO students (student_name, suspended)
VALUES (?, ?);