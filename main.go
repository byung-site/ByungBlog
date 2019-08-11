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
	e.POST("/savearticle", controller.SaveArticle)
	e.POST("/addtopic", controller.AddTopic)
	e.POST("/deltopic", controller.DeleteTopic)
	e.GET("/gettopics", controller.GetTopics)
	e.Logger.Fatal(e.Start(":1323"))
}
