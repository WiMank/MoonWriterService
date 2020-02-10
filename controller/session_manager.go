package controller

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
)

func (u *User) createSession(mobileKey string, db *sqlx.DB) Sessions {
	session := Sessions{
		UserName:     u.UserName,
		LastVisit:    nowAsUnixMilliseconds(),
		AccessToken:  u.generateAccessToken(mobileKey),
		RefreshToken: u.generateRefreshToken(mobileKey),
		MobileKey:    mobileKey,
	}

	if session.checkSessionsCount(db) < 5 {
		if session.checkMobileKeyExist(db) {
			session.updateSession(db)
		} else {
			session.insertSession(db)
		}
	} else {
		session.clearSessionsAndInsertLast(db)
	}

	return session
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

func (s *Sessions) updateSession(db *sqlx.DB) {
	db.MustExec(updateSessionExec, s.AccessToken, s.RefreshToken, nowAsUnixMilliseconds(), s.UserName, s.MobileKey)
}

func (s *Sessions) checkMobileKeyExist(db *sqlx.DB) bool {
	var mkExist bool
	err := db.QueryRowx(checkSessionMkQuery, s.MobileKey).Scan(&mkExist)
	if err != nil {
		log.Error("checkMobileKeyExist: ", err)
		return false
	}
	return mkExist
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
