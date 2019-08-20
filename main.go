package main

import (
	controller "byung-cn/byung/controllers"
	_ "byung-cn/byung/models"

	"github.com/labstack/echo"
)

func main() {
	e := echo.New()

	e.Static("/", "views")
	//user
	e.POST("/login", controller.Login)
	e.POST("/register", controller.Register)
	//article
	e.POST("/savearticle", controller.SaveArticle)
	e.GET("/getArticles", controller.GetArticles)
	e.GET("/getArticle/:key", controller.GetArticle)
	e.GET("/getNewest", controller.GetNewestArticle)
	e.GET("/getHottest", controller.GetHottestArticle)
	e.GET("/getArticlesByTopicID/:id", controller.GetArticlesByTopicID)
	e.POST("/delArticle", controller.DeleteArticle)
	//topic
	e.GET("/getTopics", controller.GetTopics)
	e.Logger.Fatal(e.Start(":1323"))
}
