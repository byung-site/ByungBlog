package controllers

import (
	"byung/config"
	"byung/logger"
	"io"
	"os"

	"github.com/labstack/echo"
)

func UploadArticleImage(c echo.Context) error {
	result := "图片上传失败"

	userId := c.FormValue("userId")
	key := c.FormValue("key")
	file, err := c.FormFile("file")
	if err != nil {
		logger.Error(err)
		return ResponseError(c, result)
	}

	if userId == "" {
		logger.Error("user ID can not be empty")
		return ResponseError(c, result)
	}
	if key == "" {
		logger.Error("key can not be empty")
		return ResponseError(c, result)
	}

	src, err := file.Open()
	if err != nil {
		logger.Error(err)
		return ResponseError(c, result)
	}
	defer src.Close()

	updir := config.Conf.DataDirectory + "/uploads/" + userId + "/" + key
	_, err = os.Stat(updir)
	if os.IsNotExist(err) {
		logger.Info("mkdir -p  " + updir)
		if err = os.MkdirAll(updir, os.ModePerm); err != nil {
			logger.Error(err)
			return ResponseError(c, result)
		}
	}
	//destination
	dst, err := os.Create(updir + "/" + file.Filename)
	if err != nil {
		logger.Error(err)
		return ResponseError(c, result)
	}
	defer dst.Close()

	//copy
	if _, err = io.Copy(dst, src); err != nil {
		logger.Error(err)
		return ResponseError(c, result)
	}

	url := userId + "/" + key + "/" + file.Filename
	return ResponseOk(c, url)
}

func UploadArticleAttachImage(c echo.Context) error {
	result := "图片上传失败"

	userId := c.Param("userId")
	key := c.Param("key")
	file, err := c.FormFile("file")
	if err != nil {
		logger.Error(err)
		return ResponseError(c, result)
	}

	if userId == "" {
		logger.Error("user ID can not be empty")
		return ResponseError(c, result)
	}
	if key == "" {
		logger.Error("key can not be empty")
		return ResponseError(c, result)
	}

	src, err := file.Open()
	if err != nil {
		logger.Error(err)
		return ResponseError(c, result)
	}
	defer src.Close()

	updir := config.Conf.DataDirectory + "/uploads/" + userId + "/" + key
	_, err = os.Stat(updir)
	if os.IsNotExist(err) {
		logger.Info("mkdir -p  " + updir)
		if err = os.MkdirAll(updir, os.ModePerm); err != nil {
			logger.Error(err)
			return ResponseError(c, result)
		}
	}
	//destination
	dst, err := os.Create(updir + "/" + file.Filename)
	if err != nil {
		logger.Error(err)
		return ResponseError(c, result)
	}
	defer dst.Close()

	//copy
	if _, err = io.Copy(dst, src); err != nil {
		logger.Error(err)
		return ResponseError(c, result)
	}

	url := userId + "/" + key + "/" + file.Filename
	return ResponseOk(c, url)
}
