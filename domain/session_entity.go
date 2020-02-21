package domain

type SessionEntity struct {
	UserId       string `bson:"user_id"`
	UserName     string `bson:"user_name"`
	UserRole     string `bson:"user_role"`
	RefreshToken string `bson:"refresh_token"`
	ExpiresInR   int64  `bson:"expires_in_r"`
	AccessToken  string `bson:"access_token"`
	ExpiresInA   int64  `bson:"expires_in_a"`
	LastVisit    int64  `bson:"last_visit"`
	MobileKey    string `bson:"mobile_key"`
}

func (se *SessionEntity) CheckMkExist(newMk string) bool {
	if se.MobileKey == newMk {
		return true
	} else {
		return false
	}
}
