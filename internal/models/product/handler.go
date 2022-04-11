package product

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/cagrikilicoglu/shopping-basket/internal/api"
	"github.com/cagrikilicoglu/shopping-basket/internal/models"
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
	// r.Use(middleware.AuthMiddleware(cfg.JWTConfig.SecretKey))
	h := &productHandler{repo: repo}
	r.GET("/", h.getAll)
	r.POST("/create", middleware.AdminAuthMiddleware(cfg.JWTConfig.SecretKey), h.create)
	r.POST("/upload", middleware.AdminAuthMiddleware(cfg.JWTConfig.SecretKey), h.createFromFile)
	r.GET("/id/:id", h.getByID)
	r.GET("/sku/:sku", h.getBySKU)
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

	response.RespondWithJson(c, http.StatusCreated, productToResponse(product))
	// c.JSON(http.StatusOK, productsToResponse(products))
}

// TODO
func (p *productHandler) getByID(c *gin.Context) {
	id := c.Param("id")

	product, err := p.repo.getByID(id)
	if err != nil {
		response.RespondWithError(c, err)
		return
	}
	response.RespondWithJson(c, http.StatusOK, productToResponse(product))
}

func (p *productHandler) createFromFile(c *gin.Context) {
	data, err := c.FormFile("file")
	if err != nil {
		response.RespondWithError(c, err)
		return
	}
	// TODO content type'ı check et
	// fileType := data.Header.Get("Content-Type")
	// if fileType != "application/CSV" {
	// 	response.RespondWithError(c, errors.New("wrong file type"))
	// 	return
	// }
	results, err := readProductsWithWorkerPool(data)
	if err != nil {
		response.RespondWithError(c, errors.New("file cannot be read"))
	}
	// TODO hangilerinin succesful hangilerin unsuccesfull olduğunu ekle
	var successfulCategories []models.Product
	var unsuccessfulCategories []models.Product
	for i, v := range results {
		product, err := p.repo.Create(&results[i])
		if err != nil {
			// tempCat := results[i]
			unsuccessfulCategories = append(unsuccessfulCategories, v)
			continue
		}
		successfulCategories = append(successfulCategories, *product)
	}
	if len(successfulCategories) > 0 {
		response.RespondWithJson(c, http.StatusCreated, productsToResponse(&successfulCategories))
	}
	//TODO burası pointer arrayi dönüyor Mutlaka bak
	if len(unsuccessfulCategories) > 0 {
		zap.L().Debug("unsuccessfulCategories", zap.Reflect("uns", unsuccessfulCategories))
		// responseCat := categoriesToResponse(&unsuccessfulCategories)
		response.RespondWithError(c, fmt.Errorf("Products %v already exists", productsToResponse(&unsuccessfulCategories)))
	}

}

//---------------BURADA KALDIN--------
func (p *productHandler) getBySKU(c *gin.Context) {
	sku := c.Param("sku")

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
