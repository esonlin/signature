package middleware

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"signature/common"
	"strconv"
	"strings"
	"time"

	"signature/service"

	"github.com/gin-gonic/gin"
)

var (
	// 目前仅支持两种签名方法
	SigMethodDefault    = "HmacSHA1"   // 默认签名方法
	SigMethodHmacSHA256 = "HmacSHA256" // 支持的签名方法，需指定
	// 目前仅支持一种签名版本
	DefaultVersion = "20191001" // 默认的签名版本

	SignatureExpireTime = 5 * 60 //签名过期时间，只允许五分钟误差
)

func sign(source string, secretKey string, signatureMethod string) (sign string) {
	// 使用HmacSHA256加密
	if signatureMethod == SigMethodHmacSHA256 {
		hmacObj := hmac.New(sha256.New, []byte(secretKey))
		hmacObj.Write([]byte(source))
		sign = base64.StdEncoding.EncodeToString(hmacObj.Sum(nil))
		return sign
	}

	// 使用HmacSHA1加密
	hmacObj := hmac.New(sha1.New, []byte(secretKey))
	hmacObj.Write([]byte(source))
	sign = base64.StdEncoding.EncodeToString(hmacObj.Sum(nil))
	return sign
}

// SignatureMiddleware 验证签名
func SignatureMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		requestID := common.GenerateRequestID()

		// 步骤1、获取请求签名
		requestSig := c.Query("Signature")
		if requestSig == "" {
			//没找到签名信息，拒绝访问
			fmt.Println("not found url'param Signature.")
			common.ErrorResponse(c, requestID, common.SignatureNotFound, "not found Signature")
			c.Abort()
			return
		}
		fmt.Println("Signature:", requestSig)

		//todo: 2、获取时间戳
		strTimestamp := c.Query("Timestamp")
		var errMsg string
		var expireFlag bool
		if strTimestamp != "" {
			timestamp, err := strconv.ParseInt(strTimestamp, 10, 64)
			if err != nil {
				//时间戳解析失败，也认为签名过期，拒绝访问
				expireFlag = true
				errMsg = "invalid param: Timestamp."
			} else {
				//签名时间有效性验证，允许5分钟偏差；超过则拒绝访问
				curTime := time.Now().Unix()
				if curTime-timestamp > int64(SignatureExpireTime) || timestamp-curTime > int64(SignatureExpireTime) {
					expireFlag = true
					errMsg = "signature expire."
				} else {
					expireFlag = false
				}
			}
		} else {
			//没找到时间戳，也认为签名过期，拒绝访问
			expireFlag = true
			errMsg = "not found url'param Timestamp."
		}
		if expireFlag {
			fmt.Println("Timestamp err:", errMsg)
			common.ErrorResponse(c, requestID, common.SignatureExpire, errMsg)
			c.Abort()
			return
		}

		// todo: 3、获取SecretId
		secretId := c.Query("SecretId")
		if secretId == "" {
			//没找到SecretId，拒绝访问
			fmt.Println("not found url'param SecretId.")
			common.ErrorResponse(c, requestID, common.SecretIDNotFound, "not found SecretId")
			c.Abort()
			return
		}
		// 通过secretId查找SecretKey
		secretKey, err := service.GetSecretKeyBySecretID(secretId)
		if err != nil {
			//密钥信息查找失败，拒绝访问
			errMsg := fmt.Sprintf("mysql operation(GetSecretInfoByID) error:%s", err.Error())
			fmt.Println(errMsg)
			common.ErrorResponse(c, requestID, common.InvalidSecetID, errMsg)
			c.Abort()
			return
		}

		// todo: 4、获取签名算法, 默认使用HmacSHA1
		// 计算签名的方法有两种：HmacSHA256 和 HmacSHA1
		signatureMethod := c.Query("SignatureMethod")
		if signatureMethod != "" {
			if strings.ToUpper(signatureMethod) == strings.ToUpper(SigMethodHmacSHA256) {
				signatureMethod = SigMethodHmacSHA256
			} else {
				signatureMethod = SigMethodDefault
			}
		} else {
			signatureMethod = SigMethodDefault
		}

		// todo: 5、获取签名版本
		version := c.Query("Version")
		if version == "" {
			version = DefaultVersion
		}

		// todo: 6、获取包体签名 -- 没有包体签名时，不验包体签名
		hashedRequestPayload := c.Query("HashedRequestPayload")
		if hashedRequestPayload != "" {

			body, err := ioutil.ReadAll(c.Request.Body)
			if err != nil {
				errMsg := fmt.Sprintf("Request' Payload parse err:%v", err.Error())
				fmt.Println(errMsg)
				common.ErrorResponse(c, requestID, common.PayLoadSigFailure, errMsg)
				c.Abort()
				return
			}
			c.Request.Body.Close()
			fmt.Printf("requestPayload:[%v]\n", string(body))

			// 包体读完必须重新赋值，否则后续流程将读取不到包体
			c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))

			//todo: 6.1、验证包体签名
			payLoadSignature := sign(string(body), secretKey, signatureMethod)
			fmt.Printf("requestSig(%v)--payloadSig(%v)\n", hashedRequestPayload, payLoadSignature)
			if payLoadSignature != hashedRequestPayload {
				fmt.Printf("Request' Payload Signature verify failed:req(%v)-verify(%v)\n", hashedRequestPayload, payLoadSignature)
				common.ErrorResponse(c, requestID, common.PayLoadSigFailure, "Request' Payload Signature verify failed")
				c.Abort()
				return
			}
		}

		// todo: 7、获取Nonce值
		nonce := c.Query("Nonce")
		if nonce == "" {
			// nonce暂时不用，但不能缺失
			errMsg := fmt.Sprintf("not found url'param Nonce.")
			fmt.Println(errMsg)
			common.ErrorResponse(c, requestID, common.NonceNotFound, errMsg)
			c.Abort()
			return
		}
		// todo: Nonce一次性有效校验--暂时不实现

		index := strings.Index(c.Request.RequestURI, "Signature=")
		stringToSig := strings.ToUpper(c.Request.Method) + c.Request.Host + c.Request.RequestURI[0:index-1]

		fmt.Printf("string to signature:(%v)\n", stringToSig)
		fmt.Printf("secretKey:(%v)\n", secretKey)

		// todo: 8、请求签名验证
		signature := sign(stringToSig, secretKey, signatureMethod)

		fmt.Printf("reqSig(%v)-verifySig(%v)\n", requestSig, signature)
		if signature != requestSig {
			fmt.Printf("Signature verify failed:req(%v)-verify(%v)\n", requestSig, signature)
			common.ErrorResponse(c, requestID, common.SignatureFailure, "Signature verify failed")
			c.Abort()
			return
		}

		// todo: 9、按需设置验证成功的信息，便于后续流程使用，此处省略......
		//var paramAppID gin.Param
		//paramAppID.Key = "application_id"
		//paramAppID.Value = fmt.Sprint(secretInfo.APPID)
		//c.Params = append(c.Params, paramAppID)
		//
		//var paramUser gin.Param
		//paramUser.Key = "user_id"
		//paramUser.Value = secretInfo.UserID
		//c.Params = append(c.Params, paramUser)

		fmt.Println("SignatureMiddleware success.")
		c.Next()
		return
	}
}
