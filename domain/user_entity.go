package domain

type UserEntity struct {
	Id            string `bson:"_id" json:"-"`
	UserName      string `bson:"user_name" json:"user_name" validate:"required,gte=2,lte=25"`
	UserPass      string `bson:"user_pass" json:"user_pass" validate:"required,gte=6,lte=50"`
	UserRole      string `bson:"user_role" json:"-"`
	IsPremiumUser bool   `bson:"is_premium_user" json:"-"`
}

func (ue *UserEntity) CheckUserExist(newUser UserEntity) bool {
	if ue != nil {

		if ue.UserName == newUser.UserName {
			return true
		} else {
			return false
		}

	} else {
		return false
	}
}

func (ue *UserEntity) CheckUserNameAndPass(newUser UserEntity) bool {
	if ue != nil {

		if (ue.UserName == newUser.UserName) && (ue.UserPass == newUser.UserPass) {
			return true
		} else {
			return false
		}

	} else {
		return false
	}
}
