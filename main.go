package main

import (
	"log"

	"github.com/MehmoodNadeemKhan1/Page_Pagination_API/config"
	"github.com/MehmoodNadeemKhan1/Page_Pagination_API/controller"
	"github.com/MehmoodNadeemKhan1/Page_Pagination_API/routes"
	"github.com/MehmoodNadeemKhan1/Page_Pagination_API/service"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	Loadconfig, err := config.LoadConfig(".")
	if err != nil {
		return
	}
	db, err := config.ConnectingMySQLDb(&Loadconfig)
	if err != nil {
		log.Fatal("SQL ERROR: ", err)
	}

	service := service.NewProductService(db)
	controller := controller.NewProductController(service)
	routes := routes.NewRouter(app, controller)

	routes.Listen(":8000")
}
