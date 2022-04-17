package category

import (
	"github.com/cagrikilicoglu/shopping-basket/internal/models"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

// getAll fetches categories with pagination parameters from the database
func (cr *CategoryRepository) getAll(pageIndex, pageSize int) (*[]models.Category, int, error) {

	zap.L().Debug("category.repo.getAll")
	var categories *[]models.Category
	var count int64

	if err := cr.db.Order("name").Offset((pageIndex - 1) * pageSize).Limit(pageSize).Find(&categories).Offset(-1).Limit(-1).Count(&count).Error; err != nil {
		zap.L().Error("category.repo.getAll failed to get categories", zap.Error(err))
		return nil, -1, err
	}
	return categories, int(count), nil
}

// getByNameWithProducts fetches product data by category name from the database
func (cr *CategoryRepository) getByNameWithProducts(name string) (*models.Category, error) {
	zap.L().Debug("category.repo.getByNameWithProducts", zap.Reflect("name", name))

	var category *models.Category
	if result := cr.db.Preload("Products").Where("name = ?", name).First(&category); result.Error != nil {
		zap.L().Error("category.repo.getByNameWithProducts failed to get category", zap.Error(result.Error))
		return nil, result.Error
	}

	return category, nil
}

// create creates a category in the database
func (cr *CategoryRepository) create(c *models.Category) (*models.Category, error) {
	zap.L().Debug("Category.repo.create", zap.Reflect("Category", c))

	if err := cr.db.Create(c).Error; err != nil {
		zap.L().Error("Category.repo.Create failed to create Category", zap.Error(err))
		return nil, err
	}
	return c, nil
}

// batchCreate creates categories as a batch in the database
func (cr *CategoryRepository) batchCreate(cs []models.Category) ([]models.Category, error) {
	zap.L().Debug("Category.repo.batchCreate", zap.Reflect("Categories", cs))

	if err := cr.db.Clauses(clause.OnConflict{DoNothing: true}).Create(&cs).Error; err != nil {
		zap.L().Error("Category.repo.batchCreate failed to create Category", zap.Error(err))
		return nil, err
	}

	return cs, nil
}
