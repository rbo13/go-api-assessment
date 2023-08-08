CREATE TABLE IF NOT EXISTS notifications (
  id INT AUTO_INCREMENT NOT NULL PRIMARY KEY,
  teacherID INT,
  NotificationText TEXT,
  FOREIGN KEY (teacherID) REFERENCES teachers(id)
);