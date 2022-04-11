package category

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/cagrikilicoglu/shopping-basket/internal/api"
	"github.com/cagrikilicoglu/shopping-basket/internal/models"
	"github.com/cagrikilicoglu/shopping-basket/internal/models/response"
	"github.com/cagrikilicoglu/shopping-basket/pkg/pagination"
	"go.uber.org/zap"

	// "github.com/cagrikilicoglu/shopping-basket/pkg/csvHelper"
	"github.com/gin-gonic/gin"
	"github.com/go-openapi/strfmt"
)

type categoryHandler struct {
	repo *CategoryRepository
}

// type ApiResponse struct {
// 	Payload interface{} `json:"data"`
// }

func NewCategoryHandler(r *gin.RouterGroup, repo *CategoryRepository) {
	h := &categoryHandler{repo: repo}
	r.POST("/create", h.create)
	r.POST("/upload", h.createFromFile)
	r.GET("/", h.getAll)
	// r.POST("/create", h.create)
	// r.GET("/:id", h.getByID)
	// // r.GET("", h.getBySKU)
	// r.GET("", h.getByName)
}
func (ch *categoryHandler) getAll(c *gin.Context) {
	pageIndex, pageSize := pagination.GetPaginationParametersFromRequest(c)

	categories, count, err := ch.repo.getAll(pageIndex, pageSize)
	if err != nil {
		response.RespondWithError(c, err)
		return
	}
	paginatedResult := pagination.NewFromGinRequest(c, count)
	paginatedResult.Items = categoriesToResponse(categories)
	c.Header("Page Links", paginatedResult.BuildLinkHeader(c.Request.URL.Path, pageSize))
	response.RespondWithJson(c, http.StatusOK, paginatedResult)
}

func (ch *categoryHandler) createFromFile(c *gin.Context) {
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
	results, err := readCategoriesWithWorkerPool(data)
	if err != nil {
		response.RespondWithError(c, errors.New("file cannot be read"))
	}
	// TODO hangilerinin succesful hangilerin unsuccesfull olduğunu ekle
	var successfulCategories []models.Category
	var unsuccessfulCategories []models.Category
	for i, v := range results {
		category, err := ch.repo.Create(&results[i])
		if err != nil {
			// tempCat := results[i]
			unsuccessfulCategories = append(unsuccessfulCategories, v)
			continue
		}
		successfulCategories = append(successfulCategories, *category)
	}
	if len(successfulCategories) > 0 {
		response.RespondWithJson(c, http.StatusCreated, categoriesToResponse(&successfulCategories))
	}
	//TODO burası pointer arrayi dönüyor Mutlaka bak
	if len(unsuccessfulCategories) > 0 {
		zap.L().Debug("unsuccessfulCategories", zap.Reflect("uns", unsuccessfulCategories))
		// responseCat := categoriesToResponse(&unsuccessfulCategories)
		response.RespondWithError(c, fmt.Errorf("Categories %v already exists", categoriesToResponse(&unsuccessfulCategories)))
	}

}

func (ch *categoryHandler) create(c *gin.Context) {
	categoryBody := &api.Category{}

	if err := c.Bind(&categoryBody); err != nil {
		response.RespondWithError(c, err)
		return
	}
	if err := categoryBody.Validate(strfmt.NewFormats()); err != nil {
		response.RespondWithError(c, err)
		return
	}

	category, err := ch.repo.Create(responseToCategory(categoryBody))
	if err != nil {
		response.RespondWithError(c, err)
	}

	response.RespondWithJson(c, http.StatusCreated, category)
	// c.JSON(http.StatusOK, productsToResponse(products))
}
