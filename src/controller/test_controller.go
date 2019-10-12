package controller

import (
	"signature/common"

	"fmt"

	"github.com/gin-gonic/gin"
)

//SayHelloRsp SayHello-响应包体
type SayHelloRsp struct {
	Response string
}

func SayHello(c *gin.Context) {

	requestID := common.GenerateRequestID()
	responseData := SayHelloRsp{
		Response: "Hello World!",
	}
	common.GeneralResponse(c, requestID, responseData)
}

//SomethingReq DoSomething-请求包体
type SomethingReq struct {
	Action string `json:"Action" binding:"required"`
}

//SomethingRsp DoSomething-响应包体
type SomethingRsp struct {
	Action string
	Result string
}

func DoSomething(c *gin.Context) {
	var request SomethingReq
	var err error
	requestID := common.GenerateRequestID()
	if err = c.ShouldBindJSON(&request); err != nil {
		fmt.Println("Param Error:", err.Error())
		common.ErrorResponse(c, requestID, common.InvalidParameter, "invalid param")
		return
	}

	fmt.Printf("Request: %#v\n", request)

	var response SomethingRsp
	response.Action = request.Action

	if request.Action == "ping" {
		response.Result = "pong"
	} else {
		response.Result = "Unknown Action " + request.Action
	}

	common.GeneralResponse(c, requestID, response)
}
