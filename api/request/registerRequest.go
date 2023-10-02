package request

type RegisterRequest struct {
	Email    string `json:"email"`
	UserName string `json:"user_name"`
	DOB      string `json:"dob"`
	Password string `json:"password"`
}
