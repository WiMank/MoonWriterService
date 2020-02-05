package controller

import "net/http"

type User struct {
	UserId    int
	UserName  string
	UserPass  string
	LastVisit int64
	Role      string
}

func Authentication(w http.ResponseWriter, r *http.Request) {

}

func RegisterNewUser(w http.ResponseWriter, r *http.Request) {

}
