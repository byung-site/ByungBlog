package controllers

import (
	"net/http"

	"github.com/labstack/echo"
)

type RequestResult struct {
	Ok   bool
	Data interface{}
}

func ResponseOk(c echo.Context, data interface{}) error {
	ret := RequestResult{
		Ok:   true,
		Data: data,
	}
	return c.JSON(http.StatusOK, ret)
}

func ResponseError(c echo.Context, data interface{}) error {
	ret := RequestResult{
		Ok:   false,
		Data: data,
	}
	return c.JSON(http.StatusOK, ret)
}
