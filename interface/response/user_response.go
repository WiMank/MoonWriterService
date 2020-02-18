package response

import log "github.com/sirupsen/logrus"

type UserResponse struct {
	AppResponse AppResponse `json:"user_response"`
}

type UserCreatedResponse struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Desc    string `json:"desc"`
}

func (uc *UserCreatedResponse) PrintLog() {
	log.Info(uc.Message)
}

func (uc *UserCreatedResponse) GetStatusCode() int {
	return uc.Code
}

type UserExistResponse struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Desc    string `json:"desc"`
}

func (uex *UserExistResponse) PrintLog() {
	log.Info(uex.Message)
}

func (uex *UserExistResponse) GetStatusCode() int {
	return uex.Code
}

type UserFindResponse struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Desc    string `json:"desc"`
}

func (ufe *UserFindResponse) PrintLog() {
	log.Info(ufe.Message)
}

func (ufe *UserFindResponse) GetStatusCode() int {
	return ufe.Code
}

type UserInsertErrorResponse struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Desc    string `json:"desc"`
}

func (ur *UserInsertErrorResponse) PrintLog() {
	log.Info(ur.Message)
}

func (ur *UserInsertErrorResponse) GetStatusCode() int {
	return ur.Code
}
