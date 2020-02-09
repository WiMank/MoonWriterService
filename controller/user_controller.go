package controller

import (
	"github.com/jmoiron/sqlx"
	"net/http"
)

func (controller *AuthenticationController) Authentication(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(contentTypeHeader, applicationJsonType)
	db := controller.BaseController.OpenAppDataBase()
	defer controller.BaseController.CloseAppDataBase(db)
	var request UserRequest
	request.decodeUserRequestJson(r)
	request.authenticateUser(w, db)
}

func (controller *AuthenticationController) RegisterNewUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(contentTypeHeader, applicationJsonType)
	db := controller.BaseController.OpenAppDataBase()
	defer controller.BaseController.CloseAppDataBase(db)
	var authUser UserRequest
	authUser.decodeUserRequestJson(r)
	authUser.registerUser(w, db)
}

func (controller *AuthenticationController) GetUserProfile(_ http.ResponseWriter, _ *http.Request) {

}

func (controller *AuthenticationController) UpdateUser(_ http.ResponseWriter, _ *http.Request) {

}

func (controller *AuthenticationController) DeleteUser(_ http.ResponseWriter, _ *http.Request) {

}

func (ur *UserRequest) authenticateUser(w http.ResponseWriter, db *sqlx.DB) {
	exist, user := ur.getAndCheckExistUser(db)
	if exist {
		w.WriteHeader(http.StatusOK)
		session := user.createSession(ur.MobileKey, db)
		encodeJson(w, UserResponse{user.UserName, session.AccessToken, session.RefreshToken})
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		encodeJson(w,
			AuthenticationResponse{
				"User not registered in the system",
				ur.UserName,
				http.StatusText(http.StatusUnauthorized),
			})
	}
}

func (ur *UserRequest) registerUser(w http.ResponseWriter, db *sqlx.DB) {
	exist, _ := ur.getAndCheckExistUser(db)
	if !exist {
		w.WriteHeader(http.StatusCreated)
		//TODO: Записать юзера в БД
		encodeJson(w,
			AuthenticationResponse{
				"Successful registration",
				ur.UserName,
				http.StatusText(http.StatusCreated),
			})
	} else {
		w.WriteHeader(http.StatusBadRequest)
		encodeJson(w,
			AuthenticationResponse{
				"A user with this name is already registered",
				ur.UserName,
				http.StatusText(http.StatusBadRequest),
			})
	}
}
