package domain

type SessionEntity struct {
	Id           string `bson:"_id"`
	UserId       string `bson:"user_id"`
	UserName     string `bson:"user_name"`
	UserRole     string `bson:"user_role"`
	RefreshToken string `bson:"refresh_token"`
	AccessToken  string `bson:"access_token"`
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
