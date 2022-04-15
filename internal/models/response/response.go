package response

import (
	"strconv"

	"github.com/cagrikilicoglu/shopping-basket/internal/httpErrors"
	"github.com/gin-gonic/gin"
)

// respondWithJson: creates responses to the request in a standardized structure
func RespondWithJson(c *gin.Context, code int, payload interface{}) {
	c.Header("code", strconv.Itoa(code))
	c.JSON(code, payload)
}

// respondWithError: creates responses when an error occurs in a standardized structure
func RespondWithError(c *gin.Context, err error) {
	a := httpErrors.ParseErrors(err)
	RespondWithJson(c, a.Status(), a.Error())
}
