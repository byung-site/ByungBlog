package controllers

import (
	"io"
	"net/http"
	"os"

	"github.com/labstack/echo"
)

type UploadResult struct {
	Code    int
	Message string
	Url     string
}

const UploadDir = "assets/uploads/images"

func UploadImage(c echo.Context) error {
	var result UploadResult

	key := c.FormValue("key")
	file, err := c.FormFile("image")
	if err != nil {
		result.Code = -1
		result.Message = "图片上传失败"
		return c.JSON(http.StatusInternalServerError, result)
	}

	src, err := file.Open()
	if err != nil {
		result.Code = -1
		result.Message = "图片上传失败"
		return c.JSON(http.StatusInternalServerError, result)
	}
	defer src.Close()

	_, err = os.Stat(UploadDir + "/" + key)
	if exist := os.IsExist(err); exist == false {
		if err = os.MkdirAll(UploadDir+"/"+key, os.ModePerm); err != nil {
			result.Code = -1
			result.Message = "创建目录失败"
			return c.JSON(http.StatusInternalServerError, result)
		}
	}
	//destination
	dst, err := os.Create(UploadDir + "/" + key + "/" + file.Filename)
	if err != nil {
		result.Code = -1
		result.Message = "图片上传失败"
		return c.JSON(http.StatusInternalServerError, result)
	}
	defer dst.Close()

	//copy
	if _, err = io.Copy(dst, src); err != nil {
		result.Code = -1
		result.Message = "图片上传失败"
		return c.JSON(http.StatusInternalServerError, result)
	}

	result.Code = 0
	result.Message = "图片上传成功"
	result.Url = key + "/" + file.Filename
	return c.JSON(http.StatusOK, result)
}
