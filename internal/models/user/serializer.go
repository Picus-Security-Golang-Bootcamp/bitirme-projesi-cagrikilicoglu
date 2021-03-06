package user

import (
	"github.com/cagrikilicoglu/shopping-basket/internal/api"
	"github.com/cagrikilicoglu/shopping-basket/internal/models"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

var userRole = "user"

/// responseToUser converts user response model to database model
func responseToUser(u *api.User) (*models.User, error) {

	zap.L().Debug("User.serializer.responseToUser", zap.Reflect("user", u))
	encryptedPassword, err := getHash([]byte(*u.Password))
	if err != nil {
		return nil, err
	}

	return &models.User{
		Email:     u.Email,
		Password:  encryptedPassword,
		FirstName: *u.FirstName,
		LastName:  *u.LastName,
		ZipCode:   *u.ZipCode,
		Role:      userRole,
	}, nil
}

/// getHash encrypts password given by the user for security
func getHash(pwd []byte) (*string, error) {

	zap.L().Debug("User.serializer.getHash")
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		zap.L().Error("User.serializer.getHash cannot generate hash from password")
		return nil, err
	}

	hashStr := string(hash)
	return &hashStr, nil
}
