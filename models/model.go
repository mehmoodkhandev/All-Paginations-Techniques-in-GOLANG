package models

import (
	"time"

	"github.com/MehmoodNadeemKhan1/Page_Pagination_API/helper"
	"gorm.io/gorm"
)

type Product struct {
	ID          string         `gorm:"type:char(36);primaryKey;default:(uuid())" json:"id"`
	Title       string         `gorm:"type:longtext;not null" json:"title"`
	Description string         `gorm:"type:longtext;not null" json:"description"`
	Image       string         `gorm:"type:longtext;not null" json:"image"`
	Price       int64          `gorm:"not null" json:"price"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

type PaginatedResult struct {
	Data       []Product `json:"data"`
	Page       int       `json:"page"`
	Limit      int       `json:"limit"`
	TotalItems int       `json:"total_items"`
	TotalPages int       `json:"total_pages"`
}

type OffsetResult struct {
	Data       []Product `json:"data"`
	Page       int       `json:"page"`
	Limit      int       `json:"limit"`
	Offset     int       `json:"offset"`
	TotalItems int       `json:"total_items"`
	TotalPages int       `json:"total_pages"`
}

type CursorResult struct {
	Success    bool                  `json:"success"`
	Data       []Product             `json:"data"`
	Pagination helper.PaginationInfo `json:"pagination"`
	//HasMore bool `json:"has_more"`
}

func (p Product) GetID() string {
	return p.ID
}

func (p Product) GetCreatedAt() time.Time {
	return p.CreatedAt
}

// func(product *Product)BeforeCreate(tx *gorm.DB)(err error){
// 	product.ID = uuid.NewString()
// 	return
// }
