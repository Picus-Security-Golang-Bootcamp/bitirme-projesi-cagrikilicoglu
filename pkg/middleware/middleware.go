package middleware

import (
	"errors"

	"github.com/cagrikilicoglu/shopping-basket/internal/models/response"
	"github.com/cagrikilicoglu/shopping-basket/pkg/jwtHelper"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(secretKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetHeader("Authorization") != "" {
			decodedClaims := jwtHelper.VerifyToken(c.GetHeader("Authorization"), secretKey)
			if decodedClaims != nil {
				for _, role := range decodedClaims.Roles {
					if role == "admin" {
						c.Next()
						c.Abort()
						return
					}
				}
			}
			response.RespondWithError(c, errors.New("You are not allowed to use this endpoint"))
			c.Abort()
			return
		} else {
			c.Abort()
			// TODO respond with error missing authorization
			return
		}

	}
}

func RefreshMiddleware(secretKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetHeader("Authorization") != "" {
			decodedClaims := jwtHelper.VerifyToken(c.GetHeader("Authorization"), secretKey)
			if decodedClaims != nil {
				for _, role := range decodedClaims.Roles {
					if role == "admin" {
						c.Next()
						c.Abort()
						return
					}
				}
			}
			response.RespondWithError(c, errors.New("You are not allowed to use this endpoint"))
			c.Abort()
			return
		} else {
			c.Abort()
			// TODO respond with error missing authorization
			return
		}

	}
}
