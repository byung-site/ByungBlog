package main

import (
	"byung/config"
	controller "byung/controllers"
	"byung/logger"
	_ "byung/models"
	"os"
	"os/signal"

	"github.com/labstack/echo"
)

func init() {

	uploadsDir := config.Conf.DataDirectory + "/uploads"
	logDir := config.Conf.DataDirectory + "/logs"

	_, err := os.Stat(uploadsDir)
	if exist := os.IsExist(err); exist == false {
		if err = os.MkdirAll(uploadsDir, os.ModePerm); err != nil {
			logger.Error(err)
			os.Exit(-1)
		}
	} else {
	}

	_, err = os.Stat(logDir)
	if exist := os.IsExist(err); exist == false {
		if err = os.MkdirAll(logDir, os.ModePerm); err != nil {
			logger.Error(err)
			os.Exit(-1)
		}
	} else {
	}

	logger.Record(logDir, 30)
}

func registerSignal() {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	go func() {
		for _ = range signalChan {
			logger.Close()
			os.Exit(1)
		}
	}()

}

func main() {
	e := echo.New()

	registerSignal()

	e.Static("/", "build")
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
	e.GET("/getArticlesByUserID/:userid", controller.GetArticlesByUserID)
	e.POST("/delArticle", controller.DeleteArticle)
	//e.POST("/publish", controller.PublishArticle)
	e.POST("/saveAndPublish", controller.SaveAndPublishArticle)
	e.POST("/updateVisit", controller.UpdateVisit)
	//topic
	e.GET("/getTopics", controller.GetTopics)
	e.GET("/getTopicsByUserID/:userId", controller.GetTopicsByUserID)
	e.POST("/addTopic", controller.AddTopic)
	//upload
	e.POST("/uploadArticleImage", controller.UploadArticleImage)
	e.POST("/uploadArticleAttachImage/:userId/:key", controller.UploadArticleAttachImage)
	//view
	e.GET("/viewArticleImage/:userId/:key/:name", controller.ViewArticleImage)
	e.GET("/viewAvatar/:userId/:name", controller.ViewAvatar)
	e.Logger.Fatal(e.Start(config.Conf.ListenAddress))
}
