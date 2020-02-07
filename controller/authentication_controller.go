package controller

import (
	"encoding/json"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

const contentTypeHeader = "Content-Type"
const applicationJsonType = "application/json"

type User struct {
	UserId    int    `db:"user_id" json:"user_id"`
	UserName  string `db:"user_name" json:"user_name" validate:"required,min=2,max=25"`
	UserPass  string `db:"user_pass" json:"user_pass" validate:"passwd, required,min=6,max=50"`
	LastVisit int64  `db:"last_visit" json:"last_visit"`
	Role      string `db:"role" json:"role"`
}

type AuthenticationController struct {
	BaseController BaseController
}

type AuthenticationResponse struct {
	Message  string `json:"message"`
	UserName string `json:"user_name"`
	Reason   string `json:"reason"`
}

func (controller *AuthenticationController) Authentication(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(contentTypeHeader, applicationJsonType)
	db := controller.BaseController.OpenAppDataBase()
	defer controller.BaseController.CloseAppDataBase(db)

	var aUser User
	aUser.decodeUserJson(r)

	if aUser.canAuthenticationUser(db) {
		w.WriteHeader(http.StatusOK)
		//TODO: Пускаем и даем токен
		log.Info("Такой юзер есть: ", aUser)
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		encodeUserJson(
			w,
			AuthenticationResponse{
				"User not registered in the system",
				aUser.UserName,
				http.StatusText(http.StatusUnauthorized)})
	}
}

func (controller *AuthenticationController) RegisterNewUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(contentTypeHeader, applicationJsonType)
	db := controller.BaseController.OpenAppDataBase()
	defer controller.BaseController.CloseAppDataBase(db)

	var newUser User
	newUser.decodeUserJson(r)

	if newUser.canRegisterUser(db) {
		w.WriteHeader(http.StatusCreated)
		newUser.insertUserFromDb(db)
		encodeUserJson(w,
			AuthenticationResponse{
				"Successful registration",
				newUser.UserName,
				http.StatusText(http.StatusCreated)})
	} else {
		w.WriteHeader(http.StatusBadRequest)
		encodeUserJson(w,
			AuthenticationResponse{
				"A user with this name is already registered",
				newUser.UserName,
				http.StatusText(http.StatusBadRequest)})
	}
}

func (u *User) canAuthenticationUser(db *sqlx.DB) bool {
	dbUser := u.getUserFromDb(db)
	if (dbUser.UserName == u.UserName) && (dbUser.UserPass == u.UserPass) {
		return true
	} else {
		return false
	}
}

func (u *User) canRegisterUser(db *sqlx.DB) bool {
	dbUser := u.getUserFromDb(db)
	if (dbUser.UserName == u.UserName) && (dbUser.UserPass == u.UserPass) {
		return false
	} else {
		return true
	}
}

func (u *User) getUserFromDb(db *sqlx.DB) User {
	var dbUser User
	err := db.QueryRowx(`SELECT * FROM "user" WHERE user_name=$1 AND user_pass=$2`, u.UserName, u.UserPass).StructScan(&dbUser)
	if err != nil {
		log.Error("getUserFromDb: ", err)
	}
	return dbUser
}

func (u *User) insertUserFromDb(db *sqlx.DB) {
	insertUserQuery := `INSERT INTO "user" (user_name, user_pass, last_visit, role) VALUES ($1, $2, $3, $4)`
	db.MustExec(insertUserQuery, u.UserName, u.UserPass, nowAsUnixMilliseconds(), "user")
}

func encodeUserJson(w http.ResponseWriter, v interface{}) {
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Error("encodePersonJson: ", err)
	}
}

func (u *User) decodeUserJson(r *http.Request) {
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		log.Error("decodePersonJson: ", err)
	}
}

func nowAsUnixMilliseconds() int64 {
	return time.Now().Round(time.Millisecond).UnixNano() / 1e6
}
