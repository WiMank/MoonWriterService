package domain

type User struct {
	UserId   int    `db:"user_id" `
	UserName string `db:"user_name" json:"user_name"`
	UserPass string `db:"user_pass" json:"user_pass"`
	UserRole string `db:"user_role" json:"user_role"`
}

type UserResponse struct {
	UserName string `json:"user_name"`
	Message  string `json:"message"`
	Code     int    `json:"code"`
}
