package controller

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type UserInterface interface {
	decodeJson(r *http.Request)
	getNameAndPass() (string, string)
}

type User struct {
	UserId    int    `db:"user_id" json:"user_id"`
	UserName  string `db:"user_name" json:"user_name" validate:"required,min=2,max=25"`
	UserPass  string `db:"user_pass" json:"user_pass" validate:"passwd, required,min=6,max=50"`
	LastVisit int64  `db:"last_visit" json:"last_visit"`
	Role      string `db:"role" json:"role"`
}

type UserData struct {
	UserName string
	UserPass string
}

type UserAuthRequest struct {
	UserName string `json:"user_name" validate:"required,min=2,max=25"`
	UserPass string `json:"user_pass" validate:"passwd, required,min=6,max=50"`
}

func (uar *UserAuthRequest) decodeJson(r *http.Request) {
	if err := json.NewDecoder(r.Body).Decode(&uar); err != nil {
		log.Error("UserAuthRequest decodeJson: ", err)
	}
}

func (uar *UserAuthRequest) getNameAndPass() (string, string) {
	return uar.UserName, uar.UserPass
}

type UserRegistrationRequest struct {
	UserName  string `json:"user_name" validate:"required,min=2,max=25"`
	UserPass  string `json:"user_pass" validate:"passwd, required,min=6,max=50"`
	MobileKey string `json:"mobile_key" validate:"required,min=6,max=50"`
}

func (urr *UserRegistrationRequest) decodeJson(r *http.Request) {
	if err := json.NewDecoder(r.Body).Decode(&urr); err != nil {
		log.Error("UserRegistrationRequest decodeJson: ", err)
	}
}

func (urr *UserRegistrationRequest) getNameAndPass() (string, string) {
	return urr.UserName, urr.UserPass
}
