CREATE TABLE IF NOT EXISTS students (
  id INT AUTO_INCREMENT NOT NULL PRIMARY KEY,
  student_name VARCHAR(100),
  student_email VARCHAR(150) NOT NULL,
  suspended BOOLEAN NOT NULL DEFAULT false
);

CREATE INDEX idx_student_name ON students(student_email);