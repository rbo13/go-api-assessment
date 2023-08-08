-- name: CreateTeacher :execresult
INSERT INTO teachers (teacher_name, email)
VALUES (?, ?);

-- name: ListTeachers :many
SELECT * FROM teachers
ORDER BY id DESC;

-- name: GetTeacher :one
SELECT * FROM teachers
WHERE id = ? LIMIT 1;

-- name: GetTeacherByEmail :one
SELECT * FROM teachers
WHERE email = ? LIMIT 1;

-- name: DeleteTeacher :exec
DELETE FROM teachers
WHERE id = ?;