package controllers

import (
	"log"
	"net/http"

	"github.com/labstack/echo"
)

func SaveArticle(c echo.Context) error {
	userId := c.FormValue("userId")
	topicId := c.FormValue("topicId")
	title := c.FormValue("title")
	content := c.FormValue("content")

	log.Println(userId)
	log.Println(topicId)
	log.Println(title)
	log.Println(content)
	return c.String(http.StatusOK, "ok")
}

func DeleteArticle(c echo.Context) error {
	return nil
}
