package user

import (
	"github.com/cagrikilicoglu/shopping-basket/internal/models"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func (ur *UserRepository) Migration() {
	ur.db.AutoMigrate(&models.User{})
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (pr *UserRepository) Create(u *models.User) (*models.User, error) {
	zap.L().Debug("User.repo.create", zap.Reflect("User", u))

	if err := pr.db.Create(u).Error; err != nil {
		zap.L().Error("User.repo.Create failed to create User", zap.Error(err))
		return nil, err
	}
	return u, nil
}
