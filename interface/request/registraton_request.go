package request

import "github.com/WiMank/MoonWriterService/domain"

type UserRegistrationRequest struct {
	User domain.UserEntity `json:"new_user"`
}
