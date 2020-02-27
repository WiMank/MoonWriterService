package domain

type SessionEntity struct {
	Id           string `db:"session_id"`
	UserId       string `db:"user_id"`
	UserName     string `db:"user_name"`
	UserRole     string `db:"user_role"`
	AccessToken  string `db:"access_token"`
	RefreshToken string `db:"refresh_token"`
	LastVisit    int64  `db:"last_visit"`
	MobileKey    string `db:"mobile_key"`
}
