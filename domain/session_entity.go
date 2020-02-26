package domain

type SessionEntity struct {
	Id           string `bson:"_id"`
	UserId       string `bson:"user_id"`
	UserName     string `bson:"user_name"`
	UserRole     string `bson:"user_role"`
	AccessToken  string `bson:"access_token"`
	RefreshToken string `bson:"refresh_token"`
	LastVisit    int64  `bson:"last_visit"`
	MobileKey    string `bson:"mobile_key"`
}
