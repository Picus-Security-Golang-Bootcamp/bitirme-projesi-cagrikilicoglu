package user

import (
	"log"

	"github.com/cagrikilicoglu/shopping-basket/internal/api"
	"github.com/cagrikilicoglu/shopping-basket/internal/models"
	"golang.org/x/crypto/bcrypt"
)

func userToResponse(u *models.User) *api.User {

	return &api.User{
		Email:     u.Email,
		Password:  u.Password,
		FirstName: &u.FirstName,
		LastName:  &u.LastName,
		ZipCode:   &u.ZipCode,
	}

}

func responseToUser(u *api.User) *models.User {

	encryptedPassword := getHash([]byte(*u.Password))
	// var roles []models.Role
	// roles = append(roles, models.Role{Role: "user"})
	role := "user"
	return &models.User{

		Email:     u.Email,
		Password:  &encryptedPassword,
		FirstName: *u.FirstName,
		LastName:  *u.LastName,
		ZipCode:   *u.ZipCode,
		Role:      role,
		// Cart:      models.Cart{}, //TODO düzelt
	}

}

func getHash(pwd []byte) string {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.Println(err) // TODO başka şekilde handle et
	}
	return string(hash)
}
