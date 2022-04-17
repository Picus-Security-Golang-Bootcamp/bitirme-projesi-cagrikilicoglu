package category

import (
	"errors"
	"net/http"

	"github.com/cagrikilicoglu/shopping-basket/internal/api"
	"github.com/cagrikilicoglu/shopping-basket/internal/models/product"
	"github.com/cagrikilicoglu/shopping-basket/internal/models/response"
	"github.com/cagrikilicoglu/shopping-basket/pkg/config"
	"github.com/cagrikilicoglu/shopping-basket/pkg/middleware"
	"github.com/cagrikilicoglu/shopping-basket/pkg/pagination"
	"github.com/gin-gonic/gin"
	"github.com/go-openapi/strfmt"
	"go.uber.org/zap"
)

type categoryHandler struct {
	repo *CategoryRepository
}

func NewCategoryHandler(r *gin.RouterGroup, repo *CategoryRepository, cfg *config.Config) {
	h := &categoryHandler{repo: repo}

	r.GET("/", h.getAll)
	r.POST("/create", middleware.AdminAuthMiddleware(cfg.JWTConfig.SecretKey), h.create)
	r.POST("/upload", middleware.AdminAuthMiddleware(cfg.JWTConfig.SecretKey), h.createFromFile)
	r.GET("/:name", h.getProductsByCategoryName)
}

func (ch *categoryHandler) getAll(c *gin.Context) {

	pageIndex, pageSize := pagination.GetPaginationParametersFromRequest(c)
	zap.L().Debug("category.handler.getAll with pagination", zap.Reflect("pageIndex", pageIndex), zap.Reflect("pageSize", pageSize))

	categories, count, err := ch.repo.getAll(pageIndex, pageSize)
	if err != nil {
		response.RespondWithError(c, err)
		return
	}

	paginatedResult := pagination.NewFromGinRequest(c, count, categoriesToResponse(categories))

	response.RespondWithJson(c, http.StatusOK, paginatedResult)
}

func (ch *categoryHandler) getProductsByCategoryName(c *gin.Context) {

	name := c.Param("name")
	zap.L().Debug("category.handler.getByName", zap.Reflect("name", name))

	category, err := ch.repo.getByNameWithProducts(name)
	if err != nil {
		response.RespondWithError(c, err)
		return
	}
	// TODO prdouctları serialize etmek gerekir
	response.RespondWithJson(c, http.StatusOK, product.ProductsToResponse(&category.Products))
}

func (ch *categoryHandler) createFromFile(c *gin.Context) {

	zap.L().Debug("category.handler.createFromFile")
	data, err := c.FormFile("file")
	if err != nil {
		response.RespondWithError(c, err)
		return
	}

	results, err := readCategoriesWithWorkerPool(data)
	if err != nil {
		response.RespondWithError(c, errors.New("file cannot be read"))
		return
	}

	// TODO batch create var olanları göster eklenebilir.
	categories, err := ch.repo.batchCreate(results)
	if err != nil {
		response.RespondWithError(c, err)
		return
	}

	response.RespondWithJson(c, http.StatusCreated, categoriesToResponse(&categories))
}

func (ch *categoryHandler) create(c *gin.Context) {
	zap.L().Debug("category.handler.create")
	categoryBody := &api.Category{}

	zap.L().Debug("category.handler.create.Bind", zap.Reflect("categoryBody", categoryBody))
	if err := c.Bind(&categoryBody); err != nil {
		response.RespondWithError(c, err)
		return
	}
	zap.L().Debug("category.handler.create.Validate")
	if err := categoryBody.Validate(strfmt.NewFormats()); err != nil {
		response.RespondWithError(c, err)
		return
	}

	category, err := ch.repo.create(responseToCategory(categoryBody))
	if err != nil {
		response.RespondWithError(c, err)
		return
	}

	response.RespondWithJson(c, http.StatusCreated, categoryToResponse(category))
}
