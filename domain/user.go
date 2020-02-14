package domain

type User struct {
	UserName string `db:"user_name" json:"user_name" bson:"user_name"`
	UserPass string `db:"user_pass" json:"user_pass" bson:"user_pass"`
	UserRole string `db:"user_role" json:"user_role" bson:"user_role"`
}

type UserResponse struct {
	UserName string `json:"user_name"`
	Message  string `json:"message"`
	Code     int    `json:"code"`
}
