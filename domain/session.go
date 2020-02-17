package domain

type Session struct {
	RefreshToken string `bson:"refresh_token"`
	AccessToken  string `bson:"access_token"`
	LastVisit    int64  `bson:"last_visit"`
	MobileKey    string `bson:"mobile_key"`
}
