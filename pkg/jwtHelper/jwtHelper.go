package jwtHelper

import (
	"encoding/json"

	"github.com/golang-jwt/jwt"
	"go.uber.org/zap"
)

type DecodedToken struct {
	Iat    int    `json:"iat"`
	Roles  string `json:"roles"`
	UserId string `json:"userId"`
	Email  string `json:"email"`
	Iss    string `json:"iss"`
}

func GenerateToken(claims *jwt.Token, secret string) string {
	zap.L().Debug("jwtHelper.GenerateToken", zap.Reflect("token", &claims))

	hmacSecretStr := secret
	hmacSecret := []byte(hmacSecretStr)
	token, err := claims.SignedString(hmacSecret)

	// TODO erroru handle et
	if err != nil {
		// response.RespondWithError(c, err)
	}
	return token
}

func VerifyToken(token string, secret string) *DecodedToken {
	hmacSecretStr := secret
	hmacSecret := []byte(hmacSecretStr)
	decoded, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) { return hmacSecret, nil })

	// TODO handle et
	if err != nil {
		return nil
	}
	// TODO handle et
	if !decoded.Valid {
		return nil
	}

	decodedClaims := decoded.Claims.(jwt.MapClaims)

	var decodedToken DecodedToken
	jsonString, err := json.Marshal(decodedClaims)
	// TODO handle et
	if err != nil {
		return nil
	}
	json.Unmarshal(jsonString, &decodedToken)
	return &decodedToken

}
