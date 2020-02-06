package controller

import (
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
	"net/http"
	"regexp"
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

	aUser := decodeUserJson(r)
	if aUser.checkUser(db) {
		w.WriteHeader(http.StatusOK)
		//TODO: Пускаем и даем токен
		log.Info("Такой юзер есть: ", aUser)
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		encodeUserJson(w, AuthenticationResponse{"User not registered in the system", aUser.UserName, http.StatusText(http.StatusUnauthorized)})
	}
}

func (controller *AuthenticationController) RegisterNewUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(contentTypeHeader, applicationJsonType)
	db := controller.BaseController.OpenAppDataBase()
	defer controller.BaseController.CloseAppDataBase(db)

	newUser := decodeUserJson(r)
	checkName(newUser.UserName)

	//TODO: Валидация имени и пароля
	if newUser.checkUser(db) {
		w.WriteHeader(http.StatusBadRequest)
		encodeUserJson(w, AuthenticationResponse{"A user with this name is already registered", newUser.UserName, http.StatusText(http.StatusBadRequest)})
	} else {
		w.WriteHeader(http.StatusCreated)
		insertUser := `INSERT INTO "user" (user_name, user_pass, last_visit, role) VALUES ($1, $2, $3, $4)`
		db.MustExec(insertUser, newUser.UserName, newUser.UserPass, nowAsUnixMilliseconds(), "user")
		encodeUserJson(w, AuthenticationResponse{"Successful registration", newUser.UserName, http.StatusText(http.StatusCreated)})
	}
}

func checkName(name string) {
	d := func(pattern string, text string) {
		matched, _ := regexp.Match(pattern, []byte(text))
		if matched {
			fmt.Println("√", pattern, ":", text)
		} else {
			fmt.Println("X", pattern, ":", text)
		}
	}

	pattern := `/^[a-zA-Z0-9]+([_ -]?[a-zA-Z0-9])*$/`
	d(pattern, name)
}

func (u *User) checkUser(db *sqlx.DB) bool {
	var dbUser User
	err := db.QueryRowx(`SELECT * FROM "user" WHERE user_name=$1 AND user_pass=$2`, u.UserName, u.UserPass).StructScan(&dbUser)
	if err != nil {
		log.Error("CheckUser", err)
	}
	if (dbUser.UserName == u.UserName) && (dbUser.UserPass == u.UserPass) {
		return true
	} else {
		return false
	}
}

func encodeUserJson(w http.ResponseWriter, v interface{}) {
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Error("encodePersonJson()", err)
	}
}

func decodeUserJson(r *http.Request) User {
	var lUser User
	if err := json.NewDecoder(r.Body).Decode(&lUser); err != nil {
		log.Error("decodePersonJson()", err)
	}
	return lUser
}

func nowAsUnixMilliseconds() int64 {
	return time.Now().Round(time.Millisecond).UnixNano() / 1e6
}
