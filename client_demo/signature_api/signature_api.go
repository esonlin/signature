package signature_api

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var debug bool = true

var (
	// HTTP请求方法
	MethodGET  = "GET"
	MethodPOST = "POST"
	// 目前仅支持两种签名方法
	SigMethodDefault    = "HmacSHA1"   // 默认签名方法
	SigMethodHmacSHA256 = "HmacSHA256" // 支持的签名方法，需指定
	// 目前仅支持一种签名版本
	DefaultVersion = "20191001" // 默认的签名版本
)

// Sign 签名函数， 返回加密后的信息
// source：源字符串
// secretKey：用户密钥
// signatureMethod：签名算法，目前支持HmacSHA1和HmacSHA256
func Sign(source string, secretKey string, signatureMethod string) (sign string) {
	// 使用HmacSHA256加密
	if signatureMethod == SigMethodHmacSHA256 {
		hmac256Obj := hmac.New(sha256.New, []byte(secretKey))
		hmac256Obj.Write([]byte(source))
		sign = base64.StdEncoding.EncodeToString(hmac256Obj.Sum(nil))
		return sign
	}

	// 使用HmacSHA1加密
	hmacObj := hmac.New(sha1.New, []byte(secretKey))
	hmacObj.Write([]byte(source))
	sign = base64.StdEncoding.EncodeToString(hmacObj.Sum(nil))
	return sign
}

// SendRequest 发送请求函数，返回json格式后的回包包体信息或错误信息
// host：目标host
// path: 目标url路径
// method: http请求方法，仅支持GET和POST
// secretId: 密钥ID
// secretKey: 密钥信息
// body: json格式化后的包体信息
// signatureMethod：签名算法，目前支持HmacSHA1和HmacSHA256
func SendRequest(host string, path string, method string, body []byte, secretId string, secretKey string,
	signatureMethod string) ([]byte, error) {

	var err error

	//版本号，必填
	version := DefaultVersion
	timestamp := fmt.Sprintf("%v", time.Now().Unix())
	rand.Seed(time.Now().UnixNano())
	nonce := fmt.Sprintf("%v", rand.Int())

	// 请求方法, 除非指定为GET， 否则，默认为POST
	var requestMethod string
	if strings.ToUpper(method) == MethodGET {
		requestMethod = MethodGET
	} else {
		requestMethod = MethodPOST
	}

	requestHost := host
	requestPath := path

	// 组url参数
	rawQuery := "Version=" + url.QueryEscape(version) + "&SecretId=" + url.QueryEscape(secretId) +
		"&Timestamp=" + url.QueryEscape(timestamp) + "&Nonce=" + url.QueryEscape(nonce)
	if signatureMethod == SigMethodHmacSHA256 {
		rawQuery += "&SignatureMethod=" + url.QueryEscape(SigMethodHmacSHA256)
	} else {
		signatureMethod = SigMethodDefault
	}
	// 包体非空时，进行包体签名
	if body != nil {
		if debug {
			fmt.Printf("payload[%v]\n", string(body))
		}
		hashedRequestPayload := Sign(string(body), secretKey, signatureMethod)
		rawQuery += "&HashedRequestPayload=" + url.QueryEscape(hashedRequestPayload)
		if debug {
			fmt.Printf("hashedRequestPayload[%v]\n", hashedRequestPayload)
		}
	}
	urlStr := requestHost + requestPath + "?" + rawQuery

	// 整包签名
	stringToSign := strings.ToUpper(requestMethod) + urlStr
	if debug {
		fmt.Printf("stringToSign[%v]\n", stringToSign)
	}
	signature := Sign(stringToSign, secretKey, signatureMethod)
	if debug {
		fmt.Printf("signature[%v]\n", signature)
	}
	urlStr += "&Signature=" + url.QueryEscape(signature)

	// 发送请求包
	httpUrl := "http://" + urlStr
	if debug {
		fmt.Printf("httpUrl(%v), method(%v)\n", httpUrl, requestMethod)
	}

	var rsp *http.Response

	// 执行GET请求
	if requestMethod == MethodGET {
		rsp, err = http.Get(httpUrl)
		if err != nil {
			fmt.Printf("http GET error:%v\n", err.Error())
			return nil, err
		}

		defer rsp.Body.Close()

		if rsp.StatusCode != http.StatusOK {
			errMsg := fmt.Sprintf("http's statusCode error:%v", rsp.StatusCode)
			fmt.Println(errMsg)
			return nil, errors.New(errMsg)
		}

		retData, err := ioutil.ReadAll(rsp.Body)
		if err != err {
			return nil, err
		}

		if debug {
			fmt.Printf("rsp[%v]\n", string(retData))
		}

		return retData, nil
	}

	// 执行POST请求
	rsp, err = http.Post(httpUrl, "application/json", bytes.NewBuffer(body))
	if err != nil {
		fmt.Printf("http POST error:%v\n", err.Error())
		return nil, err
	}

	defer rsp.Body.Close()

	if rsp.StatusCode != http.StatusOK {
		errMsg := fmt.Sprintf("http's statusCode error:%v", rsp.StatusCode)
		fmt.Println(errMsg)
		return nil, errors.New(errMsg)
	}

	retData, err := ioutil.ReadAll(rsp.Body)
	if err != err {
		return nil, err
	}

	if debug {
		fmt.Printf("rsp[%v]\n", string(retData))
	}

	return retData, nil
}
