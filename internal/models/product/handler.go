package product

import (
	"errors"
	"net/http"

	"github.com/cagrikilicoglu/shopping-basket/internal/api"
	"github.com/cagrikilicoglu/shopping-basket/internal/models/response"
	"github.com/cagrikilicoglu/shopping-basket/pkg/pagination"
	"github.com/gin-gonic/gin"
	"github.com/go-openapi/strfmt"
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
	r.POST("/create", h.create)
	r.GET("/:id", h.getByID)
	// r.GET("", h.getBySKU)
	r.GET("", h.getByName)
}

func (p *productHandler) getAll(c *gin.Context) {

	pageIndex, pageSize := pagination.GetPaginationParametersFromRequest(c)

	products, count, err := p.repo.getAll(pageIndex, pageSize)
	if err != nil {
		response.RespondWithError(c, err)
		return
	}
	paginatedResult := pagination.NewFromGinRequest(c, count)
	paginatedResult.Items = productsToResponse(products)
	c.Header("Page Links", paginatedResult.BuildLinkHeader(c.Request.URL.Path, pageSize))
	response.RespondWithJson(c, http.StatusOK, paginatedResult)
	// c.JSON(http.StatusOK, productsToResponse(products))
}

func (p *productHandler) create(c *gin.Context) {
	productBody := &api.Product{}

	if err := c.Bind(&productBody); err != nil {
		response.RespondWithError(c, err)
		return
	}
	if err := productBody.Validate(strfmt.NewFormats()); err != nil {
		response.RespondWithError(c, err)
		return
	}

	product, err := p.repo.Create(responseToProduct(productBody))
	if err != nil {
		response.RespondWithError(c, err)
	}

	response.RespondWithJson(c, http.StatusCreated, product)
	// c.JSON(http.StatusOK, productsToResponse(products))
}

func (p *productHandler) getByID(c *gin.Context) {
	id := c.Param("id")

	product, err := p.repo.getByID(id)
	if err != nil {
		response.RespondWithError(c, err)
		return
	}
	response.RespondWithJson(c, http.StatusOK, productToResponse(product))
}

// func (p *productHandler) getByQuery(c *gin.Context) {
// 	queryParams := c.Request.URL.Query()
// 	sku := queryParams["sku"]
// 	if sku != "" {

// 	}
// }

//---------------BURADA KALDIN--------
func (p *productHandler) getBySKU(c *gin.Context) {
	sku, ok := c.GetQuery("sku")
	if !ok {
		response.RespondWithError(c, errors.New("not Found"))
		return
	}
	product, err := p.repo.getBySKU(sku)
	if err != nil {
		response.RespondWithError(c, err)
		return
	}
	response.RespondWithJson(c, http.StatusOK, productToResponse(product))
}

func (p *productHandler) getByName(c *gin.Context) {
	name := c.Query("name")
	// TODO  arrayın uzunluğu 0sa error ver
	// if !ok {
	// 	response.RespondWithError(c, errors.New("not Found"))
	// 	return
	// }

	products, err := p.repo.getByName(name)
	if len(*products) == 0 {
		response.RespondWithError(c, errors.New("not Found"))
		return

	}
	if err != nil {
		response.RespondWithError(c, err)
		return
	}
	response.RespondWithJson(c, http.StatusOK, productsToResponse(products))
}

// // respondWithJson: creates responses to the request in a standardized structure
// func respondWithJson(c *gin.Context, code int, payload interface{}) {
// 	// data := ApiResponse{
// 	// 	Payload: payload,
// 	// }
// 	// response, err := json.Marshal(data)
// 	// if err != nil {
// 	// 	// respondWithError(w, httpErrors.ParseErrors(err))
// 	// 	return
// 	// }

// 	codeStr := strconv.Itoa(code) // TODO daha iyi handle et
// 	c.Header("code", codeStr)
// 	c.JSON(code, payload)
// }

// // respondWithError: creates responses when an error occurs in a standardized structure
// func respondWithError(c *gin.Context, err error) {
// 	a := httpErrors.ParseErrors(err)
// 	respondWithJson(c, a.Status(), a.Error())
// }
