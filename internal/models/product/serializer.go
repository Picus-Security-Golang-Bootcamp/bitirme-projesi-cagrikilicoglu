package product

import (
	"github.com/cagrikilicoglu/shopping-basket/internal/api"
	"github.com/cagrikilicoglu/shopping-basket/internal/models"
	"go.uber.org/zap"
)

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

func ProductsToResponse(ps *[]models.Product) []*api.Product {
	zap.L().Debug("Product.serializer.productsToResponse", zap.Reflect("Products", ps))

	products := make([]*api.Product, 0)
	for i := range *ps {
		productsDeref := *ps
		products = append(products, ProductToResponse(&productsDeref[i]))
	}
	return products
}
func productsToResponseForAdmin(ps *[]models.Product) []*api.Product {
	zap.L().Debug("Product.serializer.productsToResponseForAdmin", zap.Reflect("Products", ps))

	products := make([]*api.Product, 0)
	for i := range *ps {
		productsDeref := *ps
		products = append(products, ProductToResponseForAdmin(&productsDeref[i]))
	}
	return products
}

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
