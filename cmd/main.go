package main

import (
	"simple_mysql_elasticsearch/config"
	"simple_mysql_elasticsearch/routes"

	"github.com/labstack/echo/v4"
)

func main() {
	cfg := config.LoadConfig()

	db, esClient := config.InitDBAndElastic(cfg)

	e := echo.New()
	routes.RegisterRoutes(e, cfg, db, esClient)
	e.Logger.Fatal(e.Start(":8085"))
}
