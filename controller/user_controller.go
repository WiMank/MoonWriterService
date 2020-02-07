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

func (controller *AuthenticationController) GetUserProfile(_ http.ResponseWriter, _ *http.Request) {

}

func (controller *AuthenticationController) UpdateUser(_ http.ResponseWriter, _ *http.Request) {

}

func (controller *AuthenticationController) DeleteUser(_ http.ResponseWriter, _ *http.Request) {

}

func (uar *UserAuthRequest) authenticateUser(w http.ResponseWriter, db *sqlx.DB) {
	if userExist(db, uar) {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		encodeJson(w, AuthenticationResponse{
			"User not registered in the system",
			uar.UserName,
			http.StatusText(http.StatusUnauthorized)})
	}
}

func (urr *UserRegistrationRequest) registerUser(w http.ResponseWriter, db *sqlx.DB) {
	if !userExist(db, urr) {
		log.Info("userExist if")
		w.WriteHeader(http.StatusCreated)
		newUser := User{UserName: urr.UserName, UserPass: urr.UserPass}
		newUser.insertUserFromDb(db)
		encodeJson(w, AuthenticationResponse{
			"Successful registration",
			urr.UserName,
			http.StatusText(http.StatusCreated)})
	} else {
		log.Info("userExist else")
		w.WriteHeader(http.StatusBadRequest)
		encodeJson(w, AuthenticationResponse{
			"A user with this name is already registered",
			urr.UserName,
			http.StatusText(http.StatusBadRequest)})
	}
}

func userExist(db *sqlx.DB, userInterface UserRequestInterface) bool {
	dbUser := getUserFromDb(db, userInterface)
	if (dbUser.UserName == userInterface.getName()) && (dbUser.UserPass == userInterface.getPass()) {
		return true
	} else {
		return false
	}
}

func getUserFromDb(db *sqlx.DB, userInterface UserRequestInterface) UserNameAndPass {
	var dbUserNameAndPass UserNameAndPass
	err := db.QueryRowx(`SELECT user_name, user_pass FROM "user" WHERE user_name=$1 AND user_pass=$2`,
		userInterface.getName(),
		userInterface.getPass()).StructScan(&dbUserNameAndPass)
	if err != nil {
		log.Error("getUserFromDb: ", err)
	}
	return dbUserNameAndPass
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
