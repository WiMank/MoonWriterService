package repository

import (
	"encoding/json"
	"github.com/WiMank/AlarmService/domain"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type userRepository struct {
	db *sqlx.DB
}

type UserRepository interface {
	DecodeUser(r *http.Request) domain.User
	InsertUser(user domain.User)
	DeleteUser(user domain.User)
}

func NewUserRepository(db *sqlx.DB) UserRepository {
	return &userRepository{db}
}

func (ur *userRepository) DecodeUser(r *http.Request) domain.User {
	var requestUser domain.User
	if err := json.NewDecoder(r.Body).Decode(&requestUser); err != nil {
		log.Error("Decode User error! ", err)
	}
	return requestUser
}

func (ur *userRepository) InsertUser(user domain.User) {
	insertUserExec := `INSERT INTO "user" (user_name, user_pass, user_role) VALUES ($1, $2, $3)`
	ur.db.MustExec(insertUserExec, user.UserName, user.UserPass, user.UserRole)
	defer ur.CloseDataBase()
}

func (ur *userRepository) DeleteUser(user domain.User) {
	deleteUserExec := `DELETE FROM "user" WHERE user_name=$1 AND user_pass=$2 AND user_role=$3`
	ur.db.MustExec(deleteUserExec, user.UserName, user.UserPass, user.UserRole)
	defer ur.CloseDataBase()
}

func (ur *userRepository) OpenDataBase() {

}

func (ur *userRepository) CloseDataBase() {
	err := ur.db.Close()
	if err != nil {
		log.Errorf("Failed close database! ", err)
	}
}
