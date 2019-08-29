package controllers

import (
	"byung-cn/byung/models"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

type jwtUserClaims struct {
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

	log.Println(email)
	if ret := verifyEmailFormat(email); ret == false {
		return c.String(http.StatusInternalServerError, "该邮箱格式不正确!")
	}

	passHashStr, err := hash256(password)
	if err != nil {
		errinfo := fmt.Sprintf("%s", err)
		log.Println("password calculate hash256: " + errinfo)
		return c.String(http.StatusInternalServerError, "内部错误!")
	}

	user, err := models.QueryUserByEmailAndPassword(email, passHashStr)
	if err != nil {
		return c.String(http.StatusInternalServerError, "邮箱或密码错误!")
	}
	jwtToken, err := getJWTToken(&user)
	if err != nil {
		errinfo := fmt.Sprintf("%s", err)
		log.Println("generate jwt token: " + errinfo)
		return c.String(http.StatusInternalServerError, "内部错误!")
	}
	return c.String(http.StatusOK, jwtToken)
}

func Register(c echo.Context) error {
	nickname := c.FormValue("nickname")
	email := c.FormValue("email")
	password := c.FormValue("password")
	confirm := c.FormValue("confirm")

	if user, err := models.QueryUserByNickname(nickname); err == nil && user.ID > 0 {
		return c.String(http.StatusInternalServerError, "该昵称已存在!")
	}

	if ret := verifyEmailFormat(email); ret == false {
		return c.String(http.StatusInternalServerError, "该邮箱格式不正确!")
	}

	if user, err := models.QueryUserByEmail(email); err == nil && user.ID > 0 {
		return c.String(http.StatusInternalServerError, "该邮箱已存在!")
	}

	if passLen := len(password); passLen < 8 {
		return c.String(http.StatusInternalServerError, "密码要求8位以上")
	}

	if strings.Compare(password, confirm) != 0 {
		return c.String(http.StatusInternalServerError, "密码不一致!")
	}

	hash := sha256.New()
	hash.Write([]byte(password))
	passHashHex := hash.Sum(nil)
	passHashStr := hex.EncodeToString(passHashHex)

	user := &models.User{
		Nickname: nickname,
		Email:    email,
		Password: passHashStr,
		Avatar:   "defaultavatar.jpeg",
		Role:     1,
	}

	if err := models.SaveUser(user); err != nil {
		return c.String(http.StatusInternalServerError, "用户注册失败!")
	}

	*user, _ = models.QueryUserByEmail(user.Email)
	createDefaultAvatar(user.ID)
	return c.String(http.StatusOK, "注册成功")
}

func ChangeNickname(c echo.Context) error {
	userId := c.FormValue("userId")
	newNickname := c.FormValue("newNickname")

	if userId == "" {
		return c.String(http.StatusInternalServerError, "用户ID不能为空")
	}
	if newNickname == "" {
		return c.String(http.StatusInternalServerError, "昵称不能为空")
	}

	userIdInt, _ := strconv.Atoi(userId)
	user, err := models.QueryUserById(userIdInt)
	if err != nil {
		return c.String(http.StatusInternalServerError, "查询用户失败! ")
	}

	user.Nickname = newNickname

	if err = models.SaveUser(&user); err != nil {
		return c.String(http.StatusInternalServerError, "更新用户失败! ")
	}
	jwtToken, err := getJWTToken(&user)
	if err != nil {
		return c.String(http.StatusInternalServerError, "jtw生成失败!")
	}
	return c.String(http.StatusOK, jwtToken)

}

func ChangeEmail(c echo.Context) error {
	userId := c.FormValue("userId")
	newEmail := c.FormValue("newEmail")

	if userId == "" {
		return c.String(http.StatusInternalServerError, "用户ID不能为空")
	}
	if newEmail == "" {
		return c.String(http.StatusInternalServerError, "邮箱不能为空")
	}

	userIdInt, _ := strconv.Atoi(userId)
	user, err := models.QueryUserById(userIdInt)
	if err != nil {
		return c.String(http.StatusInternalServerError, "查询用户失败! ")
	}

	if ret := verifyEmailFormat(newEmail); ret == false {
		return c.String(http.StatusInternalServerError, "该邮箱格式不正确!")
	}

	user.Email = newEmail

	if err = models.SaveUser(&user); err != nil {
		return c.String(http.StatusInternalServerError, "更新用户失败! ")
	}
	jwtToken, err := getJWTToken(&user)
	if err != nil {
		return c.String(http.StatusInternalServerError, "jtw生成失败!")
	}
	return c.String(http.StatusOK, jwtToken)

}

func ChangePassword(c echo.Context) error {
	userId := c.FormValue("userId")
	oldPass := c.FormValue("oldPass")
	newPass := c.FormValue("newPass")
	confirmPass := c.FormValue("confirmPass")

	if userId == "" {
		return c.String(http.StatusInternalServerError, "用户ID不能为空")
	}
	if oldPass == "" {
		return c.String(http.StatusInternalServerError, "老密码不能为空")
	}
	if newPass == "" {
		return c.String(http.StatusInternalServerError, "新密码不能为空")
	}
	if confirmPass == "" {
		return c.String(http.StatusInternalServerError, "新密码验证不能为空")
	}
	if len(newPass) < 8 {
		return c.String(http.StatusInternalServerError, "密码长度必须大于8")
	}
	if newPass != confirmPass {
		return c.String(http.StatusInternalServerError, "两次输入的密码不同")
	}

	userIdInt, _ := strconv.Atoi(userId)
	user, err := models.QueryUserById(userIdInt)
	if err != nil {
		return c.String(http.StatusInternalServerError, "查询用户失败! ")
	}

	user.Password, err = hash256(newPass)
	if err != nil {
		return c.String(http.StatusInternalServerError, "密码hash256失败")
	}

	if err = models.SaveUser(&user); err != nil {
		return c.String(http.StatusInternalServerError, "更新用户失败! ")
	}

	return c.String(http.StatusOK, "密码更改成功")
}

func ChangeAvatar(c echo.Context) error {
	var result UploadResult

	userId := c.Param("userid")
	file, err := c.FormFile("file")
	if err != nil {
		result.Code = -1
		result.Message = "头像上传失败"
		log.Println("表单没有文件")
		return c.JSON(http.StatusInternalServerError, result)
	}

	src, err := file.Open()
	if err != nil {
		result.Code = -1
		result.Message = "头像上传失败"
		log.Println("源文件打开失败")
		return c.JSON(http.StatusInternalServerError, result)
	}
	defer src.Close()

	//destination
	dst, err := os.Create("assets/avatar/" + userId + "/" + file.Filename)
	if err != nil {
		result.Code = -1
		result.Message = "头像上传失败"
		log.Println("创建目的文件失败")
		return c.JSON(http.StatusInternalServerError, result)
	}
	defer dst.Close()

	//copy
	if _, err = io.Copy(dst, src); err != nil {
		result.Code = -1
		result.Message = "头像上传失败"
		log.Println("文件复制失败")
		return c.JSON(http.StatusInternalServerError, result)
	}

	userIdInt, _ := strconv.Atoi(userId)
	user, err := models.QueryUserById(userIdInt)
	if err != nil {
		result.Code = -1
		result.Message = "头像更改失败"
		log.Println("查询用户失败")
		return c.JSON(http.StatusInternalServerError, result)
	}

	oldAvatar := user.Avatar
	user.Avatar = file.Filename

	jwtToken, err := getJWTToken(&user)
	if err != nil {
		result.Code = -1
		result.Message = "头像更改失败"
		log.Println("生成jwt失败")
		return c.JSON(http.StatusInternalServerError, result)
	}

	err = models.SaveUser(&user)
	if err != nil {
		result.Code = -1
		result.Message = "头像更改失败"
		log.Println("更新用户失败")
		return c.JSON(http.StatusInternalServerError, result)
	}
	if oldAvatar != user.Avatar {
		os.Remove("assets/avatar/" + userId + "/" + oldAvatar)
	}

	result.Code = 0
	result.Message = jwtToken
	result.Url = userId + "/" + file.Filename
	return c.JSON(http.StatusOK, result)
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

func getJWTToken(user *models.User) (string, error) {
	//Set user claims
	claims := &jwtUserClaims{
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

func createDefaultAvatar(userId uint) error {
	if userId == 0 {
		return errors.New("user ID cannot be 0")
	}

	avatarDir := fmt.Sprintf("assets/avatar/%d", userId)
	if err := os.MkdirAll(avatarDir, os.ModePerm); err != nil {
		return err
	}

	copyFile(avatarDir+"/defaultavatar.jpeg", "assets/avatar/defaultavatar.jpeg")
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
