package category

import (
	"github.com/cagrikilicoglu/shopping-basket/internal/api"
	"github.com/cagrikilicoglu/shopping-basket/internal/models"
)

func responseToCategory(ac *api.Category) *models.Category {
	// TODO sil
	// stockNum := uint(*ap.Stock.Number)
	return &models.Category{
		Name:        ac.Name,
		Description: ac.Description,
	}

}

func categoriesToResponse(cs *[]models.Category) []*api.Category {
	categories := make([]*api.Category, 0)
	for _, c := range *cs {
		categories = append(categories, categoryToResponse(&c))
	}
	return categories
}

func categoryToResponse(p *models.Category) *api.Category {

	return &api.Category{
		Name:        p.Name,
		Description: p.Description,
	}

}
