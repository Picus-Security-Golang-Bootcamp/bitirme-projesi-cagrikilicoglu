package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Stock struct {
	SKU    string `json:"sku" gorm:"unique"`
	Number uint   `json:"number"`
	Status string `json:"status"`
}

type Product struct {
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
	ID           uuid.UUID      `json:"id"`
	Name         *string        `json:"description"`
	Price        float32        `json:"price"`
	Stock        Stock          `json:"stock" gorm:"embedded"`
	CategoryName *string        `json:"categoryName"`
	// Category     Category  `json:"category"`
}

type Category struct {
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
	ID          uuid.UUID      `json:"id"`
	Name        *string        `json:"name" gorm:"unique;primaryKey"`
	Description string         `json:"description"`
	Products    []Product      `json:"products" gorm:"foreignKey:CategoryName"`
}

type User struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	ID        uuid.UUID      `json:"id"`
	Email     *string        `json:"email" gorm:"unique"`
	Password  *string        `json:"password"`
	FirstName string         `json:"firstName"`
	LastName  string         `json:"lastName"`
	ZipCode   string         `json:"zipCode"`
	Role      string         `json:"role"`
	// CartID    uint      `json:"cartId"`
	Cart Cart `json:"cart"`
}

type Cart struct {
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
	ID         uuid.UUID      `json:"id"`
	UserID     uuid.UUID      `json:"userId"`
	Items      []Item         `json:"items"`
	TotalPrice float32        `json:"totalPrice"`
}

type Order struct {
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      gorm.DeletedAt `gorm:"index"`
	User           *User          `json:"user" gorm:"unique"`
	Items          []Item         `json:"items"`
	TotalPrice     float32        `json:"totalPrice"`
	OrderStatus    string         `json:"orderStatus"`
	TrackingNumber string         `json:"trackingNumber"`
}

type Item struct {
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
	ID         uuid.UUID      `json:"id"`
	ProductID  uuid.UUID      `json:"productID"`
	Product    Product        `json:"product" gorm:"constraint:OnUpdate:CASCADE;"`
	Quantity   uint           `json:"quantity"`
	TotalPrice float32        `json:"totalPrice"`
	CartID     uuid.UUID      `json:"cartId"`
	OrderID    string         `json:"orderID,omitempty"`
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

// TODO aşağıdaki iki fonksiyonu başka bir yere alabilir miyiz?
// func (c *Cart) AfterUpdate(tx *gorm.DB) (err error) {
// 	zap.L().Debug("cart.afterupdate")
// 	c.CalculatePrice()
// 	return
// }
// func (c *Cart) AfterDelete(tx *gorm.DB) (err error) {
// 	c.CalculatePrice()
// 	return
// }

// func (c *Cart) CalculatePrice() {
// 	c.TotalPrice = 0
// 	for i := range c.Items {
// 		c.TotalPrice += c.Items[i].TotalPrice
// 	}
// }

func (c *Category) BeforeCreate(tx *gorm.DB) (err error) {
	c.ID = uuid.New()
	// TODO erroru handle et
	// if !u.IsValid() {
	// 	err = errors.New("can't save invalid data")
	// }
	return
}

func (p *Product) BeforeCreate(tx *gorm.DB) (err error) {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	// TODO erroru handle et
	// if !u.IsValid() {
	// 	err = errors.New("can't save invalid data")
	// }
	return
}
func (i *Item) BeforeCreate(tx *gorm.DB) (err error) {
	i.ID = uuid.New()
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
