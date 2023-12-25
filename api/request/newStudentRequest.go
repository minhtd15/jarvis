package request

type NewStudentRequest struct {
	CourseId    string `json:"courseId"`
	Name        string `json:"name"`
	DOB         string `json:"DOB"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phoneNumber"`
}

type ModifyStudentRequest struct {
	StudentId   string `json:"student_id"`
	Name        string `json:"name"`
	DOB         string `json:"DOB"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
}
