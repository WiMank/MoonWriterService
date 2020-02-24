package domain

type Purchase struct {
	UserId        string `bson:"user_id"`
	IsPremiumUser bool   `bson:"is_premium_user"`
	PurchaseToken string `bson:"purchase_token"`
	OrderId       string `bson:"order_id"`
	PurchaseTime  int64  `bson:"purchase_time"`
	Sku           string `bson:"sku"`
	AccessToken   string `bson:"access_token"`
}
