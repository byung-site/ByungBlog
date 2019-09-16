package controller

import (
	"net/http"

	"github.com/labstack/echo"
)

type RequestResult struct {
	Code    int         `json:"code"`
	Message interface{} `json:"message"`
}

func ResponseOk(c echo.Context, message interface{}) error {
	ret := RequestResult{
		Code:    0,
		Message: message,
	}
	return c.JSON(http.StatusOK, ret)
}

func ResponseFailure(c echo.Context, message interface{}) error {
	ret := RequestResult{
		Code:    1,
		Message: message,
	}
	return c.JSON(http.StatusOK, ret)
}

func ResponseError(c echo.Context, message interface{}) error {
	ret := RequestResult{
		Code:    2,
		Message: message,
	}
	return c.JSON(http.StatusOK, ret)
}
