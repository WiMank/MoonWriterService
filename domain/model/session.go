package model

type Session struct {
	SessionId    int    `db:"session_id"`
	RefreshToken string `db:"refresh_token"`
	AccessToken  string `db:"access_token"`
	LastVisit    string `db:"last_visit"`
	MobileKey    string `db:"mobile_key"`
}
