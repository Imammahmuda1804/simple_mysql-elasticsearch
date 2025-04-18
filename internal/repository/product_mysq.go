package repository

import (
	"database/sql"
	"fmt"
	"simple_mysql_elasticsearch/internal/domain"
)

type ProductMySQL struct {
	DB *sql.DB
}

func (p *ProductMySQL) Create(product domain.Product) error {
	fmt.Println("url", product.Images)
	sql := "INSERT INTO products (sku,name, description, price, category, image_path) VALUES (?,?, ?, ?, ?, ?)"
	_, err := p.DB.Exec(sql, product.SKU,
		product.Name, product.Description, product.Price, product.Category, product.Images)

	fmt.Println(sql, product.SKU, product.Name, product.Description, product.Price, product.Category, product.Images)

	fmt.Println("erororor", err)
	return err
}

func (p *ProductMySQL) Update(product domain.Product) error {
	_, err := p.DB.Exec("UPDATE products SET sku=?,name=?, description=?, price=?, category=?, image_path=? WHERE id=?",
		product.SKU, product.Name, product.Description, product.Price, product.Category, product.Images, product.ID)
	return err
}

func (p *ProductMySQL) GetBySKU(id string) (domain.Product, error) {
	var product domain.Product
	err := p.DB.QueryRow("SELECT id, sku,name, description, price, category, image_path FROM products WHERE sku = ?", id).
		Scan(&product.ID, &product.SKU, &product.Name, &product.Description, &product.Price, &product.Category, &product.Images)
	if err != nil {
		return domain.Product{}, err
	}
	return product, nil
}

func (p *ProductMySQL) GetByID(id int) (domain.Product, error) {
	var product domain.Product
	err := p.DB.QueryRow("SELECT id, sku, name, description, price, category, image_path FROM products WHERE id = ?", id).
		Scan(&product.ID, &product.SKU, &product.Name, &product.Description, &product.Price, &product.Category, &product.Images)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.Product{}, fmt.Errorf("product with ID %d not found", id)
		}
		return domain.Product{}, err
	}
	return product, nil
}

func (p *ProductMySQL) GetAll() ([]domain.Product, error) {
	rows, err := p.DB.Query("SELECT id, name, description, price, category, image_path FROM products")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []domain.Product
	for rows.Next() {
		var product domain.Product
		if err := rows.Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.Category, &product.Images); err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	return products, nil
}

func (p *ProductMySQL) Delete(id int) error {
	_, err := p.DB.Exec("DELETE FROM products WHERE id = ?", id)
	return err
}
