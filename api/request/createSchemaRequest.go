package request

type CreateSchemaRequest struct {
	CityId       int    `json:"cityId"`
	UserFullName string `json:"userFullName"`
	Email        string `json:"email"`
	PhoneNumber  string `json:"phoneNumber"`
	Location     string `json:"location"`
}
