package request

type PurchaseRegisterRequest struct {
	AccessToken   string `json:"access_token"`
	PurchaseToken string `json:"purchase_token"`
	OrderId       string `json:"order_id"`
	PurchaseTime  int64  `json:"purchase_time"`
	Sku           string `json:"sku"`
}
