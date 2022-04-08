package product

import (
	"github.com/cagrikilicoglu/shopping-basket/internal/api"
	"github.com/cagrikilicoglu/shopping-basket/internal/models"
	"gorm.io/gorm"
)

func productToResponse(p *models.Product) *api.Product {

	stockNum := uint32(p.Stock.Number)
	idNum := uint32(p.ID)
	return &api.Product{
		ID:           &idNum,
		CategoryName: p.CategoryName,
		Name:         p.Name,
		Price:        &p.Price,
		Stock: &api.Stock{
			Number: &stockNum,
			Sku:    &p.Stock.SKU,
			Status: &p.Stock.Status,
		},
	}

}

func productsToResponse(ps *[]models.Product) []*api.Product {
	products := make([]*api.Product, 0)
	for _, p := range *ps {
		products = append(products, productToResponse(&p))
	}
	return products
}

func responseToProduct(ap *api.Product) *models.Product {
	stockNum := uint(*ap.Stock.Number)
	return &models.Product{
		Model: gorm.Model{ID: uint(*ap.ID)},
		Name:  ap.Name,
		Price: *ap.Price,
		Stock: models.Stock{
			SKU:    *ap.Stock.Sku,
			Number: stockNum,
			Status: *ap.Stock.Status,
		},
		CategoryName: ap.CategoryName,
	}

}
