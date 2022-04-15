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
	Cart      Cart           `json:"cart"`
	Orders    []Order        `json:"orders"`
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
	IsOrdered  bool           `json:"isOrdered" gorm:"default:false"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New()
	u.Cart.ID = uuid.New()
	return
}

func (c *Category) BeforeCreate(tx *gorm.DB) (err error) {
	c.ID = uuid.New()
	return
}

func (p *Product) BeforeCreate(tx *gorm.DB) (err error) {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return
}
func (i *Item) BeforeCreate(tx *gorm.DB) (err error) {
	i.ID = uuid.New()
	return
}

var (
	statusPlaced   = "placed"
	statusCanceled = "canceled"
)

func (o *Order) BeforeCreate(tx *gorm.DB) (err error) {
	o.ID = uuid.New()
	o.Status = statusPlaced
	return
}

func (o *Order) AfterDelete(tx *gorm.DB) (err error) {

	zap.L().Debug("order.afterdelete", zap.Reflect("id", o.ID))

	result := tx.Model(&o).Unscoped().Where("id = ?", o.ID).Select("status").Update("status", statusCanceled)
	if result.Error != nil {
		return result.Error
	}
	return
}

type Tokens struct {
	AccessToken  string
	RefreshToken string
}
