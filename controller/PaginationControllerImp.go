package controller

import (
	"strconv"

	"github.com/MehmoodNadeemKhan1/Page_Pagination_API/service"
	"github.com/gofiber/fiber/v2"
)

type ProductController struct {
	service service.ProductPaginationService
}

// CreateProducts implements ProductPaginationController.
func (p *ProductController) CreateProducts(c *fiber.Ctx) error {
	if err := p.service.CreateProducts(); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(201).JSON(fiber.Map{"message": "Products created successfully"})

}

// CursorPagination implements ProductPaginationController.
func (p *ProductController) CursorPagination(c *fiber.Ctx) error {
	cursorStr := c.Query("cursor", "")
	limitStr := c.Query("limit", "10")
	sortOrder := c.Query("sort_order", "desc")

	limitNum, err := strconv.ParseInt(limitStr, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid limit number",
		})
	}
	products, err := p.service.CursorPagination(cursorStr, limitNum, sortOrder)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(products)

}

// FilterProduct implements ProductPaginationController.
func (p *ProductController) FilterProduct(c *fiber.Ctx) error {
	category := c.Query("category")
	minPriceStr := c.Query("minPrice")
	maxPriceStr := c.Query("maxPrice")
	var minPrice, maxPrice float64
	var err error
	if minPriceStr != "" {
		minPrice, err = strconv.ParseFloat(minPriceStr, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid MinPrice value",
			})
		}
	}
	if maxPriceStr != "" {
		maxPrice, err = strconv.ParseFloat(maxPriceStr, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid MaxPrice value",
			})
		}
	}
	// Call service layer
	products, err := p.service.FilterProduct(category, minPrice, maxPrice)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(products)
}

// OffsetPagination implements ProductPaginationController.
func (p *ProductController) OffsetPagination(c *fiber.Ctx) error {
	offsetStr := c.Query("offset", "0")
	limitStr := c.Query("limit", "9")

	offsetNum, err := strconv.Atoi(offsetStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid Offset number",
		})
	}

	limitNum, err := strconv.Atoi(limitStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid limit number",
		})
	}

	products, err := p.service.OffsetPagination(offsetNum, limitNum)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(products)

}

// PagePagination implements ProductPaginationController.

func (p *ProductController) PagePagination(c *fiber.Ctx) error {
	pageStr := c.Query("page", "1")   // Default page = 1
	limitStr := c.Query("limit", "9") // Default limit = 9

	pageNum, err := strconv.Atoi(pageStr)
	if err != nil || pageNum <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid page number",
		})
	}

	limitNum, err := strconv.Atoi(limitStr)
	if err != nil || limitNum <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid limit number",
		})
	}

	// Call the service
	products, err := p.service.PagePagination(pageNum, limitNum)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(products)
}

func (p *ProductController) SearchProduct(c *fiber.Ctx) error {
	query := c.Query("query")

	// Validation: check empty search
	if query == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Query parameter is required",
		})
	}

	// Call service
	products, err := p.service.SearchProduct(query)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Handle empty result
	if len(products) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "No products found matching your search",
		})
	}

	return c.Status(fiber.StatusOK).JSON(products)
}

// SortedProduct implements ProductPaginationController.
func (p *ProductController) SortedProduct(c *fiber.Ctx) error {
	OrderBy := c.Query("orderBy")
	direction := c.Query("direction")

	products, err := p.service.SortedProduct(OrderBy, direction)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"Error": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(&products)

}

func NewProductController(s service.ProductPaginationService) ProductPaginationController {
	return &ProductController{
		service: s,
	}
}
