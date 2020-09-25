package jwt

import (
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/dgrijalva/jwt-go"
	"github.com/valyala/fasthttp"
)

// GetKeycloakPublicKey 获取keycloak公钥
func GetKeycloakPublicKey(addr string) (string, error) {
	status, resp, err := fasthttp.Get(nil, addr)
	if err != nil {
		return "", fmt.Errorf("请求失败:%#v", err.Error())
	}
	if status != fasthttp.StatusOK {
		return "", fmt.Errorf("请求没有成功:%#v", status)
	}
	respdata := make(map[string]interface{})
	if err := json.Unmarshal(resp, &respdata); err != nil {
		return "", err
	}
	return respdata["public_key"].(string), nil
}

// ParseToken 解析jwt token获取信息(未经过公钥检验)
func ParseToken(tokenStr string) (map[string]interface{}, error) {
	if len(tokenStr) == 0 {
		return nil, fmt.Errorf("token is empty")
	}
	token, err := jwt.Parse(tokenStr, func(s *jwt.Token) (interface{}, error) {
		return nil, nil
	})
	// 解析失败则返回err
	claim, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, err
	}
	return claim, nil
}

// ParseTokenWithClaims 使用公钥解析token
func ParseTokenWithClaims(tokenStr string, publicKey string) (map[string]interface{}, error) {
	if len(tokenStr) == 0 {
		return nil, fmt.Errorf("token is empty")
	}
	token, err := jwt.Parse(tokenStr, func(s *jwt.Token) (interface{}, error) {
		// 通过公钥和token返回key，交由jwt去解析token
		signature := make([]byte, base64.StdEncoding.DecodedLen(len([]byte(publicKey))))
		// 根据签名部分和公钥通过base64解析出签名长度
		n, err := base64.StdEncoding.Decode(signature, []byte(publicKey))
		if err != nil {
			return nil, err
		}
		signature = signature[:n]
		// 生成数字签名解析
		key, err := x509.ParsePKIXPublicKey(signature)
		if err != nil {
			return nil, err
		}
		return key, nil
	})
	if err != nil {
		return nil, err
	}
	claim, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		err := errors.New("cannot convert claim to mapclaim")
		return nil, err
	}
	if !token.Valid {
		err := errors.New("token is invalid")
		return nil, err
	}
	return claim, nil
}
