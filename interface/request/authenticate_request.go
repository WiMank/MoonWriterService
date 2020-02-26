package request

import (
	"github.com/WiMank/MoonWriterService/domain"
	"github.com/go-playground/validator/v10"
	log "github.com/sirupsen/logrus"
)

type AuthenticateUserRequest struct {
	User      domain.UserEntity `json:"user" validate:"required"`
	MobileKey string            `json:"mobile_key" validate:"required"`
}

func (s *AuthenticateUserRequest) ValidateRequest(validator *validator.Validate) bool {
	err := validator.Struct(s)
	if err != nil {
		log.Errorf("AuthenticateUserRequest validate error: ", err)
		return false
	}
	return true
}
