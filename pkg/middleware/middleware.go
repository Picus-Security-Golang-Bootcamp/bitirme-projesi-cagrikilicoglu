package middleware

import (
	"errors"

	"github.com/cagrikilicoglu/shopping-basket/internal/models/response"
	"github.com/cagrikilicoglu/shopping-basket/pkg/jwtHelper"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func AdminAuthMiddleware(secretKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetHeader("Authorization") != "" {
			decodedClaims := jwtHelper.VerifyToken(c.GetHeader("Authorization"), secretKey)
			if decodedClaims != nil {
				// for _, role := range decodedClaims.Roles {
				// 	if string(role) == "admin" {
				// 		c.Next()
				// 		c.Abort()
				// 		return
				// 	}
				// }
				if string(decodedClaims.Roles) == "admin" {
					userID := decodedClaims.UserId
					c.Set("userID", userID)
					c.Next()
					c.Abort()
					return
				}
			}
			response.RespondWithError(c, errors.New("You are not allowed to use this endpoint"))
			c.Abort()
			return
		} else {
			c.Abort()
			response.RespondWithError(c, errors.New("Missing authorization"))
			// TODO respond with error missing authorization
			return
		}

	}
}
func UserAuthMiddleware(secretKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetHeader("Authorization") != "" {
			decodedClaims := jwtHelper.VerifyToken(c.GetHeader("Authorization"), secretKey)
			if decodedClaims != nil {
				// for _, role := range decodedClaims.Roles {
				// 	if role == "user" || role == "admin" {
				// 		c.Next()
				// 		c.Abort()
				// 		return
				// 	}
				// }
				zap.L().Debug("userAuthMid", zap.Reflect("decodedclaims", decodedClaims.Roles))
				if string(decodedClaims.Roles) == "user" || string(decodedClaims.Roles) == "admin" {
					userID := decodedClaims.UserId
					c.Set("userID", userID)
					c.Next()
					c.Abort()
					return
				}
			}
			response.RespondWithError(c, errors.New("You are not allowed to use this endpoint"))
			c.Abort()
			return
		} else {
			c.Abort()
			response.RespondWithError(c, errors.New("Missing authorization"))
			// TODO respond with error missing authorization
			return
		}

	}
}

// TODO düzelt
func RefreshMiddleware(secretKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetHeader("Authorization") != "" {
			decodedClaims := jwtHelper.VerifyToken(c.GetHeader("Authorization"), secretKey)
			if decodedClaims != nil {
				// for _, role := range decodedClaims.Roles {
				if string(decodedClaims.Roles) == "user" || string(decodedClaims.Roles) == "admin" {
					// userID := decodedClaims.UserId
					c.Set("userID", decodedClaims.UserId)
					c.Set("email", decodedClaims.Email)
					c.Set("role", decodedClaims.Roles)
					c.Next()
					c.Abort()
					return
				}
				// }
			}
			response.RespondWithError(c, errors.New("You are not allowed to use this endpoint"))
			c.Abort()
			return
		} else {
			c.Abort()
			response.RespondWithError(c, errors.New("Missing authorization"))
			// TODO respond with error missing authorization
			return
		}

	}
}
