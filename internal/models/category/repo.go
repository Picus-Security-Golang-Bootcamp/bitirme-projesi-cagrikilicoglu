package category

import (
	"github.com/cagrikilicoglu/shopping-basket/internal/models"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type CategoryRepository struct {
	db *gorm.DB
}

func (cr *CategoryRepository) Migration() {
	cr.db.AutoMigrate(&models.Category{})
}

func NewCategoryRepository(db *gorm.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (cr *CategoryRepository) Create(c *models.Category) (*models.Category, error) {
	zap.L().Debug("Category.repo.create", zap.Reflect("Category", c))

	if err := cr.db.Create(c).Error; err != nil {
		zap.L().Error("Category.repo.Create failed to create Category", zap.Error(err))
		return nil, err
	}
	return c, nil
}

func (cr *CategoryRepository) GetCount() (int, error) {
	var count int64
	var categories *[]models.Category
	// TODO BURAYı getcount gibi bir fonksiyonla handle edebiliriz.
	if err := cr.db.Find(&categories).Count(&count).Error; err != nil {
		zap.L().Error("category.repo.getCount failed to get categories count", zap.Error(err))
		return -1, err
	}
	return int(count), nil
}

// TODO getirirken order olsun
func (cr *CategoryRepository) getAll(pageIndex, pageSize int) (*[]models.Category, int, error) {
	zap.L().Debug("category.repo.getAll")

	var categories *[]models.Category
	count, err := cr.GetCount()
	if err != nil {
		return nil, -1, err
	}

	// // TODO BURAYı getcount gibi bir fonksiyonla handle edebiliriz.
	// if err := pr.db.Find(&products).Count(&count).Error; err != nil {
	// 	zap.L().Error("product.repo.getAll failed to get products count", zap.Error(err))
	// 	return nil, -1, err
	// }
	// TODO handle -1
	if err := cr.db.Offset((pageIndex - 1) * pageSize).Limit(pageSize).Find(&categories).Error; err != nil {
		zap.L().Error("category.repo.getAll failed to get categories", zap.Error(err))
		return nil, -1, err
	}
	// result := pr.db.Offset((pageIndex - 1) * pageSize).Limit(pageSize).Find(&products).Count(&count)
	// if err := result.Error; err != nil {
	// 	zap.L().Error("product.repo.getAll failed to get products", zap.Error(err))
	// 	return nil, -1, err
	// }

	// zap.L().Debug("count")
	// zap.Reflect("count", count)
	// fmt.Println(count)
	return categories, count, nil
}

func (cr *CategoryRepository) getByNameWithProducts(name string) (*models.Category, error) {
	zap.L().Debug("category.repo.getByName", zap.Reflect("name", name))

	var category *models.Category
	if result := cr.db.Preload("Products").Where("name = ?", name).First(&category); result.Error != nil {
		zap.L().Error("category.repo.getByName failed to get products", zap.Error(result.Error))
		return nil, result.Error
	}
	zap.L().Debug("category.repo.getByName", zap.Reflect("name", category.Products))
	return category, nil
}
