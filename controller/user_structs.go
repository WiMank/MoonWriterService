package controller

type User struct {
	UserId   int    `db:"user_id" json:"user_id"`
	UserName string `db:"user_name" json:"user_name" validate:"required,min=2,max=25"`
	UserPass string `db:"user_pass" json:"user_pass" validate:"passwd, required,min=6,max=50"`
	Role     string `db:"role" json:"role"`
}

type UserNameAndPass struct {
	UserName string `db:"user_name"`
	UserPass string `db:"user_pass"`
}

type UserRequest struct {
	UserName  string `json:"user_name" validate:"required,min=2,max=25"`
	UserPass  string `json:"user_pass" validate:"passwd, required,min=6,max=50"`
	MobileKey string `json:"mobile_key" validate:"required,min=6,max=50"`
}

type UserResponse struct {
	UserName     string `json:"user_name"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type Sessions struct {
	SessionId    int    `db:"session_id" json:"session_id"`
	UserName     string `db:"user_name" json:"user_name"`
	RefreshToken string `db:"refresh_token" json:"refresh_token"`
	AccessToken  string `db:"access_token" json:"access_token"`
	MobileKey    string `db:"mobile_key" json:"mobile_key"`
	LastVisit    int64  `db:"last_visit" json:"last_visit"`
}

type AuthenticationController struct {
	BaseController BaseController
}

type AuthenticationResponse struct {
	Message  string `json:"message"`
	UserName string `json:"user_name"`
	Reason   string `json:"reason"`
}
