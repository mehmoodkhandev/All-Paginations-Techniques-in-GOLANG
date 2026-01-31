package service

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/MehmoodNadeemKhan1/Page_Pagination_API/database"
	"github.com/MehmoodNadeemKhan1/Page_Pagination_API/helper"
	"github.com/MehmoodNadeemKhan1/Page_Pagination_API/models"
	"github.com/bxcodec/faker/v3"

	"gorm.io/gorm"
)

type ProductServices struct {
	db *gorm.DB
}

func (p *ProductServices) CreateProducts() error {
	for i := 0; i < 100; i++ {
		product := models.Product{
			ID:          faker.UUIDDigit(),
			Title:       faker.Word(),
			Description: faker.Paragraph(),
			Image:       fmt.Sprintf("http://lorempixel.com/200/200?random=%s", faker.UUIDDigit()),
			Price:       int64(rand.Intn(90) + 10),
		}

		// save product to DB and check for errors
		if err := p.db.Create(&product).Error; err != nil {
			return err // stop immediately if an insert fails
		}
	}

	return nil
}

// CursorPagination implements ProductPaginationService.
// Jab tumhe infinite scrolling implement karni ho
// Jab data real-time update hota ho (like social feeds, blockchain explorers)
// Jab performance critical ho aur tum OFFSET skip karna avoid karna chaho
func (p *ProductServices) CursorPagination(cursor string, limit int64, sort string) (models.CursorResult, error) {
	var products []models.Product
	if limit < 1 || limit > 100 {
		limit = 10
	}
	isFirstPage := cursor == ""
	pointNext := false
	query := p.db

	query, pointNext, err := database.GetPaginatedQuery(query, pointNext, cursor, sort)
	if err != nil {
		return models.CursorResult{}, fmt.Errorf("unable to get the paginated query: %w", err)
	}

	err = query.Limit(int(limit) + 1).Find(&products).Error
	if err != nil {
		return models.CursorResult{}, fmt.Errorf("unable to get the Products data: %w", err)
	}

	hasPagination := len(products) > int(limit)

	if hasPagination {
		products = products[:limit]

	}

	if !isFirstPage && !pointNext {
		products = helper.Reverse(products)
	}

	pageInfo := database.CalculatePagination(isFirstPage, hasPagination, int(limit), products, pointNext)
	response := models.CursorResult{
		Success:    true,
		Data:       products,
		Pagination: pageInfo,
	}

	return response, nil

}

func (p *ProductServices) FilterProduct(category string, minPrice, maxPrice float64) ([]models.Product, error) {
	var products []models.Product

	query := p.db.Model(&models.Product{})

	// Filter by category (if provided)
	if category != "" {
		query = query.Where("title = ?", category)
	}

	// Filter by minimum price
	if minPrice > 0 {
		query = query.Where("price >= ?", minPrice)
	}

	// Filter by maximum price
	if maxPrice > 0 {
		query = query.Where("price <= ?", maxPrice)
	}

	// Execute query
	if err := query.Find(&products).Error; err != nil {
		return nil, err
	}

	return products, nil
}

// OffsetPagination implements ProductPaginationService.
//Jab tumhe data scrolling type fetch karna ho (like “Load More” button)

//Jab frontend developer offset manually handle kar raha ho

// Jab tumhe fixed or automated data sync karni ho kisi service ke beech
func (p *ProductServices) OffsetPagination(offset int, limit int) (models.OffsetResult, error) {
	// ✅ Default limit if not provided
	if limit <= 0 {
		limit = 9
	}

	if offset < 0 {
		offset = 0
	}

	var products []models.Product
	var total int64
	var result models.OffsetResult

	// ✅ Count total records
	if err := p.db.Model(&models.Product{}).Count(&total).Error; err != nil {
		return result, err
	}

	// ✅ Fetch data with offset and limit
	if err := p.db.Offset(offset).Limit(limit).Find(&products).Error; err != nil {
		return result, err
	}

	// ✅ Calculate total pages
	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	// ✅ Calculate current page number from offset
	page := int(math.Floor(float64(offset)/float64(limit))) + 1

	result = models.OffsetResult{
		Data:       products,
		Page:       page,
		Limit:      limit,
		Offset:     offset,
		TotalItems: int(total),
		TotalPages: totalPages,
	}

	return result, nil

}

// PagePagination implements ProductPaginationService.
//Jab frontend users ke liye UI buttons hain ("1 2 3 4 5")

//Jab total count dikhani hoti hai ("Showing 20–30 of 100 products")

// Jab data set zyada bada nahi hota (like 10,000 records tak)
func (p *ProductServices) PagePagination(page int, limit int) (models.PaginatedResult, error) {

	// Default values
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 9
	}

	var products []models.Product
	var total int64
	var result models.PaginatedResult
	// Count total records
	if err := p.db.Model(&models.Product{}).Count(&total).Error; err != nil {
		return result, err
	}

	// Calculate offset
	offset := (page - 1) * limit

	// Fetch paginated records
	if err := p.db.Limit(limit).Offset(offset).Find(&products).Error; err != nil {
		return result, err
	}

	// Calculate total pages
	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	// Prepare paginated response
	result = models.PaginatedResult{
		Data:       products,
		Page:       page,
		Limit:      limit,
		TotalItems: int(total),
		TotalPages: totalPages,
	}

	return result, nil

}

// SearchProduct implements ProductPaginationService.
func (p *ProductServices) SearchProduct(query string) ([]models.Product, error) {
	var products []models.Product

	if query == "" {
		// if query is empty, return all products it will store all products in products array
		if err := p.db.Find(&products).Error; err != nil {
			return nil, err
		}
		return products, nil
	}
	// Use GORM to safely perform LIKE search (prevents SQL injection)
	// now finding in that products that saved previously in array find in that specific searched product
	err := p.db.Where("title LIKE ? OR description LIKE ?", "%"+query+"%", "%"+query+"%").Find(&products).Error
	if err != nil {
		return nil, err
	}

	return products, nil
}
func (p *ProductServices) SortedProduct(orderBy string, direction string) ([]models.Product, error) {
	var products []models.Product

	validColumns := map[string]bool{
		"title":       true,
		"description": true,
		"price":       true,
		"created_at":  true,
	}
	if !validColumns[orderBy] {
		orderBy = "created_at" // default column
	}

	if direction != "asc" && direction != "desc" {
		direction = "asc" // default direction
	}

	order := fmt.Sprintf("%s %s", orderBy, direction)

	if err := p.db.Order(order).Find(&products).Error; err != nil {
		return nil, err
	}

	return products, nil
}

func NewProductService(db *gorm.DB) ProductPaginationService {
	return &ProductServices{
		db: db,
	}
}
