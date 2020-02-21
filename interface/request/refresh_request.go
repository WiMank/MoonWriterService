package request

type RefreshTokensRequest struct {
	Refresh Refresh `json:"refresh"`
}

type Refresh struct {
	SessionId    string `json:"session_id"`
	RefreshToken string `json:"refresh_token"`
	MobileKey    string `json:"mobile_key"`
}
