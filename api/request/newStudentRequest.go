package request

type NewStudentRequest struct {
	Name        string `json:"name"`
	DOB         string `json:"DOB"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phoneNumber"`
}
