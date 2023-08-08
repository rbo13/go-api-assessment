CREATE TABLE IF NOT EXISTS mentions (
  id INT AUTO_INCREMENT NOT NULL PRIMARY KEY,
  notificationID INT,
  studentID INT,
  FOREIGN KEY (notificationID) REFERENCES notifications(id),
  FOREIGN KEY (studentID) REFERENCES students(id)
);