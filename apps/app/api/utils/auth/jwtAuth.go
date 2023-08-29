package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type JwtAuth struct {
	AccessSecret []byte
	AccessExpire int64
}
type CustomClaims struct {
	ID int64
	jwt.RegisteredClaims
}

// CreateToken 创建一个token
func (j *JwtAuth) CreateToken(tokenID int64) (string, error) {
	claims := CustomClaims{
		tokenID,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(j.AccessExpire) * time.Second)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "api",
			ID:        "1",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.AccessSecret)
}

// ParseToken 解析 token
func (j *JwtAuth) ParseToken(tokenString string) (int64, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.AccessSecret, nil
	})
	if err != nil {
		return -1, err
	}
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		//fmt.Printf("%v %v", claims.ID, claims.RegisteredClaims.Issuer)
		return claims.ID, nil
	} else {
		return -1, jwt.ErrTokenSignatureInvalid
	}

}
