package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Username  string         `gorm:"uniqueIndex;not null"`
	Password  string         `gorm:"not null"`
	Name      string
	Phone     string         `gorm:"uniqueIndex"`
	Role      string         `gorm:"default:customer"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type Category struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name      string         `gorm:"not null"`
	SortOrder int            `gorm:"default:0"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type Product struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	CategoryID  uuid.UUID      `gorm:"type:uuid;not null"`
	Name        string         `gorm:"not null"`
	Description string
	Price       float64        `gorm:"not null"`
	Image       string
	Status      string         `gorm:"default:active"`
	SortOrder   int            `gorm:"default:0"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
	Category    Category       `gorm:"foreignKey:CategoryID"`
}

type Order struct {
	ID            uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID        uuid.UUID      `gorm:"type:uuid;not null"`
	OrderNo       string         `gorm:"uniqueIndex;not null"`
	TotalAmount   float64        `gorm:"not null"`
	Status        string         `gorm:"default:pending"`
	Remark        string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt `gorm:"index"`
	User          User           `gorm:"foreignKey:UserID"`
	OrderItems    []OrderItem    `gorm:"foreignKey:OrderID"`
}

type OrderItem struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	OrderID   uuid.UUID `gorm:"type:uuid;not null"`
	ProductID uuid.UUID `gorm:"type:uuid;not null"`
	Quantity  int       `gorm:"not null"`
	Price     float64   `gorm:"not null"`
	Product   Product   `gorm:"foreignKey:ProductID"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}

func (c *Category) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}

func (p *Product) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}

func (o *Order) BeforeCreate(tx *gorm.DB) error {
	if o.ID == uuid.Nil {
		o.ID = uuid.New()
	}
	return nil
}

func (oi *OrderItem) BeforeCreate(tx *gorm.DB) error {
	if oi.ID == uuid.Nil {
		oi.ID = uuid.New()
	}
	return nil
}
