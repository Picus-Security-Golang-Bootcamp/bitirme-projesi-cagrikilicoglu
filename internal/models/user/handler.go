package user

import (
	"net/http"

	"github.com/cagrikilicoglu/shopping-basket/internal/api"
	"github.com/cagrikilicoglu/shopping-basket/internal/models/response"
	"github.com/gin-gonic/gin"
	"github.com/go-openapi/strfmt"
	"go.uber.org/zap"
)

type userHandler struct {
	repo *UserRepository
}

func NewUserHandler(r *gin.RouterGroup, repo *UserRepository) {
	h := &userHandler{repo: repo}

	r.POST("/signup", h.createUser)

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
