package model

type User struct {
	UserId    int    `db:"user_id"`
	UserName  string `db:"user_name"`
	UserPass  string `db:"user_pass"`
	UserRole  string `db:"user_role"`
	MobileKey string `db:"mobile_key"`
}
