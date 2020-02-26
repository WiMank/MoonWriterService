package request

import (
	"github.com/go-playground/validator/v10"
	log "github.com/sirupsen/logrus"
)

type RefreshTokensRequest struct {
	Refresh struct {
		SessionId    string `json:"session_id" validate:"required"`
		RefreshToken string `json:"refresh_token" validate:"required,gte=100,lte=200"`
		MobileKey    string `json:"mobile_key" validate:"required"`
	} `json:"refresh"`
}

func (req *RefreshTokensRequest) ValidateRequest(validator *validator.Validate) bool {
	err := validator.Struct(req)
	if err != nil {
		log.Errorf("RefreshTokensRequest validate error: ", err)
		return false
	}
	return true
}
