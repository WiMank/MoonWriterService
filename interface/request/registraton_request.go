package request

import (
	"encoding/json"
	"github.com/WiMank/MoonWriterService/domain"
	"github.com/go-playground/validator/v10"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type UserRegistrationRequest struct {
	User domain.UserEntity `json:"new_user" validate:"required"`
}

func (urr *UserRegistrationRequest) DecodeUserRegistrationRequest(r *http.Request) {
	if err := json.NewDecoder(r.Body).Decode(&urr); err != nil {
		log.Error("DecodeUserRegistrationRequest error! ", err)
	}
}

func (urr *UserRegistrationRequest) ValidateRequest(validator *validator.Validate) bool {
	err := validator.Struct(urr)
	if err != nil {
		log.Errorf("UserRegistrationRequest validate error: ", err)
		return false
	}
	return true
}
