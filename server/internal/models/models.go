package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Username  string         `json:"username" gorm:"uniqueIndex;not null"`
	Email     string         `json:"email" gorm:"uniqueIndex;not null"`
	Password  string         `json:"-" gorm:"not null"`
	Phone     string         `json:"phone"`
	Address   string         `json:"address"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type Merchant struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"not null"`
	Description string         `json:"description"`
	Address     string         `json:"address"`
	Phone       string         `json:"phone"`
	Status      string         `json:"status" gorm:"default:'active'"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

type Dish struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	MerchantID  uint           `json:"merchant_id" gorm:"not null"`
	Name        string         `json:"name" gorm:"not null"`
	Description string         `json:"description"`
	Price       float64        `json:"price" gorm:"not null"`
	Image       string         `json:"image"`
	Category    string         `json:"category"`
	Status      string         `json:"status" gorm:"default:'available'"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
	Merchant    Merchant       `json:"merchant" gorm:"foreignKey:MerchantID"`
}

type CartItem struct {
	ID         uint           `json:"id" gorm:"primaryKey"`
	UserID     uint           `json:"user_id" gorm:"not null"`
	MerchantID uint           `json:"merchant_id" gorm:"not null"`
	DishID     uint           `json:"dish_id" gorm:"not null"`
	Quantity   int            `json:"quantity" gorm:"not null"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`
	User       User           `json:"user" gorm:"foreignKey:UserID"`
	Merchant   Merchant       `json:"merchant" gorm:"foreignKey:MerchantID"`
	Dish       Dish           `json:"dish" gorm:"foreignKey:DishID"`
}

type Order struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	UserID      uint           `json:"user_id" gorm:"not null"`
	MerchantID  uint           `json:"merchant_id" gorm:"not null"`
	OrderNumber string         `json:"order_number" gorm:"uniqueIndex;not null"`
	TotalAmount float64        `json:"total_amount" gorm:"not null"`
	Status      string         `json:"status" gorm:"default:'pending'"`
	Address     string         `json:"address" gorm:"not null"`
	Phone       string         `json:"phone" gorm:"not null"`
	Remark      string         `json:"remark"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
	User        User           `json:"user" gorm:"foreignKey:UserID"`
	Merchant    Merchant       `json:"merchant" gorm:"foreignKey:MerchantID"`
	OrderItems  []OrderItem    `json:"order_items" gorm:"foreignKey:OrderID"`
}

type OrderItem struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	OrderID   uint           `json:"order_id" gorm:"not null"`
	DishID    uint           `json:"dish_id" gorm:"not null"`
	Quantity  int            `json:"quantity" gorm:"not null"`
	Price     float64        `json:"price" gorm:"not null"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	Order     Order          `json:"order" gorm:"foreignKey:OrderID"`
	Dish      Dish           `json:"dish" gorm:"foreignKey:DishID"`
}
