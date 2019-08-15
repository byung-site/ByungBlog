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
	e.GET("/getarticles", controller.GetArticles)
	e.GET("/getarticle/:key", controller.GetArticle)
	e.GET("/getnewestarticle", controller.GetNewestArticle)
	e.GET("/gethottestarticle", controller.GetHottestArticle)
	e.POST("/delarticle", controller.DeleteArticle)
	//topic
	e.POST("/addtopic", controller.AddTopic)
	e.POST("/deltopic", controller.DeleteTopic)
	e.GET("/gettopics", controller.GetTopics)
	e.Logger.Fatal(e.Start(":1323"))
}
