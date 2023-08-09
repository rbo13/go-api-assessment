CREATE TABLE IF NOT EXISTS teachers (
  id INT AUTO_INCREMENT NOT NULL PRIMARY KEY,
  teacher_name VARCHAR(100) NOT NULL,
  email VARCHAR(100) NOT NULL
);

CREATE INDEX idx_teacher_email ON teachers(email);
