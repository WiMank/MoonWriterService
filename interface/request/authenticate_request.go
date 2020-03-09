package request

import (
	"encoding/json"
	"github.com/WiMank/MoonWriterService/domain"
	"github.com/go-playground/validator/v10"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type AuthenticateUserRequest struct {
	User      domain.UserEntity `json:"user" validate:"required"`
	MobileKey string            `json:"mobile_key" validate:"required"`
}

func (aur *AuthenticateUserRequest) ValidateRequest(validator *validator.Validate) bool {
	err := validator.Struct(aur)
	if err != nil {
		log.Errorf("AuthenticateUserRequest validate error: ", err)
		return false
	}
	return true
}

func (aur *AuthenticateUserRequest) DecodeAuthenticateUserRequest(r *http.Request) {
	if err := json.NewDecoder(r.Body).Decode(&aur); err != nil {
		log.Errorf("Decode DecodeAuthenticateUserRequest request error:\n", err)
	}
}
