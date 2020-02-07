package controller

type Sessions struct {
	UserId    int    `db:"user_id" json:"user_id"`
	MobileKey string `db:"mobile_key" json:"mobile_key"`
	UserToken string `db:"user_token" json:"user_token"`
	ValidTo   int64  `db:"valid_to" json:"valid_to"`
}

type AuthenticationController struct {
	BaseController BaseController
}

type AuthenticationResponse struct {
	Message  string `json:"message"`
	UserName string `json:"user_name"`
	Reason   string `json:"reason"`
}
