package response

import log "github.com/sirupsen/logrus"

type ValidateErrorResponse struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Desc    string `json:"desc"`
}

func (ver *ValidateErrorResponse) PrintLog() {
	log.Info(ver.Message)
}

func (ver *ValidateErrorResponse) GetStatusCode() int {
	return ver.Code
}
