package response

import log "github.com/sirupsen/logrus"

type AppResponse interface {
	PrintLog(err error)
	GetStatusCode() int
}

type UserExistError struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Desc    string `json:"desc"`
}

func (uex *UserExistError) PrintLog(err error) {
	log.Errorf(uex.Message, err)
}

func (uex *UserExistError) GetStatusCode() int {
	return uex.Code
}

type UserFindError struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Desc    string `json:"desc"`
}

func (ufe *UserFindError) PrintLog(err error) {
	log.Errorf(ufe.Message, err)
}

func (ufe *UserFindError) GetStatusCode() int {
	return ufe.Code
}

type UserCreated struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Desc    string `json:"desc"`
}

func (uc *UserCreated) PrintLog(err error) {
	log.Errorf(uc.Message)
}

func (uc *UserCreated) GetStatusCode() int {
	return uc.Code
}

type UserError struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Desc    string `json:"desc"`
}

func (ur *UserError) PrintLog(err error) {
	log.Errorf(ur.Message)
}

func (ur *UserError) GetStatusCode() int {
	return ur.Code
}
