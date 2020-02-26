package request

import (
	"github.com/WiMank/MoonWriterService/domain"
	"github.com/go-playground/validator/v10"
	log "github.com/sirupsen/logrus"
)

type UserRegistrationRequest struct {
	User domain.UserEntity `json:"new_user" validate:"required"`
}

func (req *UserRegistrationRequest) ValidateRequest(validator *validator.Validate) bool {
	err := validator.Struct(req)
	if err != nil {
		log.Errorf("UserRegistrationRequest validate error: ", err)
		return false
	}
	return true
}
