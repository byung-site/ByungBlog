package models

type User struct {
	Model
	Nickname string `gorm:"unique_index"`
	Email    string `gorm:"unique_index""`
	Avatar   string
	Password string `json:"-"`
	Role     int    `gorm:"default:0"` // 0 管理员 1正常用户
}

func QueryUserByEmailAndPassword(email, password string) (user User, err error) {
	return user, db.Where("email = ? and password = ?", email, password).Take(&user).Error
}

func QueryUserByNickname(nickname string) (user User, err error) {
	return user, db.Where("nickname = ?", nickname).Take(&user).Error
}

func QueryUserByEmail(email string) (user User, err error) {
	return user, db.Where("email = ?", email).Take(&user).Error
}

func QueryUserById(id int) (user User, err error) {
	return user, db.Where("id = ?", id).Take(&user).Error
}

func SaveUser(user *User) error {
	return db.Save(user).Error
}
