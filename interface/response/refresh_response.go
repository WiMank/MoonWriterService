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

type InvalidToken struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Desc    string `json:"desc"`
}

func (is *InvalidToken) PrintLog() {
	log.Info(is.Message)
}

func (is *InvalidToken) GetStatusCode() int {
	return is.Code
}

type RefreshSessionErrorResponse struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Desc    string `json:"desc"`
}

func (rser *RefreshSessionErrorResponse) PrintLog() {
	log.Info(rser.Message)
}

func (rser *RefreshSessionErrorResponse) GetStatusCode() int {
	return rser.Code
}

type PurchaseUserExistResponse struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Desc    string `json:"desc"`
}

func (puer *PurchaseUserExistResponse) PrintLog() {
	log.Info(puer.Message)
}

func (puer *PurchaseUserExistResponse) GetStatusCode() int {
	return puer.Code
}
