package main

import (
	pb "github.com/Jimmy01010/protocol/shippy-user"
	"github.com/golang-jwt/jwt"
	"go-micro.dev/v4/util/log"
)

var (
	// Define a secure key string used
	// as a salt when hashing our tokens.
	// Please make your own way more secure than this,
	// use a randomly generated md5 hash or something.
	// 注意生产环境中替换成更安全的值而且妥善保管
	privateKey = []byte("mySuperSecretKey")
)

// CustomClaims is our custom metadata, which will be hashed
// and sent as the second segment in our JWT
type CustomClaims struct {
	User *pb.User
	// 嵌入标准的 payload
	jwt.StandardClaims
}

type TokenService struct {
	repo Repository
}

// Decode 解码JWT字符串为 CustomClaims 对象
func (srv *TokenService) Decode(tokenString string) (*CustomClaims, error) {
	// jwt.ParseWithClaims 解析、验证并返回一个Token, keyFunc将接收已解析的令牌并应返回用于验证的密钥。
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		// 我们从token的Claims中获取用户元数据以验证该用户。如果验证通过，那keyFunc的err就返回nil
		// 这里我们省略验证步骤，直接返回了密钥和nil
		return privateKey, nil
	})

	// 断言为CustomClaims类型，即验证其信息是不是用户信息
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		log.Infof("%v %v", claims.User, claims.StandardClaims.ExpiresAt)
		return claims, nil
	} else {
		log.Error("token valid: ", err)
		return nil, err
	}
}

// Encode a claim into a JWT 将 User 用户信息编码为 JWT 字符串
// 将我们自定义的claim(JWT的payload部分)编码进一个JWT中, 然后返回带有签名的(Signature)完整Token
func (srv *TokenService) Encode(user *pb.User) (string, error) {
	// 创建我们自定义的Claims
	claims := CustomClaims{
		user,
		jwt.StandardClaims{
			// 测试的时候不设置过期时间
			// ExpiresAt: time.Now().Add(time.Minute * 1).Unix(),
			Issuer: "shippy.service.user",
		},
	}

	// Create token
	// HS256
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token and return
	return token.SignedString(privateKey)
}
