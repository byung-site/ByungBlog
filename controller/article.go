package controller

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	uuid "github.com/satori/go.uuid"

	"byung/config"
	"byung/log"
	"byung/model"
)

//生成key
func CreateArticleKey(c echo.Context) error {
	uuidv4 := uuid.NewV4()
	return ResponseOk(c, uuidv4)
}

//保存同时发布文章
func PublishArticle(c echo.Context) error {
	key := c.FormValue("key")
	userId := c.FormValue("userId")
	topicId := c.FormValue("topicId")
	title := c.FormValue("title")
	summary := c.FormValue("summary")
	content := c.FormValue("content")
	image := c.FormValue("image")

	if title == "" || content == "" {
		log.Error("title or content can not be empty")
		return ResponseFailure(c, "标题或内容不能为空")
	}
	if key == "" {
		log.Error("key can not be empty")
		return ResponseFailure(c, "key不能为空")
	}
	if topicId == "" {
		log.Error("topic ID can not be empty")
		return ResponseFailure(c, "话题不能为空")
	}
	if userId == "0" || userId == "" {
		log.Error("user ID can not be empty")
		return ResponseFailure(c, "用户ID不能0或为空")
	}
	if summary == "" {
		log.Error("summary can not be empty")
		return ResponseFailure(c, "摘要不能为空")
	}

	userIdInt, err := strconv.Atoi(userId)
	if err != nil {
		log.Error(err)
		return ResponseError(c, "内部错误")
	}

	topicIdInt, err := strconv.Atoi(topicId)
	if err != nil {
		log.Error(err)
		return ResponseError(c, "内部错误")
	}

	var a model.Article
	article, err := model.QueryArticleByKey(key)
	if err != nil {
		if err == gorm.ErrRecordNotFound {

			a = model.Article{
				UserID:  userIdInt,
				TopicID: topicIdInt,
				Key:     key,
				Image:   image,
				Title:   title,
				Summary: summary,
				Content: content,
			}
			if image == "" {
				newDefaultAttachImage(&a)
				a.Image = userId + "/" + key + "/" + config.Conf.DefaultArticleAttachImage
			}
		} else {
			log.Error(err)
			return ResponseError(c, "内部错误")
		}
	} else {
		if err == nil {
			count, err := model.QueryArticleCountByTopicID(uint(article.TopicID))
			if err == nil && count == 0 {
				model.DeleteTopicById(topicIdInt)
				log.Info("delete topic: id is ", topicIdInt)
			}
		}
		article.Image = image
		article.Title = title
		article.Content = content
		article.TopicID = topicIdInt
		article.Summary = summary
		a = article
	}

	a.Publish = 1
	if err = model.SaveArticle(&a); err != nil {
		log.Error(err)
		return ResponseError(c, "内部错误")
	}
	log.Info("publish or update article: ", a.Title, " ", a.Key)
	return ResponseOk(c, "发布成功")
}

//保存或更新文章
func SaveArticle(c echo.Context) error {
	key := c.FormValue("key")
	userId := c.FormValue("userId")
	title := c.FormValue("title")
	summary := c.FormValue("summary")
	content := c.FormValue("content")
	image := c.FormValue("image")

	if title == "" || content == "" {
		log.Error("title or content can not be empty")
		return ResponseFailure(c, "标题或内容不能为空")
	}
	if key == "" {
		log.Error("key can not be empty")
		return ResponseFailure(c, "key不能为空")
	}
	if userId == "0" || userId == "" {
		log.Error("user ID can not be empty")
		return ResponseFailure(c, "用户ID不能0或为空")
	}
	if summary == "" {
		log.Error("summary can not be empty")
		return ResponseFailure(c, "摘要不能为空")
	}

	userIdInt, err := strconv.Atoi(userId)
	if err != nil {
		log.Error(err)
		return ResponseError(c, "内部错误")
	}

	var a model.Article
	article, err := model.QueryArticleByKey(key)
	oldAttachImage := article.Image

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			a = model.Article{
				UserID:  userIdInt,
				Key:     key,
				Image:   image,
				Title:   title,
				Summary: summary,
				Content: content,
			}
			if image == "" {
				newDefaultAttachImage(&a)
				a.Image = userId + "/" + key + "/" + config.Conf.DefaultArticleAttachImage
			}
		} else {
			log.Error(err)
			return ResponseError(c, "内部错误")
		}
	} else {
		if image != "" && article.Image != image {
			article.Image = image
			old := config.Conf.DataDirectory + "/uploads/" + oldAttachImage
			os.Remove(old)
		}
		article.Title = title
		article.Summary = summary
		article.Content = content
		a = article
	}

	if err = model.SaveArticle(&a); err != nil {
		log.Error(err)
		return ResponseError(c, "内部错误")
	}
	log.Info("save or update article: ", a.Title, " ", a.Key)
	return ResponseOk(c, "存稿成功")
}

//得到所有文章
func GetArticles(c echo.Context) error {
	articles, err := model.QueryAllArticles()
	if err != nil {
		log.Error(err)
		return ResponseError(c, "内部错误")
	}

	return ResponseOk(c, articles)
}

//得到指定用户ID的所有文章
func GetArticlesByUserID(c echo.Context) error {
	userIdStr := c.Param("userid")
	userId, _ := strconv.Atoi(userIdStr)

	articles, err := model.QueryArticlesByUserID(uint(userId))
	if err != nil {
		return ResponseError(c, "内部错误")
	}

	return ResponseOk(c, articles)
}

//得到发布的文章
func GetPublishArticles(c echo.Context) error {
	userIdStr := c.Param("userid")
	userId, _ := strconv.Atoi(userIdStr)

	articles, err := model.QueryPublishArticles(userId)
	if err != nil {
		return ResponseError(c, "内部错误")
	}

	return ResponseOk(c, articles)
}

//得到同话题的下篇文章
func GetNextArticleByKey(c echo.Context) error {
	topicIdStr := c.Param("topicId")
	key := c.Param("key")

	topicId, _ := strconv.Atoi(topicIdStr)

	articles, err := model.QueryArticlesByTopicID(uint(topicId))
	if err != nil {
		log.Error(err)
		return ResponseError(c, "内部错误")
	}

	var flag bool
	for _, article := range articles {
		if flag == true {
			return ResponseOk(c, article)
		}

		if article.Key == key {
			flag = true
		}
	}

	return ResponseFailure(c, "没有下一篇了")
}

//得到同话题的上一篇文章
func GetPreviousArticleBykey(c echo.Context) error {
	topicIdStr := c.Param("topicId")
	key := c.Param("key")

	topicId, _ := strconv.Atoi(topicIdStr)

	articles, err := model.QueryArticlesByTopicID(uint(topicId))
	if err != nil {
		log.Error(err)
		return ResponseError(c, "内部错误")
	}

	for index, article := range articles {
		if article.Key == key {
			if index-1 >= 0 {
				return ResponseOk(c, articles[index-1])
			}
			break
		}
	}

	return ResponseFailure(c, "没有上一篇了")
}

//按key查询文章
func GetArticle(c echo.Context) error {
	key := c.Param("key")

	article, err := model.QueryArticleByKey(key)
	if err != nil {
		log.Error(err)
		return ResponseError(c, "查询文章失败")
	}

	return ResponseOk(c, article)
}

//得到最热文章
func GetHottestArticle(c echo.Context) error {
	articles, err := model.QueryHottestArticle()
	if err != nil {
		return c.String(http.StatusInternalServerError, "得到最热文章失败！")
	}

	return c.JSON(http.StatusOK, articles)

}

//得到最新文章
func GetNewestArticle(c echo.Context) error {
	articles, err := model.QueryNewestArticle()
	if err != nil {
		return c.String(http.StatusInternalServerError, "得到最新文章失败！")
	}

	return c.JSON(http.StatusOK, articles)

}

/*
//发布文章
func PublishArticle(c echo.Context) error {
	key := c.FormValue("key")

	article, err := model.QueryArticleByKey(key)
	if err != nil {
		return c.String(http.StatusInternalServerError, "发布失败！")
	}
	if article.Publish == 1 {
		return c.String(http.StatusOK, "文章已发布！")
	}

	article.Publish = 1
	if err = model.SaveArticle(&article); err != nil {
		return c.String(http.StatusInternalServerError, "发布失败！")
	}
	return c.String(http.StatusOK, "发布成功！")
}
*/

func GetArticlesByTopicID(c echo.Context) error {
	topicIdStr := c.Param("id")
	topicId, _ := strconv.Atoi(topicIdStr)

	articles, err := model.QueryArticlesByTopicID(uint(topicId))
	if err != nil {
		log.Error(err)
		return ResponseError(c, "内部错误")
	}

	return ResponseOk(c, articles)
}

func UpdateVisit(c echo.Context) error {
	key := c.FormValue("key")

	ret := "更新访问量失败"
	article, err := model.QueryArticleByKey(key)
	if err != nil {
		return ResponseError(c, ret)
	}

	article.Visit++
	err = model.SaveArticle(&article)
	if err != nil {
		return ResponseError(c, ret)
	}

	ret = fmt.Sprintf("更新访问量成功(%d)", article.Visit)
	return ResponseOk(c, ret)
}

//删除文章
func DeleteArticle(c echo.Context) error {
	key := c.FormValue("key")

	if err := model.DeleteArticleByKey(key); err != nil {
		log.Error(err)
		return ResponseError(c, "内部错误")
	}
	log.Info("delete article: ", key)

	return ResponseOk(c, "删除文章成功")
}

func newDefaultAttachImage(article *model.Article) error {
	attachDir := fmt.Sprintf("/uploads/%d/%s/", article.UserID, article.Key)
	if err := os.MkdirAll(config.Conf.DataDirectory+attachDir, os.ModePerm); err != nil {
		return err
	}

	attachImage := attachDir + config.Conf.DefaultArticleAttachImage
	copyFile(config.Conf.DataDirectory+attachImage, config.Conf.Statics+"/"+config.Conf.DefaultArticleAttachImage)
	return nil
}
