CREATE TABLE IF NOT EXISTS students (
  id INT AUTO_INCREMENT NOT NULL PRIMARY KEY,
  student_name VARCHAR(100),
  suspended BOOLEAN
);

CREATE INDEX idx_student_name ON students(student_name);