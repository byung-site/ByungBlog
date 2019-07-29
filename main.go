package main

import (
	controller "byung-cn/byung/controllers"
	_ "byung-cn/byung/models"

	"github.com/labstack/echo"
)

func main() {
	e := echo.New()

	e.Static("/", "views")
	e.POST("/login", controller.Login)
	e.POST("/register", controller.Register)
	e.Logger.Fatal(e.Start(":1323"))
}
