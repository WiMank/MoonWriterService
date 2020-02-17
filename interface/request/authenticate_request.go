package request

import "github.com/WiMank/MoonWriterService/domain"

type AuthenticateUserRequest struct {
	User      domain.UserEntity `json:"user"`
	MobileKey string            `json:"mk"`
}
