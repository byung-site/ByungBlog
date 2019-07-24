package controller

import (
	"log"
	"net/http"

	"github.com/labstack/echo"
)

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

	log.Println("nickname: ", nickname)
	log.Println("email: ", email)
	log.Println("password: ", password)
	log.Println("repeat: ", repeat)
	return c.String(http.StatusOK, "test")
}
