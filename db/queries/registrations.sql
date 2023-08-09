-- name: RegisterStudent :execresult
INSERT INTO registrations (teacherID, studentID)
VALUES (?, ?);