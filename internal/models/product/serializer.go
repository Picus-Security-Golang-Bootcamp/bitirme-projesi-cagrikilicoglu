package product

import (
	"github.com/cagrikilicoglu/shopping-basket/internal/api"
	"github.com/cagrikilicoglu/shopping-basket/internal/models"
	"go.uber.org/zap"
)

// ProductToResponse converts product database model to response model
func ProductToResponse(p *models.Product) *api.Product {
	zap.L().Debug("Product.serializer.ProductToResponse", zap.Reflect("Products", p))
	return &api.Product{
		CategoryName: p.CategoryName,
		Name:         p.Name,
		Price:        &p.Price,
		Stock: &api.Stock{
			Sku: &p.Stock.SKU,
		},
	}
}

// ProductToResponseForAdmin converts product database model to response model for admin
// note that the result show also the stock number of a product
func ProductToResponseForAdmin(p *models.Product) *api.Product {
	zap.L().Debug("Product.serializer.ProductToResponseForAdmin", zap.Reflect("Products", p))

	stockNum := uint32(p.Stock.Number)
	return &api.Product{
		CategoryName: p.CategoryName,
		Name:         p.Name,
		Price:        &p.Price,
		Stock: &api.Stock{
			Number: stockNum,
			Sku:    &p.Stock.SKU,
		},
	}
}

/// ProductToResponse converts product database model to response model
func ProductsToResponse(ps *[]models.Product) []*api.Product {
	zap.L().Debug("Product.serializer.productsToResponse", zap.Reflect("Products", ps))

	products := make([]*api.Product, 0)
	for i := range *ps {
		productsDeref := *ps
		products = append(products, ProductToResponse(&productsDeref[i]))
	}
	return products
}

// productsToResponseForAdmin converts product database model to response model as a batch for admin
// note that the result show also the stock number of a product
func productsToResponseForAdmin(ps *[]models.Product) []*api.Product {
	zap.L().Debug("Product.serializer.productsToResponseForAdmin", zap.Reflect("Products", ps))

	products := make([]*api.Product, 0)
	for i := range *ps {
		productsDeref := *ps
		products = append(products, ProductToResponseForAdmin(&productsDeref[i]))
	}
	return products
}

// responseToProduct converts product response model to database model
func responseToProduct(ap *api.Product) *models.Product {
	zap.L().Debug("Product.serializer.responseToCategory", zap.Reflect("apiProducts", ap))

	stockNum := uint(ap.Stock.Number)
	return &models.Product{
		Name:  ap.Name,
		Price: *ap.Price,
		Stock: models.Stock{
			SKU:    *ap.Stock.Sku,
			Number: stockNum,
		},
		CategoryName: ap.CategoryName,
	}
}
