package product

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type productHandler struct {
	repo *ProductRepository
}

type ApiResponse struct {
	Payload interface{} `json:"data"`
}

func NewProductHandler(r *gin.RouterGroup, repo *ProductRepository) {
	h := &productHandler{repo: repo}
	r.GET("/", h.getAll)
}

func (p *productHandler) getAll(c *gin.Context) {
	products, err := p.repo.GetAll()
	if err != nil {
		return // TODO respond with error ekle
	}
	respondWithJson(c, http.StatusOK, products)
}

// respondWithJson: creates responses to the request in a standardized structure
func respondWithJson(c *gin.Context, code int, payload interface{}) {
	data := ApiResponse{
		Payload: payload,
	}
	response, err := json.Marshal(data)
	if err != nil {
		// respondWithError(w, httpErrors.ParseErrors(err))
		return
	}

	codeStr := strconv.Itoa(code) // TODO daha iyi handle et
	c.Header("code", codeStr)
	c.JSON(http.StatusOK, response)
}

// respondWithError: creates responses when an error occurs in a standardized structure
// func respondWithError(w http.ResponseWriter, a httpErrors.ApiErr) {
// 	respondWithJson(w, a.Status(), a.Error())
// }
