package request

type ForgotPasswordRequest struct {
	Email string `json:"email"`
}

type CheckDigitCode struct {
	Email     string `json:"email"`
	DigitCode int    `json:"digitCode"`
}
