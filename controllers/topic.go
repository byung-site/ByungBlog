package controllers

import (
	"bytes"
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

//得到所有话题
func GetTopics(c echo.Context) error {
	var topicArray []*models.Topic

	topics, err := models.QueryTopics()
	if err != nil {
		return c.String(http.StatusInternalServerError, "查询话题失败")
	}

	var itemCount int
	var count int
	for index, _ := range topics {
		var buffer bytes.Buffer

		itemCount, err = models.QueryArticleCountByTopicID(topics[index].ID)
		if err != nil {
			return c.String(http.StatusInternalServerError, "查询话题失败")
		}
		itemCountStr := strconv.Itoa(itemCount)
		buffer.WriteString(topics[index].Name)
		buffer.WriteString("(")
		buffer.WriteString(itemCountStr)
		buffer.WriteString(")")
		topics[index].Name = buffer.String()

		count += itemCount
	}

	firstTopic := &models.Topic{
		Name: "全部",
	}

	topicArray = append(topicArray, firstTopic)
	topicArray = append(topicArray, topics...)
	return c.JSON(http.StatusOK, topicArray)
}

//删除话题
func DeleteTopic(c echo.Context) error {
	topicId := c.FormValue("topicId")

	id, _ := strconv.Atoi(topicId)
	if err := models.DeleteTopicById(id); err != nil {
		return c.String(http.StatusInternalServerError, "删除话题失败")
	}
	return c.String(http.StatusOK, "删除话题成功")
}
