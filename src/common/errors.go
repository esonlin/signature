package common

// ErrorCode 全局统一错误返回码
type ErrorCode int64

const (
	Success            ErrorCode = iota //成功返回
	InvalidParameter                    //无效参数
	SignatureNotFound                   //未找到签名
	InvalidSecetID                      //无效密钥ID
	SecretIDNotFound                    //密钥不存在。请到控制台查看密钥是否被禁用，是否少复制了字符或者多了字符。
	SignatureExpire                     //签名过期。Timestamp 与服务器接收到请求的时间相差不得超过五分钟。
	SignatureFailure                    //签名错误。可能是签名计算错误，或者签名与实际发送的内容不相符合，也有可能是密钥 SecretKey 错误导致的。
	PayLoadSigFailure                   //包体签名错误。可能是签名计算错误，或者签名与实际发送的内容不相符合，也有可能是密钥 SecretKey 错误导致的。
	NonceNotFound                       //未找到Nonce值
	StatusUnauthorized                  //未授权
)
