package domain

type SessionEntity struct {
	RefreshToken string `bson:"refresh_token"`
	ExpiresInR   int64  `bson:"expires_in_r"`
	AccessToken  string `bson:"access_token"`
	ExpiresInA   int64  `bson:"expires_in_a"`
	LastVisit    int64  `bson:"last_visit"`
	MobileKey    string `bson:"mobile_key"`
}
