package api_response

type UserDto struct {
	UserId       string `json:"user_id"`
	UserName     string `json:"user_name"`
	Email        string `json:"email"`
	Role         string `json:"role"`
	DOB          string `json:"dob"`
	JobPosition  string `json:"jobPosition"`
	StartingDate string `json:"startingDate"`
}
