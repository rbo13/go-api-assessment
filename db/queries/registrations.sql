-- name: RegisterStudent :execresult
INSERT INTO registrations (teacher_id, student_id)
VALUES (?, ?);