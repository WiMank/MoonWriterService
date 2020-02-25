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

type PurchaseTokenExistResponse struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Desc    string `json:"desc"`
}

func (pter *PurchaseTokenExistResponse) PrintLog() {
	log.Info(pter.Message)
}

func (pter *PurchaseTokenExistResponse) GetStatusCode() int {
	return pter.Code
}

type PurchaseValidResponse struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Desc    string `json:"desc"`
}

func (pvr *PurchaseValidResponse) PrintLog() {
	log.Info(pvr.Message)
}

func (pvr *PurchaseValidResponse) GetStatusCode() int {
	return pvr.Code
}

type VerificationPurchaseErrorResponse struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Desc    string `json:"desc"`
}

func (vper *VerificationPurchaseErrorResponse) PrintLog() {
	log.Info(vper.Message)
}

func (vper *VerificationPurchaseErrorResponse) GetStatusCode() int {
	return vper.Code
}

type PurchaseNotFoundResponse struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Desc    string `json:"desc"`
}

func (pnf *PurchaseNotFoundResponse) PrintLog() {
	log.Info(pnf.Message)
}

func (pnf *PurchaseNotFoundResponse) GetStatusCode() int {
	return pnf.Code
}

type CheckPaymentDataErrorResponse struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Desc    string `json:"desc"`
}

func (cpdr *CheckPaymentDataErrorResponse) PrintLog() {
	log.Info(cpdr.Message)
}

func (cpdr *CheckPaymentDataErrorResponse) GetStatusCode() int {
	return cpdr.Code
}

type InsertPurchaseErrorResponse struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Desc    string `json:"desc"`
}

func (iper *InsertPurchaseErrorResponse) PrintLog() {
	log.Info(iper.Message)
}

func (iper *InsertPurchaseErrorResponse) GetStatusCode() int {
	return iper.Code
}
