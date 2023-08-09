CREATE TABLE IF NOT EXISTS notifications (
  id INT AUTO_INCREMENT NOT NULL PRIMARY KEY,
  teacher_id INT,
  notification_text TEXT,
  FOREIGN KEY (teacher_id) REFERENCES teachers(id)
);