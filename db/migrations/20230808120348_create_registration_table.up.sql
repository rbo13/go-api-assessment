CREATE TABLE IF NOT EXISTS registrations (
  id INT AUTO_INCREMENT NOT NULL PRIMARY KEY,
  teacherID INT,
  studentID INT,
  FOREIGN KEY (teacherID) REFERENCES teachers(id),
  FOREIGN KEY (studentID) REFERENCES students(id)
);