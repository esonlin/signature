package common

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code      ErrorCode   `json:"Code"`
	Message   string      `json:"Message"`
	Data      interface{} `json:"Data"`
	RequestID string      `json:"RequestID"`
}

// GeneralResponse 通用返回--正确返回
func GeneralResponse(c *gin.Context, requestId string, data interface{}) {
	var rsp Response
	rsp.Code = Success
	rsp.Message = "success"
	rsp.RequestID = requestId
	log.Println("RequestID:", rsp.RequestID)
	rsp.Data = data
	c.JSON(http.StatusOK, rsp)
}

// ErrorResponse 通用返回--错误返回
func ErrorResponse(c *gin.Context, requestId string, code ErrorCode, message string) {
	var rsp Response
	rsp.Code = code
	rsp.Message = message
	rsp.RequestID = requestId
	log.Println("RequestID:", rsp.RequestID)
	c.JSON(http.StatusOK, rsp)
}
