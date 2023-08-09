CREATE TABLE IF NOT EXISTS mentions (
  id INT AUTO_INCREMENT NOT NULL PRIMARY KEY,
  notification_id INT,
  student_id INT,
  FOREIGN KEY (notification_id) REFERENCES notifications(id),
  FOREIGN KEY (student_id) REFERENCES students(id)
);