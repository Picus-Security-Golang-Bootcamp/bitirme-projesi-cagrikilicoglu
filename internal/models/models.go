package models

import (
	"github.com/Rhymond/go-money"
	"gorm.io/gorm"
)

type Stock struct {
	SKU    string `json:"sku" gorm:"unique"`
	Number uint   `json:"number"`
	Status string `json:"status"`
}

type Product struct {
	gorm.Model
	Name         *string `json:"description" gorm`
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
	Email     *string `json:"email" gorm:"unique;primaryKey"`
	Password  *string `json:"password"`
	FirstName string  `json:"firstName"`
	LastName  string  `json:"lastName"`
	ZipCode   string  `json:"zipCode"`
	Role      string  `json:"role"`
	Cart      Cart    `json:"cart"`
}

type Cart struct {
	*gorm.Model
	User  *User       `json:"user" gorm:"unique"`
	Items []Item      `json:"items"`
	Total money.Money `json:"total"`
}

type Order struct {
	*gorm.Model
	User           *User       `json:"user" gorm:"unique"`
	Items          []Item      `json:"items"`
	Total          money.Money `json:"total"`
	OrderStatus    string      `json:"orderStatus"`
	TrackingNumber string      `json:"trackingNumber"`
}

type Item struct {
	*gorm.Model
	Product    Product     `json:"prodcut"`
	Quantity   uint        `json:"quantity"`
	TotalPrice money.Money `json:"totalPrice"`
}

type Price struct {
	Amount       uint   `json:"amount"`
	CurrencyCode string `json:"currencyCode"`
}
