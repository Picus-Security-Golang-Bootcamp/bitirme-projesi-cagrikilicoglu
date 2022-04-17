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

// GenerateToken creates a new token by signing it
func GenerateToken(claims *jwt.Token, secret string) (*string, error) {

	zap.L().Debug("jwtHelper.GenerateToken", zap.Reflect("claims", &claims))

	hmacSecret := []byte(secret)
	token, err := claims.SignedString(hmacSecret)

	if err != nil {
		zap.L().Error("jwtHelper.GenerateToken failed to generate Token", zap.Error(err))
		return nil, err
	}
	return &token, nil
}

// VerifyToken verifies a token by decoding it
func VerifyToken(token string, secret string) *DecodedToken {
	hmacSecretStr := secret
	hmacSecret := []byte(hmacSecretStr)
	decoded, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) { return hmacSecret, nil })

	if err != nil {
		return nil
	}

	if !decoded.Valid {
		return nil
	}

	decodedClaims := decoded.Claims.(jwt.MapClaims)

	var decodedToken DecodedToken
	jsonString, err := json.Marshal(decodedClaims)

	if err != nil {
		return nil
	}
	json.Unmarshal(jsonString, &decodedToken)
	return &decodedToken

}
