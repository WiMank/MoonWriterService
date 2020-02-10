package controller

import (
	"encoding/json"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
	"net/http"
)

const insertSessionExec = `INSERT INTO sessions (user_name, last_visit, refresh_token, access_token, mobile_key) VALUES ($1, $2, $3, $4, $5)`
const deleteSessionsByNameExec = `DELETE FROM sessions WHERE user_name = $1`
const getUserByNameAndPassSelect = `SELECT * FROM "user" WHERE user_name=$1 AND user_pass=$2`
const sessionsCountSelect = `SELECT COUNT(*) FROM sessions WHERE user_name=$1`
const updateSessionExec = `UPDATE sessions SET access_token=$1, refresh_token=$2, last_visit=$3 WHERE user_name=$4 AND mobile_key=$5`
const checkSessionMkQuery = `SELECT EXISTS (SELECT mobile_key FROM sessions WHERE mobile_key=$1)::bool`

func (ur *UserRequest) getAndCheckExistUser(db *sqlx.DB) (bool, User) {
	dbUser := ur.getUserFromDb(db)
	if (dbUser.UserName == ur.UserName) && (dbUser.UserPass == ur.UserPass) {
		return true, dbUser
	} else {
		return false, User{}
	}
}

func (ur *UserRequest) getUserFromDb(db *sqlx.DB) User {
	var dbUser User
	err := db.QueryRowx(getUserByNameAndPassSelect, ur.UserName, ur.UserPass).StructScan(&dbUser)
	if err != nil {
		log.Error("getUserFromDb: ", err)
	}
	return dbUser
}

func (ur *UserRequest) insertUserFromDb(db *sqlx.DB, urr *UserRequest) {

}

func (ur *UserRequest) decodeUserRequestJson(r *http.Request) {
	if err := json.NewDecoder(r.Body).Decode(&ur); err != nil {
		log.Error("decodeUserRequestJson: ", err)
	}
}

func encodeJson(w http.ResponseWriter, v interface{}) {
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Error("encodeJson: ", err)
	}
}
