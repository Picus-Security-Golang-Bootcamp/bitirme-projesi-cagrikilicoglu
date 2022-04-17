package middleware

import (
	"errors"

	"github.com/cagrikilicoglu/shopping-basket/internal/models/response"
	"github.com/cagrikilicoglu/shopping-basket/pkg/jwtHelper"
	"github.com/gin-gonic/gin"
)

// AdminAuthMiddleware checks authorization of the request and allows admin to continue
func AdminAuthMiddleware(secretKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetHeader("Authorization") != "" {
			decodedClaims := jwtHelper.VerifyToken(c.GetHeader("Authorization"), secretKey)
			if decodedClaims != nil {
				if string(decodedClaims.Roles) == "admin" {
					userID := decodedClaims.UserId
					c.Set("userID", userID)
					c.Next()
					c.Abort()
					return
				}
			}
			c.Abort()
			response.RespondWithError(c, errors.New("You are not allowed to use this endpoint"))
			return
		} else {
			c.Abort()
			response.RespondWithError(c, errors.New("Missing authorization"))
			return
		}

	}
}

// UserAuthMiddleware checks authorization of the request and allows user to continue
func UserAuthMiddleware(secretKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetHeader("Authorization") != "" {
			decodedClaims := jwtHelper.VerifyToken(c.GetHeader("Authorization"), secretKey)
			if decodedClaims != nil {
				if string(decodedClaims.Roles) == "user" || string(decodedClaims.Roles) == "admin" {
					userID := decodedClaims.UserId
					c.Set("userID", userID)
					c.Next()
					c.Abort()
					return
				}
			}
			c.Abort()
			response.RespondWithError(c, errors.New("You are not allowed to use this endpoint"))
			return
		} else {
			c.Abort()
			response.RespondWithError(c, errors.New("Missing authorization"))
			return
		}

	}
}

// RefreshMiddleware checks authorization to refresh tokens
func RefreshMiddleware(secretKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetHeader("Authorization") != "" {
			decodedClaims := jwtHelper.VerifyToken(c.GetHeader("Authorization"), secretKey)
			if decodedClaims != nil {
				if string(decodedClaims.Roles) == "user" || string(decodedClaims.Roles) == "admin" {
					c.Set("userID", decodedClaims.UserId)
					c.Set("email", decodedClaims.Email)
					c.Set("role", decodedClaims.Roles)
					c.Next()
					c.Abort()
					return
				}
			}
			c.Abort()
			response.RespondWithError(c, errors.New("You are not allowed to use this endpoint"))
			return
		} else {
			c.Abort()
			response.RespondWithError(c, errors.New("Missing authorization"))
			return
		}

	}
}
