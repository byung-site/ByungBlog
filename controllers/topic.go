package controllers

import (
	"byung-cn/byung/models"
	"log"
	"net/http"

	"github.com/labstack/echo"
)

func AddTopic(c echo.Context) error {
	topicName := c.FormValue("topic")

	topic := &models.Topic{
		Name: topicName,
	}
	models.SaveTopic(topic)
	return c.String(http.StatusOK, "ok")
}

func GetTopics(c echo.Context) error {
	topics, err := models.QueryTopics()
	if err != nil {
		log.Println(err)
		return err
	}

	return c.JSON(http.StatusOK, topics)
}

func DeleteTopic(c echo.Context) error {
	return nil
}
