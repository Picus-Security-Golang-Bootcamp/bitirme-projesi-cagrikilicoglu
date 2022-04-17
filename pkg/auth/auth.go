package auth

import (
	"net/http"
	"os"
	"time"

	"github.com/cagrikilicoglu/shopping-basket/internal/models"
	"github.com/cagrikilicoglu/shopping-basket/internal/models/response"
	"github.com/cagrikilicoglu/shopping-basket/pkg/config"
	"github.com/cagrikilicoglu/shopping-basket/pkg/jwtHelper"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Authenticator struct {
	Cfg *config.Config
}

func NewAuthenticator(cfg *config.Config) *Authenticator {
	return &Authenticator{Cfg: cfg}
}

// Authenticate returns access and refresh tokens by id, email and role inputs from the user
func (a *Authenticator) Authenticate(id uuid.UUID, email, role string) (*models.Tokens, error) {

	zap.L().Debug("authenticator.Authenticate.jwtNewWithClaims.Access",
		zap.Reflect("id", id),
		zap.Reflect("email", email),
		zap.Reflect("role", role))
	jwtAccessClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": id,
		"email":  email,
		"roles":  role,
		"iss":    os.Getenv("APP_ENV"),
		"iat":    time.Now().Unix(),
		"exp":    time.Now().Add(time.Duration(a.Cfg.JWTConfig.AccessTokenDurationMins) * time.Minute).Unix(),
	})

	zap.L().Debug("authenticator.Authenticate.GenerateToken.access", zap.Reflect("jwtAccessClaims", jwtAccessClaims))
	accessToken, err := jwtHelper.GenerateToken(jwtAccessClaims, a.Cfg.JWTConfig.SecretKey)
	if err != nil {
		return nil, err
	}

	zap.L().Debug("authenticator.Authenticate.jwtNewWithClaims.Refresh")
	jwtRefreshClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": id,
		"email":  email,
		"roles":  role,
		"iss":    os.Getenv("APP_ENV"),
		"iat":    time.Now().Unix(),
		"exp":    time.Now().Add(time.Duration(a.Cfg.JWTConfig.RefreshTokenDurationHours) * time.Hour).Unix(),
	})

	zap.L().Debug("authenticator.Authenticate.GenerateToken.refresh", zap.Reflect("jwtRefreshClaims", jwtRefreshClaims))
	refreshToken, err := jwtHelper.GenerateToken(jwtRefreshClaims, a.Cfg.JWTConfig.RefreshSecretKey)
	if err != nil {
		return nil, err
	}

	tokens := models.Tokens{
		AccessToken:  *accessToken,
		RefreshToken: *refreshToken,
	}
	return &tokens, err
}

// VerifyAccessToken verifies the access token by checking secret key
func (a *Authenticator) VerifyAccessToken(c *gin.Context) {
	token := c.GetHeader("Authorization")
	decodedToken := jwtHelper.VerifyToken(token, a.Cfg.JWTConfig.SecretKey)
	response.RespondWithJson(c, http.StatusOK, decodedToken)

}

// VerifyRefreshToken verifies the refresh token by checking secret key
func (a *Authenticator) VerifyRefreshToken(c *gin.Context) {
	token := c.GetHeader("Authorization")
	decodedToken := jwtHelper.VerifyToken(token, a.Cfg.JWTConfig.RefreshSecretKey)
	response.RespondWithJson(c, http.StatusOK, decodedToken)

}
