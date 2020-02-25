package request

type RefreshTokensRequest struct {
	Refresh struct {
		SessionId    string `json:"session_id"`
		RefreshToken string `json:"refresh_token"`
		MobileKey    string `json:"mobile_key"`
	} `json:"refresh"`
}
