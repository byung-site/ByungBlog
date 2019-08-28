package controllers

import (
	"net/http"
	"os"

	"github.com/labstack/echo"
)

func ViewImage(c echo.Context) error {
	key := c.Param("key")
	filename := c.Param("filename")

	file, err := os.Open(UploadDir + "/" + key + "/" + filename)
	if err != nil {
		return c.String(http.StatusInternalServerError, "failure")
	}
	defer file.Close()

	return c.Stream(200, "image/jpeg", file)
}

func GetAvatar(c echo.Context) error {
	userId := c.Param("userId")
	filename := c.Param("filename")

	file, err := os.Open("assets/avatar/" + userId + "/" + filename)
	if err != nil {
		return c.String(http.StatusInternalServerError, "failure")
	}
	defer file.Close()

	return c.Stream(200, "image/jpeg", file)

}
