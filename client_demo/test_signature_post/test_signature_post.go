package main

import (
	"encoding/json"
	"fmt"

	"client_demo/signature_api"

	"github.com/mitchellh/mapstructure"
)

//SomethingReq DoSomething-请求包体
type SomethingReq struct {
	Action string `json:"Action" binding:"required"`
}

//SomethingRsp DoSomething-响应包体
type SomethingRsp struct {
	Action string
	Result string
}

type Response struct {
	Code      int64       `json:"Code"`
	Message   string      `json:"Message"`
	Data      interface{} `json:"Data"`
	RequestID string      `json:"RequestID"`
}

func main() {

	// 以获取库类型列表为例，这里开始组json包体
	path := "/do-something"
	request := SomethingReq{
		Action: "ping",
	}
	// 将请求包体序列化成bytes
	bodyBytes, err := json.Marshal(request)
	fmt.Println(bodyBytes)

	// 实际请求中，需替换实际的 SecretId 和 SecretKey，向平台侧申请
	secretId := "SKIDz8krbsJ5yKBZQpn74WFkmLPx3EXAMPLE"
	secretKey := "Gu5t9xGARNpq86cd98joQYCN3EXAMPLE"

	// 发送请求， 需修改为实际的host
	retData, err := signature_api.SendRequest("localhost:7777", path, signature_api.MethodPOST,
		bodyBytes, secretId, secretKey, signature_api.SigMethodHmacSHA256)
	if err != nil {
		fmt.Print("Error.", err)
		return
	}

	fmt.Println("response' data:", string(retData))

	// 打印回包方法1：解包到具体结构体
	var response = &Response{}
	err = json.Unmarshal(retData, response)
	if err != nil {
		fmt.Println("json(response) unmarshal failed:", err)
		return
	}

	// 打印包体信息
	fmt.Println("response.Rsp.RequestID:", response.RequestID)
	if response.Code == 0 {
		// 返回成功时，可解析具体包体
		var rspBody SomethingRsp
		err = mapstructure.Decode(response.Data, &rspBody)
		if err != nil {
			fmt.Println("response.Rsp.Data fecode error:", err.Error())
		}

		fmt.Println("response.Data.Action:", rspBody.Action)
		fmt.Println("response.Data.Result:", rspBody.Result)

	} else {
		// 返回失败时，仅关心返回码及错误信息即可
		fmt.Println("response.Code:", response.Code)
		fmt.Println("response.Message:", response.Message)
	}

	//todo: 打印回包方法2：通用解包并Json
	//var jsonObj interface{}
	//err = json.Unmarshal([]byte(retData), &jsonObj)
	//if err != nil {
	//	fmt.Println("json unmarshal failed:", err)
	//	return
	//}
	//jsonOut, _ := json.MarshalIndent(jsonObj, "", "  ")
	//b2 := append(jsonOut, '\n')
	//os.Stdout.Write(b2)

	return
}
