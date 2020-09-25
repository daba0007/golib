# JWT

jwt包用于处理JWT-TOKEN

## 获取keycloak的公钥

```Go
publicKey, err := GetKeycloakPublicKey("http://localhost:8080//auth/realms/Master")
```

## 直接获取token的信息

```Go
// 直接获取token的信息，会忽略掉`jwt.Parse` 函数中的valid错误，仅在无token或token非法时报错
// token -> jwttoken string
data, err := ParseToken(token)
// 获取iss信息
iss,ok := data["iss"].(string)
```

## 通过公解析jtoken的信息

```Go
// 通过公钥解析jtoken的信息，会捕捉`jwt.Parse` 函数中的所有错误，仅在正确解析时返回
// token -> jwttoken
// publicKey -> public key 
data, err := ParseTokenWithClaims(token, publicKey)
// 获取iss信息
iss,ok := data["iss"].(string)
```