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

func (ur *UserRepository) Create(u *models.User) (*models.User, error) {
	zap.L().Debug("User.repo.create", zap.Reflect("User", u))

	if err := ur.db.Create(u).Error; err != nil {
		zap.L().Error("User.repo.Create failed to create User", zap.Error(err))
		return nil, err
	}
	return u, nil
}

func (ur *UserRepository) GetUser(email, password string) (*models.User, error) {
	zap.L().Debug("User.repo.getUser", zap.Reflect("email", email))
	zap.Reflect("password", password)

	var user *models.User
	if err := ur.db.Where(&models.User{Email: &email, Password: &password}).First(&user).Error; err != nil {
		zap.L().Error("User.repo.getUser failed to get User", zap.Error(err))
		return nil, err
	}
	return user, nil
}
