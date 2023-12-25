package request

type RegisterRequest struct {
	Email    string `json:"email"`
	UserName string `json:"user_name"`
	DOB      string `json:"dob"`
	Password string `json:"password"`
	Gender   string `json:"gender"`
	FullName string `json:"fullName"`
}

type NewPasswordRequest struct {
	Email       string `json:"email"`
	NewPassword string `json:"newPassword"`
}

type NewUserAddedByAdmin struct {
	Email       string `json:"email"`
	UserName    string `json:"user_name"`
	DOB         string `json:"dob"`
	Gender      string `json:"gender"`
	FullName    string `json:"fullName"`
	JobPosition string `json:"job_position"`
}
