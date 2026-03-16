package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func OK(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{Code: 0, Message: "success", Data: data})
}

func OKMsg(c *gin.Context, msg string) {
	c.JSON(http.StatusOK, Response{Code: 0, Message: msg})
}

func FailMsg(c *gin.Context, msg string) {
	c.JSON(http.StatusOK, Response{Code: 500, Message: msg})
}

func Fail(c *gin.Context, code int, msg string) {
	c.JSON(http.StatusOK, Response{Code: code, Message: msg})
}

func Error(c *gin.Context, httpCode int, msg string) {
	c.JSON(httpCode, Response{Code: -1, Message: msg})
}

type PageResult struct {
	List  interface{} `json:"list"`
	Total int64       `json:"total"`
	Page  int         `json:"page"`
	Size  int         `json:"size"`
}

func OKPage(c *gin.Context, list interface{}, total int64, page, size int) {
	OK(c, PageResult{List: list, Total: total, Page: page, Size: size})
}
