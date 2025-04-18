package routes

import (
	"database/sql"
	"simple_mysql_elasticsearch/config"
	"simple_mysql_elasticsearch/internal/handler"
	rp "simple_mysql_elasticsearch/internal/repository"
	uc "simple_mysql_elasticsearch/internal/usecase"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Echo, cfg *config.Config, db *sql.DB, esClient *elasticsearch.Client) {
	mysqlRepo := &rp.ProductMySQL{DB: db}
	elasticRepo := &rp.ProductElastic{Client: esClient}
	usecase := &uc.ProductElastic{RepoMysql: mysqlRepo, RepoElastic: elasticRepo}
	handler := &handler.ProductHandler{Config: cfg, Usecase: usecase}

	e.POST("/products", handler.Create)
	e.PUT("/products", handler.Update)
	e.GET("/products", handler.GetAll)
	e.GET("/products/:id", handler.GetByID)
	e.DELETE("/products/:id", handler.Delete)

	e.GET("/products/search", handler.SearchProductHandler)

	// Serve static images
	e.Static("/uploads", "uploads")
}
