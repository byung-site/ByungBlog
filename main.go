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
	e.POST("/articlesave", controller.SaveArticle)
	e.POST("/topicadd", controller.AddTopic)
	e.GET("/gettopics", controller.GetTopics)
	e.Logger.Fatal(e.Start(":1323"))
}
