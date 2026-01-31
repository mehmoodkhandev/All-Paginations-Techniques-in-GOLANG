package main

import (
	"fmt"
	"log"

	"ariga.io/atlas-provider-gorm/gormschema"
	"github.com/MehmoodNadeemKhan1/Page_Pagination_API/models"
)

func main() {
	ddl, err := gormschema.New("mysql").Load(
		&models.Product{},
	)
	if err != nil {
		log.Fatal("Loader Error:", err)
	}
	fmt.Print(ddl)
}
