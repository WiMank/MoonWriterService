package controller

import (
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

const insertSessionExec = `INSERT INTO sessions (user_name, last_visit, refresh_token, access_token, mobile_key) VALUES ($1, $2, $3, $4, $5)`
const deleteSessionsByNameExec = `DELETE FROM sessions WHERE user_name = $1`
const getUserByNameAndPassSelect = `SELECT * FROM "user" WHERE user_name=$1 AND user_pass=$2`
const sessionsCountSelect = `SELECT COUNT(*) FROM sessions WHERE user_name=$1`

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
	//insertUserQuery := `INSERT INTO "user" (user_name, user_pass, last_visit, role) VALUES ($1, $2, $3, $4)`
	//db.MustExec(insertUserQuery, u.UserName, u.UserPass, nowAsUnixMilliseconds(), "user")
	//log.Info("TOKEN: ", generateToken(urr.MobileKey))
	//insertSessionQuery := `INSERT INTO sessions (user_name, mobile_key, user_token, valid_to) VALUES ($1, $2, $3, $4)`
	//db.MustExec(insertSessionQuery, u.UserName, urr.MobileKey, "token", 123123123)
}

func (u *User) createSession(mobileKey string) Sessions {
	return Sessions{
		UserName:     u.UserName,
		LastVisit:    nowAsUnixMilliseconds(),
		AccessToken:  u.generateAccessToken(mobileKey),
		RefreshToken: u.generateRefreshToken(mobileKey),
		MobileKey:    mobileKey,
	}
}

func (s *Sessions) insertSession(db *sqlx.DB) {
	db.MustExec(insertSessionExec, s.UserName, s.LastVisit, s.RefreshToken, s.AccessToken, s.MobileKey)
}

func (s *Sessions) checkSessionsCount(db *sqlx.DB) int {
	var count int
	err := db.QueryRowx(sessionsCountSelect, s.UserName).Scan(&count)
	if err != nil {
		log.Error("checkSessionsCount: ", err)
		return 0
	}
	return count
}

func (s *Sessions) clearSessionsAndInsertLast(db *sqlx.DB) {
	tx := db.MustBegin()
	tx.MustExec(deleteSessionsByNameExec, s.UserName)
	tx.MustExec(insertSessionExec, s.UserName, s.LastVisit, s.RefreshToken, s.AccessToken, s.MobileKey)
	err := tx.Commit()
	if err != nil {
		log.Error("clearSessionsAndInsertLast: ", err)
		return
	}
	log.Info(fmt.Sprintf(`User %s has more than 5 sessions. Sessions cleared!`, s.UserName))
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

func nowAsUnixMilliseconds() int64 {
	return time.Now().Round(time.Millisecond).UnixNano() / 1e6
}
