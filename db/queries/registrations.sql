-- name: RegisterStudent :execresult
INSERT INTO registrations (teacher_id, student_id)
VALUES (?, ?);

-- name: GetRegistrationByID :one
SELECT * FROM registrations
WHERE id = ?
LIMIT 1;

-- name: GetStudentsByCommonTeacher :many
SELECT s.student_email AS email
FROM students AS s
JOIN registrations AS r ON s.id = r.student_id 
JOIN teachers AS t ON r.teacher_id = t.id
WHERE t.email IN (sqlc.slice(emails))
GROUP BY s.student_email 
HAVING COUNT(DISTINCT t.id) = ?;