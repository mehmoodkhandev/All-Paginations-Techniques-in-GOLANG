package service

import (
	"github.com/MehmoodNadeemKhan1/Page_Pagination_API/models"
)

type ProductPaginationService interface {
	CreateProducts() error
	SortedProduct(orderBy string, direction string) ([]models.Product, error)
	FilterProduct(category string, minPrice, maxPrice float64) ([]models.Product, error)
	SearchProduct(query string) ([]models.Product, error)
	PagePagination(page, limit int) (models.PaginatedResult, error)
	OffsetPagination(offset, limit int) (models.OffsetResult, error)
	CursorPagination(cursor string, limit int64, sort string) (models.CursorResult, error)
}
