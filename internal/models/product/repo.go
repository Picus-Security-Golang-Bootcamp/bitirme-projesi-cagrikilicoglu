package product

import (
	"github.com/cagrikilicoglu/shopping-basket/internal/models"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (pr *ProductRepository) Create(p *models.Product) (*models.Product, error) {
	zap.L().Debug("product.repo.create", zap.Reflect("product", p))

	if err := pr.db.Create(p).Error; err != nil {
		zap.L().Error("product.repo.Create failed to create product", zap.Error(err))
		return nil, err
	}
	return p, nil
}

func (pr *ProductRepository) getAll() (*[]models.Product, error) {
	zap.L().Debug("product.repo.getAll")

	var products *[]models.Product
	if err := pr.db.Find(&products).Error; err != nil {
		zap.L().Error("product.repo.getAll failed to get products", zap.Error(err))
		return nil, err
	}

	return products, nil
}

func (pr *ProductRepository) getByID(id string) (*models.Product, error) {
	zap.L().Debug("product.repo.getByID", zap.Reflect("id", id))

	var product *models.Product
	if result := pr.db.First(&product, product.Stock.SKU); result.Error != nil {
		return nil, result.Error
	}
	return product, nil
}

func (pr *ProductRepository) getBySKU(sku string) (*models.Product, error) {
	zap.L().Debug("product.repo.getBySKU", zap.Reflect("SKU", sku))

	var product *models.Product
	if result := pr.db.First(&product, product.Stock.SKU); result.Error != nil {
		return nil, result.Error
	}
	return product, nil
}

func (pr *ProductRepository) getByName(name string) ([]*models.Product, error) {
	zap.L().Debug("product.repo.getByName", zap.Reflect("name", name))

	var products = []*models.Product{}
	if result := pr.db.Where("name ILIKE ?", "%"+name+"%").Find(&products); result.Error != nil {
		return nil, result.Error
	}
	return products, nil
}

func (pr *ProductRepository) update(p *models.Product) (*models.Product, error) {
	zap.L().Debug("product.repo.update", zap.Reflect("product", p))

	if result := pr.db.Save(&p); result.Error != nil {
		return nil, result.Error
	}

	return p, nil
}

func (pr *ProductRepository) delete(sku string) error {
	zap.L().Debug("product.repo.delete", zap.Reflect("sku", sku))

	product, err := pr.getBySKU(sku)
	if err != nil {
		return err
	}

	if result := pr.db.Delete(&product); result.Error != nil {
		return result.Error
	}

	return nil
}

func (pr *ProductRepository) Migration() {
	pr.db.AutoMigrate(&models.Product{})
}
