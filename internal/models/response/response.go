package response

import (
	"strconv"

	"github.com/cagrikilicoglu/shopping-basket/internal/httpErrors"
	"github.com/gin-gonic/gin"
)

// RespondWithJson creates responses to the http requests in a standardized structure
func RespondWithJson(c *gin.Context, code int, payload interface{}) {
	c.Header("code", strconv.Itoa(code))
	c.JSON(code, payload)
}

// RespondWithError creates responsesto the http requests in a standardized structure when an error occurs
func RespondWithError(c *gin.Context, err error) {
	a := httpErrors.ParseErrors(err)
	RespondWithJson(c, a.Status(), a.Error())
}
