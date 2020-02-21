package response

import (
	log "github.com/sirupsen/logrus"
)

type SessionResponse struct {
	AppResponse AppResponse `json:"auth_response"`
}

type TokenResponse struct {
	Message      string `json:"message"`
	Code         int    `json:"code"`
	Desc         string `json:"desc"`
	SessionId    string `json:"session_id"`
	RefreshToken string `json:"refresh_token"`
	AccessToken  string `json:"access_token"`
}

func (tr *TokenResponse) PrintLog() {
	log.Info(tr.Message)
}

func (tr *TokenResponse) GetStatusCode() int {
	return tr.Code
}

type UnauthorizedResponse struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Desc    string `json:"desc"`
}

func (ur *UnauthorizedResponse) PrintLog() {
	log.Info(ur.Message)
}

func (ur *UnauthorizedResponse) GetStatusCode() int {
	return ur.Code
}

type TokenErrorResponse struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Desc    string `json:"desc"`
}

func (ter *TokenErrorResponse) PrintLog() {
	log.Info(ter.Message)
}

func (ter *TokenErrorResponse) GetStatusCode() int {
	return ter.Code
}

type SessionUpdateFailedResponse struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Desc    string `json:"desc"`
}

func (suf *SessionUpdateFailedResponse) PrintLog() {
	log.Info(suf.Message)
}

func (suf *SessionUpdateFailedResponse) GetStatusCode() int {
	return suf.Code
}

type SessionInsertFailedResponse struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Desc    string `json:"desc"`
}

func (sif *SessionInsertFailedResponse) PrintLog() {
	log.Info(sif.Message)
}

func (sif *SessionInsertFailedResponse) GetStatusCode() int {
	return sif.Code
}
