package common

import (
	"github/Connector-Gamefi/ConnectorGoSDK/config"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

var jwtSecret = []byte(config.GetServerConfig().JWTSecret)

type MyClaims struct {
	Accont  string `json:"account"`
	Address string `json:"address"`
	jwt.StandardClaims
}

func GenerateToken(account, address string) (string, error) {
	now := time.Now()
	expireTime := now.Add(time.Duration(config.GetServerConfig().JWTExpireTimeMinute * int64(time.Minute)))

	claims := MyClaims{
		account,
		address,
		jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    "GOSDK",
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(jwtSecret)

	return token, err
}

func ParseToken(token string) (*MyClaims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*MyClaims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}

	return nil, err
}
