package domain

type UserEntity struct {
	UserName      string `db:"user_name" json:"user_name" validate:"required,gte=2,lte=25"`
	UserPass      string `db:"user_pass" json:"user_pass" validate:"required,gte=6,lte=50"`
	UserRole      string `db:"user_role" json:"-"`
	IsPremiumUser bool   `db:"premium" json:"-"`
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
