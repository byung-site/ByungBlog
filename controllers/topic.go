package controllers

import (
	"byung-cn/byung/models"
	"net/http"
	"strconv"

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
		return c.String(http.StatusInternalServerError, "查询话题失败")
	}

	return c.JSON(http.StatusOK, topics)
}

func DeleteTopic(c echo.Context) error {
	topicId := c.FormValue("topicId")

	id, _ := strconv.Atoi(topicId)
	if err := models.DeleteTopicById(id); err != nil {
		return c.String(http.StatusInternalServerError, "删除话题失败")
	}
	return c.String(http.StatusOK, "删除话题成功")
}
