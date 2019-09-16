package controller

import (
	"bytes"
	"net/http"
	"strconv"

	"github.com/labstack/echo"

	"byung/log"
	"byung/model"
)

func AddTopic(c echo.Context) error {
	name := c.FormValue("topic")
	userId := c.FormValue("userId")

	_, err := model.QueryTopicByName(name)
	if err != nil {
		userIdInt, _ := strconv.Atoi(userId)
		topic := model.Topic{
			Name:   name,
			UserId: userIdInt,
		}
		if err := model.SaveTopic(&topic); err != nil {
			log.Error(err)
			return ResponseOk(c, "新建话题失败")
		}

		topic, err := model.QueryTopicByName(name)
		if err != nil {
			log.Error(err)
			return ResponseOk(c, "新建话题失败")
		}

		log.Infof("added \"%s\" topic, id=%d\n", name, topic.ID)
		return ResponseOk(c, topic.ID)
	}
	log.Infof("\"%s\" topic aready exist\n", name)
	return ResponseError(c, "话题已存在")
}

//得到所有话题
func GetTopics(c echo.Context) error {
	topics, err := model.QueryTopics()
	if err != nil {
		log.Error(err)
		return c.String(http.StatusInternalServerError, "查询话题失败")
	}

	queryArticleCountPerTopic(topics)
	return c.JSON(http.StatusOK, topics)
}

//得到指定用户ID的所有话题
func GetTopicsByUserID(c echo.Context) error {
	userId := c.Param("userId")

	userIdInt, _ := strconv.Atoi(userId)
	topics, err := model.QueryTopicsByUserID(userIdInt)
	if err != nil {
		log.Error(err)
		return c.String(http.StatusInternalServerError, "查询话题失败")
	}

	queryArticleCountPerTopic(topics)

	return c.JSON(http.StatusOK, topics)
}

//删除话题
func DeleteTopic(c echo.Context) error {
	topicId := c.FormValue("topicId")

	id, _ := strconv.Atoi(topicId)
	if err := model.DeleteTopicById(id); err != nil {
		return c.String(http.StatusInternalServerError, "删除话题失败")
	}
	return c.String(http.StatusOK, "删除话题成功")
}

func queryArticleCountPerTopic(topics []*model.Topic) {
	var itemCount int
	var count int
	var err error

	for index, _ := range topics {
		var buffer bytes.Buffer

		itemCount, err = model.QueryArticleCountByTopicID(topics[index].ID)
		if err != nil {
			log.Error(err)
			continue
		}
		itemCountStr := strconv.Itoa(itemCount)
		buffer.WriteString(topics[index].Name)
		buffer.WriteString("(")
		buffer.WriteString(itemCountStr)
		buffer.WriteString(")")
		topics[index].Name = buffer.String()

		count += itemCount
	}
}
