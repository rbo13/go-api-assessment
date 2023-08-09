package domain

type Registration struct {
	ID int32

	TeacherID int32
	StudentID int32
}

type RegistrationResponse struct {
	ID       int32
	Teacher  Teacher
	Students Students
}

type Registrations []Registration
