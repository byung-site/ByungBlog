package controller

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"byung.cn/blog-byung/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

type jwtUserClaims struct {
	Nickname string `json:"nickname"`
	ID       uint   `json:"id"`
	Role     int    `json:"role"`
	jwt.StandardClaims
}

func Login(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	log.Println("username: ", username)
	log.Println("password: ", password)
	return c.String(http.StatusOK, "test")
}

func Register(c echo.Context) error {
	nickname := c.FormValue("nickname")
	email := c.FormValue("email")
	password := c.FormValue("password")
	repeat := c.FormValue("repeat")

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
	if strings.Compare(password, repeat) != 0 {
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
		Avatar:   "/static/images/info-image.png",
		Role:     1,
	}
	if err := models.SaveUser(user); err != nil {
		errinfo := fmt.Sprintf("%s", err)
		log.Println("user register: " + errinfo)
		return c.String(http.StatusInternalServerError, "用户注册失败!")
	}

	jwtToken, err := getJWTToken(user)
	if err != nil {
		errinfo := fmt.Sprintf("%s", err)
		log.Println("generate jwt token: " + errinfo)
		return c.String(http.StatusInternalServerError, "内部错误!")
	}
	return c.String(http.StatusOK, jwtToken)
}

func verifyEmailFormat(email string) bool {
	pattern := `\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*` //匹配电子邮箱
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(email)
}

func getJWTToken(user *models.User) (string, error) {
	//Set user claims
	claims := &jwtUserClaims{
		user.Nickname,
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
