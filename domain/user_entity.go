package domain

type UserEntity struct {
	Id       string `bson:"_id" json:"-"`
	UserName string `bson:"user_name" json:"user_name"`
	UserPass string `bson:"user_pass" json:"user_pass"`
	UserRole string `bson:"user_role" json:"-"`
}

func (ue *UserEntity) CheckUserExist(newUser UserEntity) bool {
	if ue.UserName == newUser.UserName {
		return true
	} else {
		return false
	}
}

func (ue *UserEntity) CheckUserNameAndPass(newUser UserEntity) bool {
	if (ue.UserName == newUser.UserName) && (ue.UserPass == newUser.UserPass) {
		return true
	} else {
		return false
	}
}
