package controller

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"

	"byung/config"
	"byung/log"
	"byung/model"
)

type JWTUserClaims struct {
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	Avatar   string `json:"avatar"`
	ID       uint   `json:"id"`
	Role     int    `json:"role"`
	jwt.StandardClaims
}

func Login(c echo.Context) error {
	email := c.FormValue("email")
	password := c.FormValue("password")

	log.Info("logining:", email)
	if ret := verifyEmailFormat(email); ret == false {
		return ResponseFailure(c, "该邮箱格式不正确")
	}

	passHashStr, err := hash256(password)
	if err != nil {
		log.Error(err)
		return ResponseError(c, "内部错误")
	}

	user, err := model.QueryUserByEmailAndPassword(email, passHashStr)
	if err != nil {
		log.Error(err)
		return ResponseFailure(c, "邮箱或密码错误")
	}
	jwtToken, err := getJWTToken(&user)
	if err != nil {
		log.Error(err)
		return ResponseError(c, "内部错误")
	}

	log.Info("login sucess:", email)
	return ResponseOk(c, jwtToken)
}

func Register(c echo.Context) error {
	nickname := c.FormValue("nickname")
	email := c.FormValue("email")
	password := c.FormValue("password")
	confirm := c.FormValue("confirm")

	log.Info("registering:", nickname, "  ", email)
	userCount, err := model.QueryUserCount()
	if err != nil {
		log.Error(err)
		return ResponseError(c, "服务器内部错误")
	}
	if userCount >= config.Conf.MaxUsers {
		log.Error("reach the maximum number of registrations")
		return ResponseFailure(c, "用户数已满")
	}
	u, err := model.QueryUserByNickname(nickname)
	if err != nil {
		log.Error(err)
		return ResponseError(c, "服务器内部错误")
	}
	if u.ID > 0 {
		log.Error("the nickname aready exsits")
		return ResponseFailure(c, "昵称已存在")
	}

	if ret := verifyEmailFormat(email); ret == false {
		log.Error("the email is not correct")
		return ResponseFailure(c, "该邮箱格式不正确")
	}

	u, err = model.QueryUserByEmail(email)
	if err != nil {
		log.Error(err)
		return ResponseError(c, "服务器内部错误")
	}
	if u.ID > 0 {
		log.Error("the email aready exsits")
		return ResponseFailure(c, "邮箱已存在")
	}

	if passLen := len(password); passLen < 8 {
		log.Error("passwrod  to short")
		return ResponseFailure(c, "密码要求8位以上")
	}

	if strings.Compare(password, confirm) != 0 {
		log.Error("two passwords are inconsistent")
		return ResponseFailure(c, "密码不一致")
	}

	hash := sha256.New()
	hash.Write([]byte(password))
	passHashHex := hash.Sum(nil)
	passHashStr := hex.EncodeToString(passHashHex)

	user := &model.User{
		Nickname: nickname,
		Email:    email,
		Password: passHashStr,
		Avatar:   config.Conf.DefaultAvatar,
		Role:     1,
	}

	if err = model.SaveUser(user); err != nil {
		log.Error(err)
		return ResponseError(c, "服务器内部错误")
	}

	*user, err = model.QueryUserByEmail(user.Email)
	if err != nil {
		log.Info("register success but create default avatar failed(", user.ID, " ", nickname, " ", email, ")")
		return ResponseOk(c, "注册成功(创建默认头像失败)")
	}
	err = newDefaultAvatar(user)
	if err != nil {
		log.Info("register success but create default avatar failed(", user.ID, " ", nickname, " ", email, ")")
		return ResponseOk(c, "注册成功(创建默认头像失败)")
	}
	log.Info("register success(", user.ID, " ", nickname, " ", email, ")")
	return ResponseOk(c, "注册成功")
}

func ChangeNickname(c echo.Context) error {
	userId := c.FormValue("userId")
	newNickname := c.FormValue("newNickname")

	if userId == "" {
		log.Error("user ID can not be empty")
		return ResponseFailure(c, "用户ID不能为空")
	}
	if newNickname == "" {
		log.Error("nickname can not be empty")
		return ResponseFailure(c, "昵称不能为空")
	}

	userIdInt, _ := strconv.Atoi(userId)
	user, err := model.QueryUserById(userIdInt)
	if err != nil {
		log.Error(err)
		return ResponseError(c, "内部错误")
	}

	user.Nickname = newNickname

	if err = model.SaveUser(&user); err != nil {
		log.Error(err)
		return ResponseError(c, "内部错误 ")
	}
	jwtToken, err := getJWTToken(&user)
	if err != nil {
		log.Error(err)
		return ResponseError(c, "内部错误")
	}
	log.Info("change nickname successfully")
	return ResponseOk(c, jwtToken)

}

func ChangeEmail(c echo.Context) error {
	userId := c.FormValue("userId")
	newEmail := c.FormValue("newEmail")

	if userId == "" {
		log.Error("user ID can not be empty")
		return ResponseFailure(c, "用户ID不能为空")
	}
	if newEmail == "" {
		log.Error("email can not be empty")
		return ResponseFailure(c, "邮箱不能为空")
	}

	userIdInt, _ := strconv.Atoi(userId)
	user, err := model.QueryUserById(userIdInt)
	if err != nil {
		log.Error(err)
		return ResponseError(c, "内部错误")
	}

	if ret := verifyEmailFormat(newEmail); ret == false {
		log.Error("email format is incorrect")
		return ResponseFailure(c, "该邮箱格式不正确")
	}

	user.Email = newEmail

	if err = model.SaveUser(&user); err != nil {
		log.Error(err)
		return ResponseError(c, "内部错误")
	}
	jwtToken, err := getJWTToken(&user)
	if err != nil {
		log.Error(err)
		return ResponseError(c, "内部错误")
	}
	log.Info("change email successfully")
	return ResponseOk(c, jwtToken)

}

func ChangePassword(c echo.Context) error {
	userId := c.FormValue("userId")
	oldPass := c.FormValue("oldPass")
	newPass := c.FormValue("newPass")
	confirmPass := c.FormValue("confirmPass")

	if userId == "" {
		log.Error("user ID can not be empty")
		return ResponseFailure(c, "用户ID不能为空")
	}
	if oldPass == "" {
		log.Error("old password can not be empty")
		return ResponseFailure(c, "老密码不能为空")
	}
	if newPass == "" {
		log.Error("new password can not be empty")
		return ResponseFailure(c, "新密码不能为空")
	}
	if confirmPass == "" {
		log.Error("comfilrm password can not be empty")
		return ResponseFailure(c, "新密码验证不能为空")
	}
	if len(newPass) < 8 {
		log.Error("password is to short")
		return ResponseFailure(c, "密码长度必须大于8")
	}
	if newPass != confirmPass {
		log.Error("two different inputs")
		return ResponseFailure(c, "两次输入的密码不同")
	}

	userIdInt, _ := strconv.Atoi(userId)
	user, err := model.QueryUserById(userIdInt)
	if err != nil {
		log.Error(err)
		return ResponseFailure(c, "内部错误 ")
	}

	passwordHash, err := hash256(oldPass)
	if err != nil {
		log.Error(err)
		return ResponseError(c, "内部错误")
	}
	if passwordHash != user.Password {
		log.Error("old password is not incorrect")
		return ResponseFailure(c, "老密码不正确")
	}

	user.Password, err = hash256(newPass)
	if err != nil {
		log.Error(err)
		return ResponseError(c, "内部错误")
	}

	if err = model.SaveUser(&user); err != nil {
		log.Error(err)
		return ResponseError(c, "内部错误")
	}

	log.Info("change password successfully")
	return ResponseOk(c, "密码更改成功")
}

func ChangeAvatar(c echo.Context) error {
	result := "头像更改失败"

	userId := c.Param("userid")
	if userId == "" {
		log.Error("user ID can not be empty")
		return ResponseError(c, result)
	}

	file, err := c.FormFile("file")
	if err != nil {
		log.Error(err)
		return ResponseError(c, result)
	}

	src, err := file.Open()
	if err != nil {
		log.Error(err)
		return ResponseError(c, result)
	}
	defer src.Close()

	//destination
	dst, err := os.Create(config.Conf.DataDirectory + "/uploads/" + userId + "/avatar/" + file.Filename)
	if err != nil {
		log.Error(err)
		return ResponseError(c, result)
	}
	defer dst.Close()

	//copy
	if _, err = io.Copy(dst, src); err != nil {
		log.Error(err)
		return ResponseError(c, result)
	}

	userIdInt, _ := strconv.Atoi(userId)
	user, err := model.QueryUserById(userIdInt)
	if err != nil {
		log.Error(err)
		return ResponseError(c, result)
	}

	oldAvatar := user.Avatar
	user.Avatar = file.Filename

	jwtToken, err := getJWTToken(&user)
	if err != nil {
		log.Error(err)
		return ResponseError(c, result)
	}

	err = model.SaveUser(&user)
	if err != nil {
		log.Error(err)
		return ResponseError(c, result)
	}
	if oldAvatar != user.Avatar {
		os.Remove(config.Conf.DataDirectory + "/uploads/" + userId + "/avatar/" + oldAvatar)
	}

	log.Info("change avatar successfully")
	return ResponseOk(c, jwtToken)
}

func verifyEmailFormat(email string) bool {
	pattern := `\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*` //匹配电子邮箱
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(email)
}

func hash256(str string) (string, error) {
	hash := sha256.New()
	if _, err := hash.Write([]byte(str)); err != nil {
		return "", err
	}
	hashHex := hash.Sum(nil)
	hashStr := hex.EncodeToString(hashHex)
	return hashStr, nil
}

func getJWTToken(user *model.User) (string, error) {
	//Set user claims
	claims := &JWTUserClaims{
		user.Nickname,
		user.Email,
		user.Avatar,
		user.ID,
		user.Role,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		},
	}

	//Create token width claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	//Generate encoded token and send it as response.
	return token.SignedString([]byte("1qaz@WSX@@@"))
}

func newDefaultAvatar(user *model.User) error {
	if user.ID == 0 {
		return errors.New("user ID cannot be 0")
	}

	avatarDir := fmt.Sprintf("/uploads/%d/avatar/", user.ID)
	if err := os.MkdirAll(config.Conf.DataDirectory+avatarDir, os.ModePerm); err != nil {
		return err
	}

	avatar := avatarDir + config.Conf.DefaultAvatar
	copyFile(config.Conf.DataDirectory+avatar, config.Conf.Statics+"/"+config.Conf.DefaultAvatar)
	return nil
}

func copyFile(dstName, srcName string) (written int64, err error) {
	src, err := os.Open(srcName)
	if err != nil {
		return
	}
	defer src.Close()
	dst, err := os.OpenFile(dstName, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return
	}
	defer dst.Close()
	return io.Copy(dst, src)
}
