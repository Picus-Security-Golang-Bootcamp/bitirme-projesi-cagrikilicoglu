package product

import (
	"fmt"

	"github.com/cagrikilicoglu/shopping-basket/internal/models"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{db: db}
}
func (pr *ProductRepository) updateBySKU(sku string, p *models.Product) (*models.Product, error) {
	zap.L().Debug("product.repo.update", zap.Reflect("product", p))

	// result := pr.db.Model(models.Product{}).Where("sku = ?", sku)

	// if result.Error != nil {
	// 	zap.L().Error("product.repo.Update failed to update product", zap.Error(result.Error))
	// 	return nil, result.Error
	// }
	// resultUpdates := result.Updates(p)
	// zap.L().Debug("product.repo.update", zap.Reflect("resultUpdates", result))
	if err := pr.db.Model(models.Product{}).Where("sku = ?", sku).Updates(p).First(p).Error; err != nil {
		zap.L().Error("product.repo.Update failed to update product", zap.Error(err))
		return nil, err
	}
	return p, nil
}

func (pr *ProductRepository) Create(p *models.Product) (*models.Product, error) {
	zap.L().Debug("product.repo.create", zap.Reflect("product", p))

	if err := pr.db.Create(p).Error; err != nil {
		zap.L().Error("product.repo.Create failed to create product", zap.Error(err))
		return nil, err
	}
	return p, nil
}

// TODO getirirken order olsun
func (pr *ProductRepository) getAll(pageIndex, pageSize int) (*[]models.Product, int, error) {
	zap.L().Debug("product.repo.getAll")

	var products *[]models.Product
	count, err := pr.GetCount()
	if err != nil {
		return nil, -1, err
	}

	// // TODO BURAYı getcount gibi bir fonksiyonla handle edebiliriz.
	// if err := pr.db.Find(&products).Count(&count).Error; err != nil {
	// 	zap.L().Error("product.repo.getAll failed to get products count", zap.Error(err))
	// 	return nil, -1, err
	// }
	// TODO handle -1
	if err := pr.db.Offset((pageIndex - 1) * pageSize).Limit(pageSize).Find(&products).Error; err != nil {
		zap.L().Error("product.repo.getAll failed to get products", zap.Error(err))
		return nil, -1, err
	}
	// result := pr.db.Offset((pageIndex - 1) * pageSize).Limit(pageSize).Find(&products).Count(&count)
	// if err := result.Error; err != nil {
	// 	zap.L().Error("product.repo.getAll failed to get products", zap.Error(err))
	// 	return nil, -1, err
	// }

	zap.L().Debug("count")
	zap.Reflect("count", count)
	fmt.Println(count)
	return products, count, nil
}

// func (pr *ProductRepository) delete(p *models.Product) error {
// 	result := pr.db.Delete(p)
// 	if result.Error != nil {
// 		return result.Error
// 	}

// 	return nil

// }
func (pr *ProductRepository) GetCount() (int, error) {
	var count int64
	var products *[]models.Product
	// TODO BURAYı getcount gibi bir fonksiyonla handle edebiliriz.
	if err := pr.db.Find(&products).Count(&count).Error; err != nil {
		zap.L().Error("product.repo.getAll failed to get products count", zap.Error(err))
		return -1, err
	}
	return int(count), nil
}

func (pr *ProductRepository) getByID(id string) (*models.Product, error) {
	zap.L().Debug("product.repo.getByID", zap.Reflect("id", id))

	var product *models.Product
	if result := pr.db.First(&product, id); result.Error != nil {
		zap.L().Error("product.repo.getByID failed to get products", zap.Error(result.Error))
		return nil, result.Error
	}
	return product, nil
}

func (pr *ProductRepository) GetBySKU(sku string) (*models.Product, error) {
	zap.L().Debug("product.repo.GetBySKU", zap.Reflect("SKU", sku))

	var product *models.Product
	result := pr.db.Where(&models.Product{Stock: models.Stock{SKU: sku}}).Find(&product)
	if result.Error != nil {
		zap.L().Error("product.repo.GetBySKU failed to get products", zap.Error(result.Error))
		return nil, result.Error
	}
	return product, nil
}

func (pr *ProductRepository) getByName(name string) (*[]models.Product, error) {
	zap.L().Debug("product.repo.getByName", zap.Reflect("name", name))

	var products = []models.Product{}
	result := pr.db.Where("Name ILIKE ?", "%"+name+"%").Find(&products)
	fmt.Println(result)
	if result.Error != nil {
		zap.L().Error("product.repo.getByName failed to get products", zap.Error(result.Error))
		return nil, result.Error
	}
	return &products, nil
}

// func (pr *ProductRepository) update(p *models.Product) (*models.Product, error) {
// 	zap.L().Debug("product.repo.update", zap.Reflect("product", p))

// 	if result := pr.db.Save(&p); result.Error != nil {
// 		return nil, result.Error
// 	}

// 	return p, nil
// }

func (pr *ProductRepository) deleteBySKU(sku string) error {
	zap.L().Debug("product.repo.delete", zap.Reflect("sku", sku))

	// product, err := pr.GetBySKU(sku)
	// if err != nil {
	// 	return err
	// }

	if result := pr.db.Where("sku = ?", sku).Delete(&models.Product{}); result.Error != nil {
		return result.Error
	}

	return nil
}

func (pr *ProductRepository) Migration() {
	pr.db.AutoMigrate(&models.Product{})
}

// TODO içine getbysku fonksiyonu alabilir
func (pr *ProductRepository) GetIDBySKU(sku string) (uuid.UUID, error) {
	zap.L().Debug("product.repo.GetIDBySKU", zap.Reflect("SKU", sku))

	product, err := pr.GetBySKU(sku)
	// var product *models.Product
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
	product, err := pr.CheckStock(sku, quantity)
	if err != nil {
		zap.L().Error("product.repo.Checkstock failed to get product", zap.Error(err))
		return err
	}
	//TODO aşağıdaki çalışmıyor olabilir
	if result := pr.db.Model(&product).Select("number").Update("number", (product.Stock.Number - quantity)); result.Error != nil {
		zap.L().Error("product.repo.UpdateStock failed to update product", zap.Error(err))
		return result.Error
	}
	return nil

}
