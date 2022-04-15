package category

import (
	"github.com/cagrikilicoglu/shopping-basket/internal/api"
	"github.com/cagrikilicoglu/shopping-basket/internal/models"
	"go.uber.org/zap"
)

func responseToCategory(ac *api.Category) *models.Category {
	zap.L().Debug("Category.serializer.responseToCategory", zap.Reflect("apiCategories", ac))

	return &models.Category{
		Name:        ac.Name,
		Description: ac.Description,
	}
}

func categoriesToResponse(cs *[]models.Category) []*api.Category {
	zap.L().Debug("Category.serializer.categoriesToResponse", zap.Reflect("Categories", cs))

	categories := make([]*api.Category, 0)
	for i := range *cs {
		categoriesDeref := *cs
		categories = append(categories, categoryToResponse(&categoriesDeref[i]))
	}
	return categories
}

func categoryToResponse(p *models.Category) *api.Category {
	zap.L().Debug("Category.serializer.categoriesToResponse", zap.Reflect("Categories", p))

	return &api.Category{
		Name:        p.Name,
		Description: p.Description,
	}
}
