CREATE TABLE IF NOT EXISTS registrations (
  id INT AUTO_INCREMENT NOT NULL PRIMARY KEY,
  teacher_id INT NOT NULL,
  student_id INT NOT NULL,
  FOREIGN KEY (teacher_id) REFERENCES teachers(id),
  FOREIGN KEY (student_id) REFERENCES students(id)
);