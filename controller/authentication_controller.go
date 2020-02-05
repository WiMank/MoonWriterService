package controller

import (
	"net/http"
)

type User struct {
	UserId    int
	UserName  string
	UserPass  string
	LastVisit int64
	Role      string
}

type AuthenticationController struct {
	BaseController
}

func (authc AuthenticationController) Authentication(w http.ResponseWriter, r *http.Request) {

}

func (authc AuthenticationController) RegisterNewUser(w http.ResponseWriter, r *http.Request) {

}
