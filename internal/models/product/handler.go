package product

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/cagrikilicoglu/shopping-basket/internal/api"
	"github.com/cagrikilicoglu/shopping-basket/internal/models/response"
	"github.com/cagrikilicoglu/shopping-basket/pkg/config"
	"github.com/cagrikilicoglu/shopping-basket/pkg/middleware"
	"github.com/cagrikilicoglu/shopping-basket/pkg/pagination"
	"github.com/gin-gonic/gin"
	"github.com/go-openapi/strfmt"
	"go.uber.org/zap"
)

type productHandler struct {
	repo *ProductRepository
}

func NewProductHandler(r *gin.RouterGroup, repo *ProductRepository, cfg *config.Config) {

	h := &productHandler{repo: repo}
	r.GET("/", h.getAll)
	r.GET("/id/:id", h.getByID)
	r.GET("/sku/:sku", h.getBySKU)
	r.POST("/create", middleware.AdminAuthMiddleware(cfg.JWTConfig.SecretKey), h.create)
	r.POST("/upload", middleware.AdminAuthMiddleware(cfg.JWTConfig.SecretKey), h.createFromFile)
	r.GET("", h.getByName)
	r.DELETE("/delete/sku/:sku", middleware.AdminAuthMiddleware(cfg.JWTConfig.SecretKey), h.deleteBySKU)
	r.PUT("/update/sku/:sku", middleware.AdminAuthMiddleware(cfg.JWTConfig.SecretKey), h.updateBySKU)
}

// getAll fetches all the products in the database and paginate the results
func (p *productHandler) getAll(c *gin.Context) {

	pageIndex, pageSize := pagination.GetPaginationParametersFromRequest(c)
	zap.L().Debug("product.handler.getAll with pagination", zap.Reflect("pageIndex", pageIndex), zap.Reflect("pageSize", pageSize))

	products, count, err := p.repo.getAll(pageIndex, pageSize)
	if err != nil {
		response.RespondWithError(c, err)
		return
	}
	paginatedResult := pagination.NewFromGinRequest(c, count, ProductsToResponse(products))

	response.RespondWithJson(c, http.StatusOK, paginatedResult)
}

// getByID fetches a product by ID
func (p *productHandler) getByID(c *gin.Context) {

	id := c.Param("id")
	zap.L().Debug("product.handler.getByID", zap.Reflect("id", id))

	product, err := p.repo.getByID(id)
	if err != nil {
		response.RespondWithError(c, err)
		return
	}

	response.RespondWithJson(c, http.StatusOK, ProductToResponse(product))
}

// getBySKU fetches a product by SKU
func (p *productHandler) getBySKU(c *gin.Context) {

	sku := c.Param("sku")
	zap.L().Debug("product.handler.getBySKU", zap.Reflect("sku", sku))

	product, err := p.repo.GetBySKU(sku)
	if err != nil {
		response.RespondWithError(c, err)
		return
	}
	response.RespondWithJson(c, http.StatusOK, ProductToResponse(product))
}

// create creates a product by the input in request body
func (p *productHandler) create(c *gin.Context) {
	zap.L().Debug("product.handler.create")
	productBody := &api.Product{}

	zap.L().Debug("product.handler.create.Bind", zap.Reflect("productBody", productBody))
	if err := c.Bind(&productBody); err != nil {
		response.RespondWithError(c, err)
		return
	}

	zap.L().Debug("product.handler.create.Validate")
	if err := productBody.Validate(strfmt.NewFormats()); err != nil {
		response.RespondWithError(c, err)
		return
	}

	product, err := p.repo.create(responseToProduct(productBody))
	if err != nil {
		response.RespondWithError(c, err)
	}

	response.RespondWithJson(c, http.StatusCreated, ProductToResponseForAdmin(product))
}

// createFromFile reads data from a csv file and create products from it
func (p *productHandler) createFromFile(c *gin.Context) {
	zap.L().Debug("product.handler.createFromFile")
	data, err := c.FormFile("file")
	if err != nil {
		response.RespondWithError(c, err)
		return
	}

	results, err := readProductsWithWorkerPool(data)
	if err != nil {
		response.RespondWithError(c, errors.New("file cannot be read"))
		return
	}

	products, err := p.repo.batchCreate(results)
	if err != nil {
		response.RespondWithError(c, err)
		return
	}

	response.RespondWithJson(c, http.StatusCreated, productsToResponseForAdmin(&products))
}

// getByName fetches products by name
func (p *productHandler) getByName(c *gin.Context) {
	name := c.Query("name")
	zap.L().Debug("product.handler.getByName", zap.Reflect("name", name))

	products, err := p.repo.getByName(name)
	if len(*products) == 0 {
		response.RespondWithError(c, errors.New("not found"))
		return
	}

	if err != nil {
		response.RespondWithError(c, err)
		return
	}
	response.RespondWithJson(c, http.StatusOK, ProductsToResponse(products))
}

// deleteBySKU deletes a product by SKU
func (p *productHandler) deleteBySKU(c *gin.Context) {
	sku := c.Param("sku")
	zap.L().Debug("product.handler.deleteBySKU", zap.Reflect("sku", sku))

	err := p.repo.deleteBySKU(sku)
	if err != nil {
		response.RespondWithError(c, err)
		return
	}
	response.RespondWithJson(c, http.StatusOK, fmt.Sprintf("Product successfully deleted"))
}

// updateBySKU updates a product by SKU
func (p *productHandler) updateBySKU(c *gin.Context) {
	sku := c.Param("sku")
	zap.L().Debug("product.handler.updateBySKU", zap.Reflect("sku", sku))

	productBody := &api.Product{}

	zap.L().Debug("product.handler.updateBySKU.Bind", zap.Reflect("productBody", productBody))
	if err := c.Bind(&productBody); err != nil {
		response.RespondWithError(c, err)
		return
	}
	zap.L().Debug("product.handler.updateBySKU.Validate")
	if err := productBody.Validate(strfmt.NewFormats()); err != nil {
		response.RespondWithError(c, err)
		return
	}

	product, err := p.repo.updateBySKU(sku, responseToProduct(productBody))
	if err != nil {
		response.RespondWithError(c, err)
		return
	}

	response.RespondWithJson(c, http.StatusOK, ProductToResponseForAdmin(product))

}
