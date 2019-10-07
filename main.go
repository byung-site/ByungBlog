package main

import (
	"os"
	"os/signal"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	"byung/config"
	controller "byung/controller"
	"byung/log"
	_ "byung/model"
)

//init configure and log
func init() {

	uploadsDir := config.Conf.DataDirectory + "/uploads"
	logDir := config.Conf.DataDirectory + "/logs"

	_, err := os.Stat(uploadsDir)
	if exist := os.IsExist(err); exist == false {
		if err = os.MkdirAll(uploadsDir, os.ModePerm); err != nil {
			log.Error(err)
			os.Exit(-1)
		}
	} else {
	}

	_, err = os.Stat(logDir)
	if exist := os.IsExist(err); exist == false {
		if err = os.MkdirAll(logDir, os.ModePerm); err != nil {
			log.Error(err)
			os.Exit(-1)
		}
	} else {
	}

	log.Record(logDir, 30)
}

// handle ctrl + c signal
func registerSignal() {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	go func() {
		for _ = range signalChan {
			log.Close()
			os.Exit(1)
		}
	}()
}

func main() {
	e := echo.New()

	registerSignal()

	e.Static("/", "build")
	// login route
	e.POST("/login", controller.Login)
	// register route
	e.POST("/register", controller.Register)
	e.GET("/getArticles", controller.GetArticles)
	e.GET("/getTopics", controller.GetTopics)
	e.GET("/getArticlesByTopicID/:id", controller.GetArticlesByTopicID)
	e.GET("/getArticle/:key", controller.GetArticle)
	e.GET("/getNext/:topicId/:key", controller.GetNextArticleByKey)
	e.GET("/getPrevious/:topicId/:key", controller.GetPreviousArticleBykey)
	e.POST("/updateVisit", controller.UpdateVisit)
	//view
	e.GET("/viewArticleImage/:userId/:key/:name", controller.ViewArticleImage)
	e.GET("/viewAvatar/:userId/:name", controller.ViewAvatar)

	// API group
	r := e.Group("/api")
	// Configure middleware with the custom claims type
	cfg := middleware.JWTConfig{
		Claims:      &controller.JWTUserClaims{},
		TokenLookup: "cookie:token",
		SigningKey:  []byte("1qaz@WSX@@@"),
	}
	r.Use(middleware.JWTWithConfig(cfg))
	r.POST("/changeNickname", controller.ChangeNickname)
	r.POST("/changeEmail", controller.ChangeEmail)
	r.POST("/changePassword", controller.ChangePassword)
	r.POST("/changeAvatar/:userid", controller.ChangeAvatar)
	//article
	r.GET("/createKey", controller.CreateArticleKey)
	r.POST("/saveArticle", controller.SaveArticle)
	r.GET("/getPublish", controller.GetPublishArticles)
	r.GET("/getNewest", controller.GetNewestArticle)
	r.GET("/getHottest", controller.GetHottestArticle)
	r.GET("/getArticlesByUserID/:userid", controller.GetArticlesByUserID)
	r.GET("/getPublishArticles/:userid", controller.GetPublishArticles)
	r.POST("/delArticle", controller.DeleteArticle)
	//e.POST("/publish", controller.PublishArticle)
	r.POST("/publishArticle", controller.PublishArticle)
	//topic
	r.GET("/getTopicsByUserID/:userId", controller.GetTopicsByUserID)
	r.POST("/addTopic", controller.AddTopic)
	r.POST("/deleteTopic", controller.DeleteTopic)
	//upload
	r.POST("/uploadArticleImage", controller.UploadArticleImage)
	r.POST("/uploadArticleAttachImage/:userId/:key", controller.UploadArticleAttachImage)

	if config.Conf.Https == false {
		e.Logger.Fatal(e.Start(config.Conf.ListenAddress))
	} else if config.Conf.Https == true {
		e.Logger.Fatal(e.StartTLS(config.Conf.ListenAddress, config.Conf.CertFile, config.Conf.KeyFile))
	} else {
		e.Logger.Fatal("configure error")
	}
}
