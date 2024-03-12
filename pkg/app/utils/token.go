package utils

import (
	"fmt"
	"github.com/go-pay/gopay/pkg/jwt"
	"time"
)

// 自定义一个字符串
var jwtKey = []byte("wuhanfanxingwangluo")
var str string

type Claims struct {
	UserId int64
	jwt.StandardClaims
}

// AddToken 颁发token
func AddToken(userID int64, expireTime time.Time) string {
	claims := &Claims{
		UserId: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(), // 过期时间
			IssuedAt:  time.Now().Unix(),
			Issuer:    "mohhom", // 签名颁发者
			Subject:   "dadmin", // 签名主题
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// fmt.Println(token)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		fmt.Println(err)
	}
	return tokenString
}

// CheckToken 解析token
func CheckToken(tokenString string) (userID int64, bo bool) {
	if tokenString == "" {
		return 0, false
	}

	token, claims, err := ParseToken(tokenString)
	if err != nil || !token.Valid {
		return 0, false
	}
	fmt.Println("CheckToken 解析token")
	fmt.Println(claims.UserId)
	return claims.UserId, true
}

// ParseToken 解析token
func ParseToken(tokenString string) (*jwt.Token, *Claims, error) {
	Claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, Claims, func(token *jwt.Token) (i interface{}, err error) {
		return jwtKey, nil
	})
	return token, Claims, err
}
