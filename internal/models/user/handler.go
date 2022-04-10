package user

import (
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/cagrikilicoglu/shopping-basket/internal/api"
	"github.com/cagrikilicoglu/shopping-basket/internal/models/response"
	"github.com/cagrikilicoglu/shopping-basket/pkg/config"
	"github.com/cagrikilicoglu/shopping-basket/pkg/jwtHelper"
	"github.com/gin-gonic/gin"
	"github.com/go-openapi/strfmt"
	"github.com/golang-jwt/jwt"
	"go.uber.org/zap"
)

type userHandler struct {
	repo *UserRepository
	cfg  *config.Config
}

// type authHandler struct {
// 	cfg *config.Config
// }

type Tokens struct {
	AccessToken  string
	RefreshToken string
}

func NewUserHandler(r *gin.RouterGroup, repo *UserRepository) {
	h := &userHandler{repo: repo}

	r.POST("/signup", h.createUser)
	r.POST("/login", h.Login)

}

func (u *userHandler) createUser(c *gin.Context) {
	userBody := &api.User{}

	zap.L().Debug("User.handler.createUser.Bind", zap.Reflect("User", userBody))
	if err := c.Bind(&userBody); err != nil {
		response.RespondWithError(c, err)
		return
	}
	zap.L().Debug("User.handler.createUser.Validate", zap.Reflect("User", userBody))
	if err := userBody.Validate(strfmt.NewFormats()); err != nil {
		response.RespondWithError(c, err)
		return
	}

	user, err := u.repo.Create(responseToUser(userBody))
	if err != nil {
		response.RespondWithError(c, err)
		return
	}

	response.RespondWithJson(c, http.StatusCreated, user)
}

// TODO ayrı bi servis olmalı
func (u *userHandler) Login(c *gin.Context) {

	loginBody := api.Login{}
	if err := c.Bind(&loginBody); err != nil {
		response.RespondWithError(c, err)
		return
	}
	if err := loginBody.Validate(strfmt.NewFormats()); err != nil {
		response.RespondWithError(c, err)
		return
	}
	user, err := u.repo.GetUser(*loginBody.Email, *loginBody.Password)
	if err != nil || user == nil {
		response.RespondWithError(c, errors.New("wrong credentials"))
		return
	}

	jwtAccessClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": user.ID,
		"email":  user.Email,
		"iat":    time.Now().Unix(),
		"iss":    os.Getenv("APP_ENV"),
		"exp":    time.Now().Add(15 * time.Minute).Unix(),
		// "exp":   time.Now().Add(30 * time.Second).Unix(),
		"roles": user.Role,
	})

	accessToken := jwtHelper.GenerateToken(jwtAccessClaims, u.cfg.JWTConfig.SecretKey)

	jwtRefreshClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": user.ID,
		"email":  user.Email,
		"iat":    time.Now().Unix(),
		"iss":    os.Getenv("APP_ENV"),
		"exp":    time.Now().Add(7 * 24 * time.Hour).Unix(),
		"roles":  user.Role,
	})

	refreshToken := jwtHelper.GenerateToken(jwtRefreshClaims, u.cfg.JWTConfig.RefreshSecretKey)

	tokens := Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	response.RespondWithJson(c, http.StatusOK, tokens)

}

func (u *userHandler) VerifyAccessToken(c *gin.Context) {
	token := c.GetHeader("Authorization")
	decodedToken := jwtHelper.VerifyToken(token, u.cfg.JWTConfig.SecretKey)
	response.RespondWithJson(c, http.StatusOK, decodedToken)

}
func (u *userHandler) VerifyRefreshToken(c *gin.Context) {
	token := c.GetHeader("Authorization")
	decodedToken := jwtHelper.VerifyToken(token, u.cfg.JWTConfig.RefreshSecretKey)
	response.RespondWithJson(c, http.StatusOK, decodedToken)

}
