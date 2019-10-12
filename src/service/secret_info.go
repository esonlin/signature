package service

import "errors"

//
func GetSecretKeyBySecretID(secretId string) (string, error) {

	// todo: 密钥信息一般需要加密存储到数据库中，此处省略......

	// 此处，仅写死测试用的密钥信息
	if secretId == "SKIDz8krbsJ5yKBZQpn74WFkmLPx3EXAMPLE" {
		return "Gu5t9xGARNpq86cd98joQYCN3EXAMPLE", nil
	} else {
		return "", errors.New("not found secret info")
	}
}
