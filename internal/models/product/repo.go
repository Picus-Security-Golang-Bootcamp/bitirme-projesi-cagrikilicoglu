package product

import (
	"fmt"

	"github.com/cagrikilicoglu/shopping-basket/internal/models"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (pr *ProductRepository) Migration() {
	pr.db.AutoMigrate(&models.Product{})
}

func (pr *ProductRepository) create(p *models.Product) (*models.Product, error) {
	zap.L().Debug("product.repo.create", zap.Reflect("product", p))

	if err := pr.db.Create(p).Error; err != nil {
		zap.L().Error("product.repo.Create failed to create product", zap.Error(err))
		return nil, err
	}
	return p, nil
}

func (pr *ProductRepository) batchCreate(ps []models.Product) ([]models.Product, error) {
	zap.L().Debug("product.repo.batchCreate", zap.Reflect("products", ps))

	if err := pr.db.Clauses(clause.OnConflict{DoNothing: true}).Preload("Categories").Create(&ps).Error; err != nil {
		zap.L().Error("product.repo.batchCreate failed to create product", zap.Error(err))
		return nil, err
	}
	return ps, nil
}

func (pr *ProductRepository) getAll(pageIndex, pageSize int) (*[]models.Product, int, error) {

	zap.L().Debug("product.repo.getAll")
	var products *[]models.Product
	var count int64

	if err := pr.db.Order("name").Offset((pageIndex - 1) * pageSize).Limit(pageSize).Find(&products).Offset(-1).Limit(-1).Count(&count).Error; err != nil {
		zap.L().Error("product.repo.getAll failed to get products", zap.Error(err))
		return nil, -1, err
	}
	return products, int(count), nil
}

func (pr *ProductRepository) getByID(id string) (*models.Product, error) {

	zap.L().Debug("product.repo.getByID", zap.Reflect("id", id))
	var product *models.Product

	if result := pr.db.First(&product, id); result.Error != nil {
		zap.L().Error("product.repo.getByID failed to get product", zap.Error(result.Error))
		return nil, result.Error
	}
	return product, nil
}

func (pr *ProductRepository) GetBySKU(sku string) (*models.Product, error) {

	zap.L().Debug("product.repo.GetBySKU", zap.Reflect("SKU", sku))
	var product *models.Product

	if result := pr.db.Where(&models.Product{Stock: models.Stock{SKU: sku}}).First(&product); result.Error != nil {
		zap.L().Error("product.repo.GetBySKU failed to get products", zap.Error(result.Error))
		return nil, result.Error
	}
	return product, nil
}

func (pr *ProductRepository) getByName(name string) (*[]models.Product, error) {
	zap.L().Debug("product.repo.getByName", zap.Reflect("name", name))

	var products = []models.Product{}

	if result := pr.db.Where("Name ILIKE ?", "%"+name+"%").Find(&products); result.Error != nil {
		zap.L().Error("product.repo.getByName failed to get products", zap.Error(result.Error))
		return nil, result.Error
	}
	return &products, nil
}
func (pr *ProductRepository) deleteBySKU(sku string) error {
	zap.L().Debug("product.repo.deleteBySKU", zap.Reflect("sku", sku))

	if result := pr.db.Where("sku = ?", sku).Delete(&models.Product{}); result.Error != nil {
		return result.Error
	}
	return nil
}

func (pr *ProductRepository) updateBySKU(sku string, p *models.Product) (*models.Product, error) {
	zap.L().Debug("product.repo.updateBySKU", zap.Reflect("product", p))

	if err := pr.db.Model(models.Product{}).Where("sku = ?", sku).Updates(p).First(p).Error; err != nil {
		zap.L().Error("product.repo.updateBySKU failed to update product", zap.Error(err))
		return nil, err
	}
	return p, nil
}

/////// ------- ////////
// TODO içine GetBySKU fonksiyonu alabilir
func (pr *ProductRepository) GetIDBySKU(sku string) (uuid.UUID, error) {
	zap.L().Debug("product.repo.GetIDBySKU", zap.Reflect("SKU", sku))

	product, err := pr.GetBySKU(sku)
	if err != nil {
		zap.L().Error("product.repo.GetIDBySKU failed to get product", zap.Error(err))
		return uuid.Nil, err
	}
	return product.ID, nil
}

func (pr *ProductRepository) CheckStock(sku string, quantity uint) (*models.Product, error) {

	zap.L().Debug("product.repo.CheckStock", zap.Reflect("SKU", sku))
	product, err := pr.GetBySKU(sku)
	if err != nil {
		zap.L().Error("product.repo.CheckStock failed to get product", zap.Error(err))
		return nil, err
	}
	if product.Stock.Number < quantity {
		return nil, fmt.Errorf("Not enough %s in the stock,please request less than %d", *product.Name, (product.Stock.Number + 1))
	}

	return product, nil
}

func (pr *ProductRepository) UpdateStock(sku string, quantity uint) error {

	if result := pr.db.Model(models.Product{}).Where("sku = ?", sku).Select("number").Update("number", gorm.Expr("number - ?", quantity)); result.Error != nil {
		zap.L().Error("product.repo.UpdateStock failed to update product", zap.Error(result.Error))
		return result.Error
	}
	return nil

}

// func (pr *ProductRepository) update(p *models.Product) (*models.Product, error) {
// 	zap.L().Debug("product.repo.update", zap.Reflect("product", p))

// 	if result := pr.db.Save(&p); result.Error != nil {
// 		return nil, result.Error
// 	}

// 	return p, nil
// }
