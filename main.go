package main

import (
	controller "byung-cn/byung/controllers"
	_ "byung-cn/byung/models"

	"github.com/labstack/echo"
)

func main() {
	e := echo.New()

	e.Static("/", "assets")
	//user
	e.POST("/login", controller.Login)
	e.POST("/register", controller.Register)
	e.POST("/changeNickname", controller.ChangeNickname)
	e.POST("/changeEmail", controller.ChangeEmail)
	e.POST("/changePassword", controller.ChangePassword)
	e.POST("/changeAvatar/:userid", controller.ChangeAvatar)
	//article
	e.GET("/createKey", controller.CreateArticleKey)
	e.POST("/saveArticle", controller.SaveArticle)
	e.GET("/getArticles", controller.GetArticles)
	e.GET("/getPublish", controller.GetPublishArticles)
	e.GET("/getArticle/:key", controller.GetArticle)
	e.GET("/getNewest", controller.GetNewestArticle)
	e.GET("/getHottest", controller.GetHottestArticle)
	e.GET("/getArticlesByTopicID/:id", controller.GetArticlesByTopicID)
	e.POST("/delArticle", controller.DeleteArticle)
	e.POST("/publish", controller.PublishArticle)
	e.POST("/saveAndPublish", controller.SaveAndPublishArticle)
	//topic
	e.GET("/getTopics", controller.GetTopics)
	//upload
	e.POST("/uploadImage", controller.UploadImage)
	//view
	e.GET("/view/:key/:filename", controller.ViewImage)
	e.GET("/getAvatar/:userId/:filename", controller.GetAvatar)
	e.Logger.Fatal(e.Start(":5678"))
}
