package client

type UserRequest struct {
	Id          int64  `json:"id"` // id = service_name + uuid + date
	PhoneNumber string `json:"phoneNumber"`
	UserName    string `json:"userName"`
	Password    string `json:"password"`
}
