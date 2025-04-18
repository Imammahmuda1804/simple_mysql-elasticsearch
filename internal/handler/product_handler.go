package handler

import (
	"fmt"
	"net/http"
	"simple_mysql_elasticsearch/config"
	"simple_mysql_elasticsearch/helper"
	"simple_mysql_elasticsearch/internal/domain"
	uc "simple_mysql_elasticsearch/internal/usecase"
	"strconv"

	"github.com/labstack/echo/v4"
)

//const ImageDir = "uploads"
// Ensure image directory exists
// func init() {
// 	if _, err := os.Stat(ImageDir); os.IsNotExist(err) {
// 		os.Mkdir(ImageDir, os.ModePerm)
// 	}
// }

type ProductHandler struct {
	Config  *config.Config
	Usecase *uc.ProductElastic
}

func (h *ProductHandler) Create(c echo.Context) error {
	var product domain.Product
	if err := c.Bind(&product); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request body"})
	}

	// Process base64 image
	if product.Images != "" {
		imagePath, err := helper.SaveBase64Image(h.Config.ImageDir, product.Images)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to save image" + err.Error()})
		}
		product.Images = imagePath
	}

	var id int
	id, err := h.Usecase.Create(product)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	product.ID = id
	product.Images = fmt.Sprintf("%s/%s", h.Config.BaseUrl, product.Images)

	return c.JSON(http.StatusCreated, product)
}

func (h *ProductHandler) Update(c echo.Context) error {
	var product domain.Product
	if err := c.Bind(&product); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request body"})
	}

	// Process base64 image if updated
	if product.Images != "" {
		imagePath, err := helper.SaveBase64Image(h.Config.ImageDir, product.Images)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to save image"})
		}
		product.Images = fmt.Sprintf("%s/%s", h.Config.BaseUrl, imagePath)
	}

	if err := h.Usecase.Update(product); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}
	return c.JSON(http.StatusOK, product)
}

func (h *ProductHandler) GetByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid product ID"})
	}

	product, err := h.Usecase.GetByID(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "Product not found"})
	}

	// Convert local image path to URL
	product.Images = fmt.Sprintf("%s/%s", h.Config.BaseUrl, product.Images)

	return c.JSON(http.StatusOK, product)
}

func (h *ProductHandler) GetAll(c echo.Context) error {
	products, err := h.Usecase.GetAll()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	// Convert image paths to URLs
	for i := range products {
		products[i].Images = fmt.Sprintf("%s/%s", h.Config.BaseUrl, products[i].Images)
	}

	return c.JSON(http.StatusOK, products)
}

func (h *ProductHandler) SearchProductHandler(c echo.Context) error {
	keyword := c.QueryParam("keyword")
	if keyword == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Keyword is required"})
	}

	products, err := h.Usecase.SearchProductByKeyword(keyword)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	// Update Images for each product
	for i := range products {
		products[i].Images = fmt.Sprintf("%s/%s", h.Config.BaseUrl, products[i].Images)
	}

	return c.JSON(http.StatusOK, products)
}

func (h *ProductHandler) Delete(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid product ID"})
	}

	err = h.Usecase.Delete(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "Product not found"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "success delete product"})
}
