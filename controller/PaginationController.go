package controller

import (
	"github.com/gofiber/fiber/v2"
)

type ProductPaginationController interface {
	CreateProducts(c *fiber.Ctx) error
	SortedProduct(c *fiber.Ctx) error
	FilterProduct(c *fiber.Ctx) error
	SearchProduct(c *fiber.Ctx) error
	PagePagination(c *fiber.Ctx) error
	OffsetPagination(c *fiber.Ctx) error
	CursorPagination(c *fiber.Ctx) error
}
