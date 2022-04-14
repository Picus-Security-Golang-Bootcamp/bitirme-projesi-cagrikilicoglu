package user

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/cagrikilicoglu/shopping-basket/internal/api"
	"github.com/cagrikilicoglu/shopping-basket/internal/models"
	"github.com/cagrikilicoglu/shopping-basket/internal/models/response"
	"github.com/cagrikilicoglu/shopping-basket/pkg/auth"
	"github.com/cagrikilicoglu/shopping-basket/pkg/middleware"
	"github.com/gin-gonic/gin"
	"github.com/go-openapi/strfmt"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type userHandler struct {
	repo *UserRepository
	auth *auth.Authenticator
}

// type authHandler struct {
// 	cfg *config.Config
// }

// type Tokens struct {
// 	AccessToken  string
// 	RefreshToken string
// }

func NewUserHandler(r *gin.RouterGroup, repo *UserRepository, auth *auth.Authenticator) {
	h := &userHandler{repo: repo,
		auth: auth}

	r.POST("/signup", h.createUser)
	r.POST("/login", h.Login)
	r.POST("/refresh", middleware.RefreshMiddleware(h.auth.Cfg.JWTConfig.RefreshSecretKey), h.Refresh)

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
	var tokens models.Tokens
	tokens = u.auth.Authenticate(user.ID, *user.Email, user.Role)
	response.RespondWithJson(c, http.StatusCreated, tokens)

	// response.RespondWithJson(c, http.StatusCreated, user)
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
	// encryptedPassword := getHash([]byte(*loginBody.Password))
	zap.L().Debug("User.handler.loginUser.PASS", zap.Reflect("encrypt", *loginBody.Password))
	// zap.L().Debug("User.handler.loginUser.encrypt", zap.Reflect("encrypt", encryptedPassword))
	user, err := u.repo.GetUser(*loginBody.Email, *loginBody.Password)
	if err != nil || user == nil {
		response.RespondWithError(c, errors.New("wrong credentials"))
		return
	}

	// jwtAccessClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
	// 	"userID": user.ID,
	// 	"email":  user.Email,
	// 	"iat":    time.Now().Unix(),
	// 	"iss":    os.Getenv("APP_ENV"),
	// 	"exp":    time.Now().Add(15 * time.Minute).Unix(),
	// 	// "exp":   time.Now().Add(30 * time.Second).Unix(),
	// 	"roles": user.Role,
	// })

	// accessToken := jwtHelper.GenerateToken(jwtAccessClaims, u.cfg.JWTConfig.SecretKey)

	// jwtRefreshClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
	// 	"userID": user.ID,
	// 	"email":  user.Email,
	// 	"iat":    time.Now().Unix(),
	// 	"iss":    os.Getenv("APP_ENV"),
	// 	"exp":    time.Now().Add(7 * 24 * time.Hour).Unix(),
	// 	"roles":  user.Role,
	// })

	// refreshToken := jwtHelper.GenerateToken(jwtRefreshClaims, u.cfg.JWTConfig.RefreshSecretKey)

	// tokens := Tokens{
	// 	AccessToken:  accessToken,
	// 	RefreshToken: refreshToken,
	// }
	var tokens models.Tokens
	tokens = u.auth.Authenticate(user.ID, *user.Email, user.Role)
	response.RespondWithJson(c, http.StatusOK, tokens)

}

func (u *userHandler) Refresh(c *gin.Context) {
	//TODO error handling yaz
	userID, _ := c.Get("userID")
	email, _ := c.Get("email")
	role, _ := c.Get("role")

	userIDParsed, err := uuid.Parse(fmt.Sprintf("%v", userID))
	if err != nil {
		response.RespondWithError(c, err)
		return
	}
	var tokens models.Tokens
	tokens = u.auth.Authenticate(userIDParsed, email.(string), role.(string))
	response.RespondWithJson(c, http.StatusOK, tokens)
}

// func (u *userHandler) VerifyAccessToken(c *gin.Context) {
// 	token := c.GetHeader("Authorization")
// 	decodedToken := jwtHelper.VerifyToken(token, u.cfg.JWTConfig.SecretKey)
// 	response.RespondWithJson(c, http.StatusOK, decodedToken)

// }
// func (u *userHandler) VerifyRefreshToken(c *gin.Context) {
// 	token := c.GetHeader("Authorization")
// 	decodedToken := jwtHelper.VerifyToken(token, u.cfg.JWTConfig.RefreshSecretKey)
// 	response.RespondWithJson(c, http.StatusOK, decodedToken)

// }

// func (u *userHandler) Authenticate(id uuid.UUID, email, role string) Tokens {

// 	jwtAccessClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
// 		"userID": id,
// 		"email":  email,
// 		"iat":    time.Now().Unix(),
// 		"iss":    os.Getenv("APP_ENV"),
// 		"exp":    time.Now().Add(15 * time.Minute).Unix(),
// 		// "exp":   time.Now().Add(30 * time.Second).Unix(),
// 		"roles": role,
// 	})

// 	accessToken := jwtHelper.GenerateToken(jwtAccessClaims, u.cfg.JWTConfig.SecretKey)

// 	jwtRefreshClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
// 		"userID": id,
// 		"email":  email,
// 		"iat":    time.Now().Unix(),
// 		"iss":    os.Getenv("APP_ENV"),
// 		"exp":    time.Now().Add(7 * 24 * time.Hour).Unix(),
// 		"roles":  role,
// 	})

// 	refreshToken := jwtHelper.GenerateToken(jwtRefreshClaims, u.cfg.JWTConfig.RefreshSecretKey)

// 	tokens := Tokens{
// 		AccessToken:  accessToken,
// 		RefreshToken: refreshToken,
// 	}

// 	return tokens
// }
