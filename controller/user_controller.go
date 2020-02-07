package controller

import (
	"encoding/json"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

func (controller *AuthenticationController) Authentication(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(contentTypeHeader, applicationJsonType)
	db := controller.BaseController.OpenAppDataBase()
	defer controller.BaseController.CloseAppDataBase(db)

	var request UserAuthRequest
	request.decodeJson(r)
	request.authenticateUser(w, db)
}

func (controller *AuthenticationController) RegisterNewUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(contentTypeHeader, applicationJsonType)
	db := controller.BaseController.OpenAppDataBase()
	defer controller.BaseController.CloseAppDataBase(db)

	var authUser UserRegistrationRequest
	authUser.decodeJson(r)
	authUser.registerUser(w, db)
}

func (controller *AuthenticationController) GetUserProfile(w http.ResponseWriter, r *http.Request) {

}

func (controller *AuthenticationController) UpdateUser(w http.ResponseWriter, r *http.Request) {

}

func (controller *AuthenticationController) DeleteUser(w http.ResponseWriter, r *http.Request) {

}

func (uar *UserAuthRequest) authenticateUser(w http.ResponseWriter, db *sqlx.DB) {
	if uar.canAuthenticationUser(db) {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		encodeJson(
			w,
			AuthenticationResponse{
				"User not registered in the system",
				uar.UserName,
				http.StatusText(http.StatusUnauthorized)})
	}
}

func (uar *UserAuthRequest) canAuthenticationUser(db *sqlx.DB) bool {
	dbUser := getUserFromDb(db, uar)
	if (dbUser.UserName == uar.UserName) && (dbUser.UserPass == uar.UserPass) {
		return true
	} else {
		return false
	}
}

func (urr *UserRegistrationRequest) registerUser(w http.ResponseWriter, db *sqlx.DB) {
	if urr.canRegisterUser(db) {
		w.WriteHeader(http.StatusCreated)
		newUser := User{UserName: urr.UserName, UserPass: urr.UserPass}
		newUser.insertUserFromDb(db)
		encodeJson(w,
			AuthenticationResponse{
				"Successful registration",
				urr.UserName,
				http.StatusText(http.StatusCreated)})
	} else {
		w.WriteHeader(http.StatusBadRequest)
		encodeJson(w,
			AuthenticationResponse{
				"A user with this name is already registered",
				urr.UserName,
				http.StatusText(http.StatusBadRequest)})
	}
}

func (urr *UserRegistrationRequest) canRegisterUser(db *sqlx.DB) bool {
	dbUser := getUserFromDb(db, urr)
	if (dbUser.UserName == urr.UserName) && (dbUser.UserPass == urr.UserPass) {
		return false
	} else {
		return true
	}
}

func getUserFromDb(db *sqlx.DB, userInterface UserInterface) UserData {
	var dbUser UserData
	userName, userPass := userInterface.getNameAndPass()
	err := db.QueryRowx(`SELECT * FROM "user" WHERE user_name=$1 AND user_pass=$2`, userName, userPass).StructScan(&dbUser)
	if err != nil {
		log.Error("getUserFromDb: ", err)
	}
	return dbUser
}

func (u *User) insertUserFromDb(db *sqlx.DB) {
	insertUserQuery := `INSERT INTO "user" (user_name, user_pass, last_visit, role) VALUES ($1, $2, $3, $4)`
	db.MustExec(insertUserQuery, u.UserName, u.UserPass, nowAsUnixMilliseconds(), "user")
}

func encodeJson(w http.ResponseWriter, v interface{}) {
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Error("encodePersonJson: ", err)
	}
}

func nowAsUnixMilliseconds() int64 {
	return time.Now().Round(time.Millisecond).UnixNano() / 1e6
}
