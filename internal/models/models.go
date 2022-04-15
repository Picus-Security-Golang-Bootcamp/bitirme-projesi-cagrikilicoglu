package models

import (
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Stock struct {
	SKU    string `json:"sku" gorm:"unique"`
	Number uint   `json:"number,omitempty"`
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
	Cart   Cart    `json:"cart"`
	Orders []Order `json:"orders"`
}

type Login struct {
	Email    *string `json:"email" gorm:"unique"`
	Password *string `json:"password"`
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
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
	ID         uuid.UUID      `json:"id"`
	UserID     uuid.UUID      `json:"userId"`
	Items      []Item         `json:"items"`
	TotalPrice float32        `json:"totalPrice"`
	Status     string         `json:"status"`
	// TrackingNumber string         `json:"trackingNumber"`
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
	OrderID    uuid.UUID      `json:"orderId,omitempty" gorm:"default:null"`
	// OrderedAt  time.Time      `json:"orderedAt" gorm:"default:null"`
	IsOrdered bool `json:"isOrdered" gorm:"default:false"`
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

var (
	statusPlaced   = "placed"
	statusCanceled = "canceled"
)

func (o *Order) BeforeCreate(tx *gorm.DB) (err error) {
	o.ID = uuid.New()
	o.Status = statusPlaced
	// TODO erroru handle et
	// if !u.IsValid() {
	// 	err = errors.New("can't save invalid data")
	// }
	return
}

func (o *Order) AfterDelete(tx *gorm.DB) (err error) {

	zap.L().Debug("order.afterdelete", zap.Reflect("id", o.ID))

	// TODO erroru handle et
	tx.Model(&o).Unscoped().Where("id = ?", o.ID).Select("status").Update("status", statusCanceled)
	// order := tx.Model(&o).First(&o)
	// zap.L().Debug("order.afterdeletew", zap.Reflect("order", order))
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
