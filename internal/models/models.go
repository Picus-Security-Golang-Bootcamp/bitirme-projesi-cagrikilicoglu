package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Stock struct {
	SKU    string `json:"sku" gorm:"unique"`
	Number uint   `json:"number"`
	Status string `json:"status"`
}

type Product struct {
	gorm.Model
	Name         *string `json:"description"`
	Price        float32 `json:"price"`
	Stock        Stock   `json:"stock" gorm:"embedded"`
	CategoryName *string
	// Category    Category `json:"category"`
}

type Category struct {
	*gorm.Model
	Name        *string   `json:"name" gorm:"unique;primarykey"`
	Description string    `json:"description"`
	Products    []Product `json:"products" gorm:"foreignKey:ProductName"`
}

type User struct {
	*gorm.Model
	ID        uuid.UUID `json:"id"`
	Email     *string   `json:"email" gorm:"unique"`
	Password  *string   `json:"password"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	ZipCode   string    `json:"zipCode"`
	Role      string    `json:"role"`
	// CartID    uint      `json:"cartId"`
	Cart Cart `json:"cart"`
}

type Cart struct {
	*gorm.Model
	ID         uuid.UUID `json:"id"`
	UserID     uuid.UUID `json:"userId"`
	Items      []Item    `json:"items"`
	TotalPrice float32   `json:"totalPrice"`
}

type Order struct {
	*gorm.Model
	User           *User   `json:"user" gorm:"unique"`
	Items          []Item  `json:"items"`
	TotalPrice     float32 `json:"totalPrice"`
	OrderStatus    string  `json:"orderStatus"`
	TrackingNumber string  `json:"trackingNumber"`
}

type Item struct {
	*gorm.Model
	ProductID  string  `json:"product"`
	Product    Product `json:"prodcut"`
	Quantity   uint    `json:"quantity"`
	TotalPrice float32 `json:"totalPrice"`
	CartID     string  `json:"cart"`
	OrderID    string  `json:"orderID, omitempty"`
}

type Price struct {
	Amount       uint   `json:"amount"`
	CurrencyCode string `json:"currencyCode"`
}

// type Role struct {
// 	Role string `json:"role"`
// }

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New()
	u.Cart.ID = uuid.New()
	// TODO erroru handle et
	// if !u.IsValid() {
	// 	err = errors.New("can't save invalid data")
	// }
	return
}

type Tokens struct {
	AccessToken  string
	RefreshToken string
}

// func (c *Cart) BeforeCreate(tx *gorm.DB) (err error) {
// 	c.ID = uuid.New()

// 	// if !u.IsValid() {
// 	// 	err = errors.New("can't save invalid data")
// 	// }
// 	return
// }
//
