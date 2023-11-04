package request

type ModifyUserInformationRequest struct {
	Email    string `json:"email"`
	Dob      string `json:"dob"`
	FullName string `json:"fullName"`
	Gender   string `json:"gender"`
}
