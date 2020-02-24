package response

import log "github.com/sirupsen/logrus"

type PurchaseResponse struct {
	AppResponse AppResponse `json:"purchase_response"`
}

type RegisterPurchaseResponse struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Desc    string `json:"desc"`
}

func (rpr *RegisterPurchaseResponse) PrintLog() {
	log.Info(rpr.Message)
}

func (rpr *RegisterPurchaseResponse) GetStatusCode() int {
	return rpr.Code
}

type RegisterPurchaseErrorResponse struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Desc    string `json:"desc"`
}

func (rper *RegisterPurchaseErrorResponse) PrintLog() {
	log.Info(rper.Message)
}

func (rper *RegisterPurchaseErrorResponse) GetStatusCode() int {
	return rper.Code
}
