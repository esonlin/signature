package main

import (
	"encoding/json"
	"fmt"
	"os"

	"client_demo/signature_api"
)

func main() {

	// 以调用接口say_hello为例
	path := "/say-hello"

	// 实际请求中，需替换实际的 SecretId 和 SecretKey，向服务提供方申请
	secretId := "SKIDz8krbsJ5yKBZQpn74WFkmLPx3EXAMPLE"
	secretKey := "Gu5t9xGARNpq86cd98joQYCN3EXAMPLE"

	// 发送请求， 需修改为实际的host
	retData, err := signature_api.SendRequest("localhost:7777", path, signature_api.MethodGET,
		nil, secretId, secretKey, signature_api.SigMethodHmacSHA256)
	if err != nil {
		fmt.Print("Error.", err)
		return
	}

	var jsonObj interface{}
	err = json.Unmarshal([]byte(retData), &jsonObj)
	if err != nil {
		fmt.Println("json unmarshal failed:", err)
		return
	}
	jsonOut, _ := json.MarshalIndent(jsonObj, "", "  ")
	b2 := append(jsonOut, '\n')
	os.Stdout.Write(b2)
}
