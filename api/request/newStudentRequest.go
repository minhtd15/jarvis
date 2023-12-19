package request

type NewStudentRequest struct {
	CourseId    string `json:"courseId"`
	Name        string `json:"name"`
	DOB         string `json:"DOB"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phoneNumber"`
}
