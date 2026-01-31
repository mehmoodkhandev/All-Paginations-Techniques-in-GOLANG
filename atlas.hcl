data "external_schema" "gorm" {
  program = [
    "go", "run", "./loader"
  ]
}

env "gorm" {
  src = data.external_schema.gorm.url

  dev = "mysql://root:Khanzada%40123@localhost:3306/dev_tmp"

  url = "mysql://root:Khanzada%40123@localhost:3306/Pagination_Practice"

  migration {
    dir = "file://migrations"
  }
}
