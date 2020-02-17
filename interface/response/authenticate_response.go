package response

import log "github.com/sirupsen/logrus"

type SessionResponse struct {
	AppResponse AppResponse `json:"session_response"`
}

type TokenResponse struct {
	Message      string `json:"message"`
	Code         int    `json:"code"`
	Desc         string `json:"desc"`
	RefreshToken string `json:"refresh_token"`
	ExpiresInR   int64  `json:"expires_in_r"`
	AccessToken  string `json:"access_token"`
	ExpiresInA   int64  `json:"expires_in_a"`
}

func (tr *TokenResponse) PrintLog(err error) {
	log.Info(tr.Message, err)
}

func (tr *TokenResponse) GetStatusCode() int {
	return tr.Code
}

type UnauthorizedResponse struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Desc    string `json:"desc"`
}

func (ur *UnauthorizedResponse) PrintLog(err error) {
	log.Errorf(ur.Message, err)
}

func (ur *UnauthorizedResponse) GetStatusCode() int {
	return ur.Code
}
