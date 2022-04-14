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

func (a *Authenticator) Authenticate(id uuid.UUID, email, role string) models.Tokens {

	// zap.L().Debug("authenticator.authenticate.jwtaccessclaims.before",
	// 	zap.Reflect("key", a.cfg.JWTConfig.SecretKey))
	// TODO aşağıdakiler debug için
	zap.L().Debug("authenticator.authenticate.jwtaccessclaims",
		zap.Reflect("id", id),
		zap.Reflect("email", email),
		zap.Reflect("role", role))
	jwtAccessClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": id,
		"email":  email,
		"iat":    time.Now().Unix(),
		"iss":    os.Getenv("APP_ENV"),
		"exp":    time.Now().Add(15 * time.Minute).Unix(),
		// "exp":   time.Now().Add(30 * time.Second).Unix(),
		"roles": role,
	})
	zap.L().Debug("authenticator.authenticate.jwtaccessclaims",
		zap.Reflect("jwtAccesClaims", jwtAccessClaims))
	zap.L().Debug("authenticator.authenticate.secret",
		zap.Reflect("secretKey", a.Cfg.JWTConfig.SecretKey))

	accessToken := jwtHelper.GenerateToken(jwtAccessClaims, a.Cfg.JWTConfig.SecretKey)
	zap.L().Debug("authenticator.authenticate.jwtaccessclaims",
		zap.Reflect("accesstoken", accessToken))

	jwtRefreshClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": id,
		"email":  email,
		"iat":    time.Now().Unix(),
		"iss":    os.Getenv("APP_ENV"),
		"exp":    time.Now().Add(7 * 24 * time.Hour).Unix(),
		"roles":  role,
	})

	refreshToken := jwtHelper.GenerateToken(jwtRefreshClaims, a.Cfg.JWTConfig.RefreshSecretKey)

	tokens := models.Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return tokens
}
func (a *Authenticator) VerifyAccessToken(c *gin.Context) {
	token := c.GetHeader("Authorization")
	decodedToken := jwtHelper.VerifyToken(token, a.Cfg.JWTConfig.SecretKey)
	response.RespondWithJson(c, http.StatusOK, decodedToken)

}
func (a *Authenticator) VerifyRefreshToken(c *gin.Context) {
	token := c.GetHeader("Authorization")
	decodedToken := jwtHelper.VerifyToken(token, a.Cfg.JWTConfig.RefreshSecretKey)
	response.RespondWithJson(c, http.StatusOK, decodedToken)

}

func (a *Authenticator) Refresh(id uuid.UUID, email, role string) models.Tokens {

	// zap.L().Debug("authenticator.authenticate.jwtaccessclaims.before",
	// 	zap.Reflect("key", a.cfg.JWTConfig.SecretKey))
	// TODO aşağıdakiler debug için
	zap.L().Debug("authenticator.refresh.jwtaccessclaims",
		zap.Reflect("id", id),
		zap.Reflect("email", email),
		zap.Reflect("role", role))
	jwtAccessClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": id,
		"email":  email,
		"iat":    time.Now().Unix(),
		"iss":    os.Getenv("APP_ENV"),
		"exp":    time.Now().Add(15 * time.Minute).Unix(),
		// "exp":   time.Now().Add(30 * time.Second).Unix(),
		"roles": role,
	})
	zap.L().Debug("authenticator.refresh.jwtaccessclaims",
		zap.Reflect("jwtAccesClaims", jwtAccessClaims))
	zap.L().Debug("authenticator.refresh.secret",
		zap.Reflect("secretKey", a.Cfg.JWTConfig.SecretKey))

	accessToken := jwtHelper.GenerateToken(jwtAccessClaims, a.Cfg.JWTConfig.SecretKey)
	zap.L().Debug("authenticator.authenticate.jwtaccessclaims",
		zap.Reflect("accesstoken", accessToken))

	jwtRefreshClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": id,
		"email":  email,
		"iat":    time.Now().Unix(),
		"iss":    os.Getenv("APP_ENV"),
		"exp":    time.Now().Add(7 * 24 * time.Hour).Unix(),
		"roles":  role,
	})

	refreshToken := jwtHelper.GenerateToken(jwtRefreshClaims, a.Cfg.JWTConfig.RefreshSecretKey)

	tokens := models.Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return tokens
}

// func (a *Authenticator) Refresh(c *gin.Context) {
// 	token := c.GetHeader("Authorization")
// 	if decodedToken := jwtHelper.VerifyToken(token, a.cfg.JWTConfig.RefreshSecretKey); decodedToken != nil {

// 		jwtAccessClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
// 			"userID": decodedToken.UserId,
// 			"email":  decodedToken.Email,
// 			"iat":    time.Now().Unix(),
// 			"iss":    os.Getenv("APP_ENV"),
// 			"exp":    time.Now().Add(15 * time.Minute).Unix(),
// 			// "exp":   time.Now().Add(30 * time.Second).Unix(),
// 			"roles": decodedToken.Roles,
// 		})

// 		accessToken := jwtHelper.GenerateToken(jwtAccessClaims, a.cfg.JWTConfig.SecretKey)

// 		jwtRefreshClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
// 			"userID": decodedToken.UserId,
// 			"email":  decodedToken.Email,
// 			"iat":    time.Now().Unix(),
// 			"iss":    os.Getenv("APP_ENV"),
// 			"exp":    time.Now().Add(7 * 24 * time.Hour).Unix(),
// 			"roles":  decodedToken.Roles,
// 		})

// 		refreshToken := jwtHelper.GenerateToken(jwtRefreshClaims, a.cfg.JWTConfig.RefreshSecretKey)

// 		tokens := Tokens{
// 			AccessToken:  accessToken,
// 			RefreshToken: refreshToken,
// 		}
// 		response.RespondWithJson(c, http.StatusOK, tokens)

// 	} else {
// 		response.RespondWithError(c, errors.New("token is expired"))
// 	}

// 	// TODO önceki tokenı iptal et

// }

// // import (
// 	"errors"
// 	"net/http"
// 	"os"
// 	"time"

// 	"github.com/cagrikilicoglu/shopping-basket/internal/api"
// 	"github.com/cagrikilicoglu/shopping-basket/internal/models/response"
// 	"github.com/cagrikilicoglu/shopping-basket/pkg/config"
// 	"github.com/cagrikilicoglu/shopping-basket/pkg/jwtHelper"
// 	"github.com/cagrikilicoglu/shopping-basket/pkg/middleware"
// 	"github.com/dgrijalva/jwt-go"
// 	"github.com/gin-gonic/gin"
// 	"github.com/go-openapi/strfmt"
// )

// type authHandler struct {
// 	cfg *config.Config
// }

// type Tokens struct {
// 	AccessToken  string
// 	RefreshToken string
// }

// func NewAuthHandler(r *gin.RouterGroup, cfg *config.Config) {
// 	a := authHandler{cfg: cfg}
// 	r.POST("/login", a.Login)
// 	r.POST("/login/refresh", a.Refresh)

// 	r.Use(middleware.AuthMiddleware(cfg.JWTConfig.SecretKey))
// 	r.POST("/decode", a.VerifyAccessToken)

// }

// func (a *authHandler) Refresh(c *gin.Context) {
// 	token := c.GetHeader("Authorization")
// 	if decodedToken := jwtHelper.VerifyToken(token, a.cfg.JWTConfig.RefreshSecretKey); decodedToken != nil {

// 		jwtAccessClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
// 			"userID": decodedToken.UserId,
// 			"email":  decodedToken.Email,
// 			"iat":    time.Now().Unix(),
// 			"iss":    os.Getenv("APP_ENV"),
// 			"exp":    time.Now().Add(15 * time.Minute).Unix(),
// 			// "exp":   time.Now().Add(30 * time.Second).Unix(),
// 			"roles": decodedToken.Roles,
// 		})

// 		accessToken := jwtHelper.GenerateToken(jwtAccessClaims, a.cfg.JWTConfig.SecretKey)

// 		jwtRefreshClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
// 			"userID": decodedToken.UserId,
// 			"email":  decodedToken.Email,
// 			"iat":    time.Now().Unix(),
// 			"iss":    os.Getenv("APP_ENV"),
// 			"exp":    time.Now().Add(7 * 24 * time.Hour).Unix(),
// 			"roles":  decodedToken.Roles,
// 		})

// 		refreshToken := jwtHelper.GenerateToken(jwtRefreshClaims, a.cfg.JWTConfig.RefreshSecretKey)

// 		tokens := Tokens{
// 			AccessToken:  accessToken,
// 			RefreshToken: refreshToken,
// 		}
// 		response.RespondWithJson(c, http.StatusOK, tokens)

// 	} else {
// 		response.RespondWithError(c, errors.New("token is expired"))
// 	}

// 	// TODO önceki tokenı iptal et

// }

// func (a *authHandler) Login(c *gin.Context) {

// 	loginBody := api.Login{}
// 	if err := c.Bind(&loginBody); err != nil {
// 		response.RespondWithError(c, err)
// 		return
// 	}
// 	if err := loginBody.Validate(strfmt.NewFormats()); err != nil {
// 		response.RespondWithError(c, err)
// 		return
// 	}
// 	user := GetUser(*loginBody.Email, *loginBody.Password)
// 	if user == nil {
// 		response.RespondWithError(c, errors.New("wrong credentials"))
// 	}

// 	jwtAccessClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
// 		"userID": user.Id,
// 		"email":  user.Email,
// 		"iat":    time.Now().Unix(),
// 		"iss":    os.Getenv("APP_ENV"),
// 		"exp":    time.Now().Add(15 * time.Minute).Unix(),
// 		// "exp":   time.Now().Add(30 * time.Second).Unix(),
// 		"roles": user.Roles,
// 	})

// 	accessToken := jwtHelper.GenerateToken(jwtAccessClaims, a.cfg.JWTConfig.SecretKey)

// 	jwtRefreshClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
// 		"userID": user.Id,
// 		"email":  user.Email,
// 		"iat":    time.Now().Unix(),
// 		"iss":    os.Getenv("APP_ENV"),
// 		"exp":    time.Now().Add(7 * 24 * time.Hour).Unix(),
// 		"roles":  user.Roles,
// 	})

// 	refreshToken := jwtHelper.GenerateToken(jwtRefreshClaims, a.cfg.JWTConfig.RefreshSecretKey)

// 	tokens := Tokens{
// 		AccessToken:  accessToken,
// 		RefreshToken: refreshToken,
// 	}
// 	response.RespondWithJson(c, http.StatusOK, tokens)

// }

// func (a *authHandler) VerifyAccessToken(c *gin.Context) {
// 	token := c.GetHeader("Authorization")
// 	decodedToken := jwtHelper.VerifyToken(token, a.cfg.JWTConfig.SecretKey)
// 	response.RespondWithJson(c, http.StatusOK, decodedToken)

// }
// func (a *authHandler) VerifyRefreshToken(c *gin.Context) {
// 	token := c.GetHeader("Authorization")
// 	decodedToken := jwtHelper.VerifyToken(token, a.cfg.JWTConfig.RefreshSecretKey)
// 	response.RespondWithJson(c, http.StatusOK, decodedToken)

// }

// // TODO fetchroles functionnı
