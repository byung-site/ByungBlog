package controllers

import (
	"byung-cn/byung/models"
	"net/http"
	"strconv"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
)

func SaveArticle(c echo.Context) error {
	key := c.FormValue("key")
	userId := c.FormValue("userId")
	topicId := c.FormValue("topicId")
	title := c.FormValue("title")
	summary := c.FormValue("summary")
	content := c.FormValue("content")

	if title == "" || content == "" {
		return c.String(http.StatusInternalServerError, "标题或内容不能为空！")
	}

	userIdInt, _ := strconv.Atoi(userId)
	topicIdInt, _ := strconv.Atoi(topicId)

	var a models.Article
	article, err := models.QueryArticleByKey(key)

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			a = models.Article{
				UserID:  userIdInt,
				TopicID: topicIdInt,
				Key:     key,
				Title:   title,
				Summary: summary,
				Content: content,
			}
		} else {
			return c.String(http.StatusInternalServerError, "保存失败！")
		}
	} else {
		article.Title = title
		article.Content = content
		a = article
	}

	if err = models.SaveArticle(&a); err != nil {
		return c.String(http.StatusInternalServerError, "保存失败！")
	}
	return c.String(http.StatusOK, "ok")
}

func DeleteArticle(c echo.Context) error {
	return nil
}
