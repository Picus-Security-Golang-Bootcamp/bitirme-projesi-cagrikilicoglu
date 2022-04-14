package user

import (
	"errors"
	"net/http"
	"net/mail"

	"github.com/cagrikilicoglu/shopping-basket/internal/api"
	"github.com/cagrikilicoglu/shopping-basket/internal/models/response"
	"github.com/cagrikilicoglu/shopping-basket/pkg/auth"
	"github.com/cagrikilicoglu/shopping-basket/pkg/middleware"
	"github.com/gin-gonic/gin"
	"github.com/go-openapi/strfmt"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type userHandler struct {
	repo *UserRepository
	auth *auth.Authenticator
}

func NewUserHandler(r *gin.RouterGroup, repo *UserRepository, auth *auth.Authenticator) {
	h := &userHandler{repo: repo,
		auth: auth}

	r.POST("/signup", h.create)
	r.POST("/login", h.login)
	r.POST("/refresh", middleware.RefreshMiddleware(h.auth.Cfg.JWTConfig.RefreshSecretKey), h.Refresh)
}

func (u *userHandler) create(c *gin.Context) {

	zap.L().Debug("User.handler.create")

	userBody := &api.User{}

	zap.L().Debug("User.handler.create.Bind", zap.Reflect("userBody", userBody))
	if err := c.Bind(&userBody); err != nil {
		response.RespondWithError(c, err)
		return
	}
	zap.L().Debug("User.handler.create.Validate")
	if err := userBody.Validate(strfmt.NewFormats()); err != nil {
		response.RespondWithError(c, err)
		return
	}

	zap.L().Debug("User.handler.create.validateEmail")
	ok := validateEmail(*userBody.Email)
	if !ok {
		zap.L().Error("User.handler.create.validateEmail invalid email", zap.Reflect("email", *userBody.Email))
		response.RespondWithError(c, errors.New("Email is not valid"))
		return
	}

	userSerialized, err := responseToUser(userBody)
	if err != nil {
		response.RespondWithError(c, err)
		return
	}

	user, err := u.repo.Create(userSerialized)
	if err != nil {
		response.RespondWithError(c, err)
		return
	}

	tokens, err := u.auth.Authenticate(user.ID, *user.Email, user.Role)
	if err != nil {
		response.RespondWithError(c, errors.New("User cannot be authenticated"))
		return
	}

	response.RespondWithJson(c, http.StatusCreated, *tokens)
}

func (u *userHandler) login(c *gin.Context) {
	zap.L().Debug("User.handler.login")
	loginBody := api.Login{}

	zap.L().Debug("User.handler.login.Bind", zap.Reflect("loginBody", loginBody))
	if err := c.Bind(&loginBody); err != nil {
		response.RespondWithError(c, err)
		return
	}
	zap.L().Debug("User.handler.login.Bind", zap.Reflect("loginBody", loginBody))
	if err := loginBody.Validate(strfmt.NewFormats()); err != nil {
		response.RespondWithError(c, err)
		return
	}

	user, err := u.repo.get(*loginBody.Email)
	if err != nil {
		response.RespondWithError(c, errors.New("Wrong credentials"))
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(*user.Password), []byte(*loginBody.Password))
	if err != nil {
		zap.L().Error("User.handler.login.CompareHashAndPassword not matched with the given password", zap.Error(err))
		response.RespondWithError(c, errors.New("Wrong credentials"))
		return
	}

	tokens, err := u.auth.Authenticate(user.ID, *user.Email, user.Role)
	if err != nil {
		response.RespondWithError(c, errors.New("User cannot be authenticated"))
		return
	}
	response.RespondWithJson(c, http.StatusOK, *tokens)

}

func (u *userHandler) Refresh(c *gin.Context) {

	userID, _ := c.Get("userID")
	email, _ := c.Get("email")
	role, _ := c.Get("role")

	zap.L().Debug("User.handler.Refresh", zap.Reflect("userID", userID), zap.Reflect("email", email), zap.Reflect("role", role))

	userIDParsed, err := uuid.Parse(userID.(string))
	if err != nil {
		zap.L().Error("User.handler.Refresh.Parse cannot parse user ID", zap.Error(err))
		response.RespondWithError(c, err)
		return
	}

	tokens, err := u.auth.Authenticate(userIDParsed, email.(string), role.(string))
	if err != nil {
		response.RespondWithError(c, err)
		return
	}
	response.RespondWithJson(c, http.StatusOK, *tokens)
}

func validateEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}
