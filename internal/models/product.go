package models

import (
	"github.com/Rhymond/go-money"
	"gorm.io/gorm"
)

type Stock struct {
	SKU         string `json:"SKU"`
	StockNumber uint   `json:"StockNumber"`
	StockStatus string `json:"StockStatus"`
}

type Product struct {
	*gorm.Model
	Name     *string     `json:"description"`
	Price    money.Money `json:"price"`
	Stock    Stock       `json:"stock" gorm:"embedded"`
	Category Category    `json:"category"`
}

type Category struct {
	*gorm.Model
	Name        *string   `json:"name" gorm:"unique"`
	Description string    `json:"description"`
	Products    []Product `json:"products"`
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
