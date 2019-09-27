package controller

import (
	"strconv"

	"github.com/labstack/echo"

	"byung/log"
	"byung/model"
)

func AddTopic(c echo.Context) error {
	name := c.FormValue("name")
	topicId, _ := strconv.ParseUint(c.FormValue("topicId"), 10, 32)
	userId, _ := strconv.ParseUint(c.FormValue("userId"), 10, 32)

	topic, err := model.QueryTopicByName(name)
	if err == nil && topicId == 0 {
		log.Errorf("\"%s\" topic aready exsit\n", name)
		return ResponseFailure(c, "话题已存在")
	}

	topic.ID = uint(topicId)
	topic.UserID = uint(userId)
	topic.Name = name

	if err := model.SaveTopic(&topic); err != nil {
		log.Error(err)
		return ResponseError(c, "内部错误")
	}

	if topic.ID == 0 {
		log.Infof("add \"%s\" topic\n", name)
	} else {
		log.Infof("update \"%s\" topic, id=%d\n", name, topic.ID)
	}
	return ResponseOk(c, "保存话题成功")
}

//得到所有话题
func GetTopics(c echo.Context) error {
	topics, err := model.QueryTopics()
	if err != nil {
		log.Error(err)
		return ResponseError(c, "查询话题失败")
	}

	queryArticleCountPerTopic(topics)
	return ResponseOk(c, topics)
}

//得到指定用户ID的所有话题
func GetTopicsByUserID(c echo.Context) error {
	userId := c.Param("userId")

	userIdInt, _ := strconv.Atoi(userId)
	topics, err := model.QueryTopicsByUserID(userIdInt)
	if err != nil {
		log.Error(err)
		return ResponseError(c, "内部错误")
	}

	queryArticleCountPerTopic(topics)

	return ResponseOk(c, topics)
}

//删除话题
func DeleteTopic(c echo.Context) error {
	topicId := c.FormValue("topicId")

	id, _ := strconv.Atoi(topicId)
	if err := model.DeleteTopicById(id); err != nil {
		log.Error(err)
		return ResponseError(c, "删除话题失败")
	}
	return ResponseOk(c, "删除话题成功")
}

func queryArticleCountPerTopic(topics []*model.Topic) {
	var itemCount int
	var err error

	for index, _ := range topics {
		itemCount, err = model.QueryArticleCountByTopicID(topics[index].ID)
		if err != nil {
			log.Error(err)
			continue
		}

		topics[index].ArticleNum = itemCount
	}
}
