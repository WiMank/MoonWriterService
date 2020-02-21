package request

type RefreshTokensRequest struct {
	SessionId    string `json:"session_id"`
	RefreshToken string `json:"refresh_token"`
	ExpiresInR   int64  `json:"expires_in_r"`
	MobileKey    string `json:"mobile_key"`
}
