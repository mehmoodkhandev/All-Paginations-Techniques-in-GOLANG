package routes

import (
	"github.com/MehmoodNadeemKhan1/Page_Pagination_API/controller"
	"github.com/gofiber/fiber/v2"
)

func NewRouter(router *fiber.App, productController controller.ProductPaginationController) *fiber.App {
	router.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.SendString("Hello world")
	})

	router.Post("/product", productController.CreateProducts)
	router.Get("/productSearch", productController.SearchProduct)
	router.Get("/productFilter", productController.FilterProduct)
	router.Get("/productSorted", productController.SortedProduct)
	router.Get("/pagePaginatedProduct", productController.PagePagination)
	router.Get("/offsetPaginatedProduct", productController.OffsetPagination)
	router.Get("/cursorPaginatedProduct", productController.CursorPagination)

	return router
}
