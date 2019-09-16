package controller

import (
	"byung/config"
	"byung/log"
	"net/http"
	"os"

	"github.com/labstack/echo"
)

//文章获得图片
func ViewArticleImage(c echo.Context) error {
	userId := c.Param("userId")
	key := c.Param("key")
	filename := c.Param("name")

	url := config.Conf.DataDirectory + "/uploads/" + userId + "/" + key + "/" + filename

	file, err := os.Open(url)
	if err != nil {
		log.Error(err)
		return c.String(http.StatusInternalServerError, "Fail to open image")
	}
	defer file.Close()

	return c.Stream(200, "image/jpeg", file)
}

//获得头像
func ViewAvatar(c echo.Context) error {
	userId := c.Param("userId")
	filename := c.Param("name")

	url := config.Conf.DataDirectory + "/uploads/" + userId + "/avatar/" + filename
	file, err := os.Open(url)
	if err != nil {
		log.Error(err)
		return c.String(http.StatusInternalServerError, "Fail to open image")
	}
	defer file.Close()

	return c.Stream(200, "image/jpeg", file)
}
