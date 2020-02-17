package domain

type SessionEntity struct {
	UserName     string `bson:"user_name"`
	RefreshToken string `bson:"refresh_token"`
	ExpiresInR   int64  `bson:"expires_in_r"`
	AccessToken  string `bson:"access_token"`
	ExpiresInA   int64  `bson:"expires_in_a"`
	LastVisit    int64  `bson:"last_visit"`
	MobileKey    string `bson:"mobile_key"`
}

type Token struct {
	Tok     string
	Expired int64
}

func (se *SessionEntity) CheckMkExist(newMk string) bool {
	if se.MobileKey == newMk {
		return true
	} else {
		return false
	}
}
