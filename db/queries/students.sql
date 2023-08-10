-- name: CreateStudent :execresult
INSERT INTO students (student_email, suspended)
VALUES (?, ?);

-- name: GetStudentById :one
SELECT * FROM students
WHERE id = ?
LIMIT 1;

-- name: GetStudentByEmail :one
SELECT * FROM students
WHERE student_email = ?
LIMIT 1;

-- name: SuspendStudent :execresult
UPDATE students
SET suspended = 1
WHERE student_email = ?;

-- name: GetStudentsByTeacherEmail :one
SELECT JSON_OBJECT(
  'teacher', t.email,
  'students', JSON_ARRAYAGG(s.student_email)
) AS result
FROM teachers AS t
JOIN registrations AS r ON t.id = r.teacher_id
JOIN students AS s ON r.student_id = s.id
WHERE t.email = ?
GROUP BY t.email;