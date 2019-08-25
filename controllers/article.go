package controllers

import (
	"byung-cn/byung/models"
	"net/http"
	"strconv"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	uuid "github.com/satori/go.uuid"
)

//生成key
func CreateArticleKey(c echo.Context) error {
	uuidv4 := uuid.NewV4()
	return c.JSON(http.StatusOK, uuidv4)
}

//保存同时发布文章
func SaveAndPublishArticle(c echo.Context) error {
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
		article.TopicID = topicIdInt
		article.Summary = summary
		a = article
	}

	a.Publish = 1
	if err = models.SaveArticle(&a); err != nil {
		return c.String(http.StatusInternalServerError, "保存失败！")
	}
	return c.String(http.StatusOK, "ok")
}

//保存或更新文章
func SaveArticle(c echo.Context) error {
	key := c.FormValue("key")
	userId := c.FormValue("userId")
	title := c.FormValue("title")
	summary := c.FormValue("summary")
	content := c.FormValue("content")

	if title == "" || content == "" {
		return c.String(http.StatusInternalServerError, "标题或内容不能为空！")
	}

	userIdInt, _ := strconv.Atoi(userId)

	var a models.Article
	article, err := models.QueryArticleByKey(key)

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			a = models.Article{
				UserID:  userIdInt,
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

//得到所有文章
func GetArticles(c echo.Context) error {
	articles, err := models.QueryAllArticles()
	if err != nil {
		return c.String(http.StatusInternalServerError, "查询文章失败！")
	}

	return c.JSON(http.StatusOK, articles)
}

//得到发布的文章
func GetPublishArticles(c echo.Context) error {
	articles, err := models.QueryPublishArticles()
	if err != nil {
		return c.String(http.StatusInternalServerError, "查询文章失败！")
	}

	return c.JSON(http.StatusOK, articles)
}

//按key查询文章
func GetArticle(c echo.Context) error {
	key := c.Param("key")

	article, err := models.QueryArticleByKey(key)
	if err != nil {
		return c.String(http.StatusInternalServerError, "查询文章失败！")
	}
	article.User, err = models.QueryUserById(article.UserID)
	if err != nil {
		return c.String(http.StatusInternalServerError, "查询文章用户失败！")
	}
	return c.JSON(http.StatusOK, article)
}

//得到最热文章
func GetHottestArticle(c echo.Context) error {
	articles, err := models.QueryHottestArticle()
	if err != nil {
		return c.String(http.StatusInternalServerError, "得到最热文章失败！")
	}

	return c.JSON(http.StatusOK, articles)

}

//得到最新文章
func GetNewestArticle(c echo.Context) error {
	articles, err := models.QueryNewestArticle()
	if err != nil {
		return c.String(http.StatusInternalServerError, "得到最新文章失败！")
	}

	return c.JSON(http.StatusOK, articles)

}

//发布文章
func PublishArticle(c echo.Context) error {
	key := c.FormValue("key")

	article, err := models.QueryArticleByKey(key)
	if err != nil {
		return c.String(http.StatusInternalServerError, "发布失败！")
	}
	if article.Publish == 1 {
		return c.String(http.StatusOK, "文章已发布！")
	}

	article.Publish = 1
	if err = models.SaveArticle(&article); err != nil {
		return c.String(http.StatusInternalServerError, "发布失败！")
	}
	return c.String(http.StatusOK, "发布成功！")
}

func GetArticlesByTopicID(c echo.Context) error {
	topicIdStr := c.Param("id")
	topicId, _ := strconv.Atoi(topicIdStr)

	articles, err := models.QueryArticlesByTopicID(uint(topicId))
	if err != nil {
		return c.String(http.StatusInternalServerError, "得到最文章失败！")
	}

	return c.JSON(http.StatusOK, articles)
}

//删除文章
func DeleteArticle(c echo.Context) error {
	key := c.FormValue("key")

	if err := models.DeleteArticleByKey(key); err != nil {
		return c.String(http.StatusInternalServerError, "删除文章失败！")
	}
	return c.String(http.StatusOK, "删除文章成功!")
}
