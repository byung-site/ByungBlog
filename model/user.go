package model

import (
	"github.com/jinzhu/gorm"
)

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
	err = db.Where("nickname = ?", nickname).Take(&user).Error
	if err == gorm.ErrRecordNotFound {
		return user, nil
	}
	return user, err
}

func QueryUserByEmail(email string) (user User, err error) {
	err = db.Where("email = ?", email).Take(&user).Error
	if err == gorm.ErrRecordNotFound {
		return user, nil
	}
	return user, err
}

func QueryUserById(id int) (user User, err error) {
	return user, db.Where("id = ?", id).Take(&user).Error
}

func QueryUserCount() (count int, err error) {
	return count, db.Table("users").Count(&count).Error
}

func SaveUser(user *User) error {
	return db.Save(user).Error
}
