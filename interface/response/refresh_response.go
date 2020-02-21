package response

import log "github.com/sirupsen/logrus"

type RefreshResponse struct {
	AppResponse AppResponse `json:"refresh_response"`
}

type InvalidSession struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Desc    string `json:"desc"`
}

func (is *InvalidSession) PrintLog() {
	log.Info(is.Message)
}

func (is *InvalidSession) GetStatusCode() int {
	return is.Code
}

type TokenExpired struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Desc    string `json:"desc"`
}

func (is *TokenExpired) PrintLog() {
	log.Info(is.Message)
}

func (is *TokenExpired) GetStatusCode() int {
	return is.Code
}
