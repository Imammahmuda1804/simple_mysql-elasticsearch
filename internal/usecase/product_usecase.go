package usecase

import (
	"fmt"
	"simple_mysql_elasticsearch/internal/domain"
	rp "simple_mysql_elasticsearch/internal/repository"
)

type ProductElastic struct {
	RepoMysql   *rp.ProductMySQL
	RepoElastic *rp.ProductElastic
}

func (uc *ProductElastic) Create(product domain.Product) (id int, error error) {
	// Simpan ke MySQL
	err := uc.RepoMysql.Create(product)
	if err != nil {
		return 0, err
	}

	// Ambil kembali ID dari MySQL (jika diperlukan)
	savedProduct, err := uc.RepoMysql.GetBySKU(product.SKU)
	if err != nil {
		return savedProduct.ID, err
	}

	// Simpan ke Elasticsearch
	err = uc.RepoElastic.Create(savedProduct)
	if err != nil {
		fmt.Println("Error saving to Elasticsearch:", err)
		return savedProduct.ID, err
	}

	return savedProduct.ID, nil
}

func (uc *ProductElastic) Update(product domain.Product) error {
	// Perbarui di MySQL
	err := uc.RepoMysql.Update(product)
	if err != nil {
		return err
	}

	// Perbarui di Elasticsearch
	err = uc.RepoElastic.Update(product)
	if err != nil {
		return err
	}

	return nil
}

func (uc *ProductElastic) GetByID(id int) (domain.Product, error) {
	// ambil dari MySQL
	product, err := uc.RepoMysql.GetByID(id)
	//if err != nil {
	return product, err
	//}

	// Jika tidak ditemukan di MySQL, coba dari Elasticsearch
	// product, err = uc.RepoElastic.GetByID(id)
	// if err != nil {
	// 	return domain.Product{}, fmt.Errorf("product not found")
	// }

	//return product, nil
}

func (uc *ProductElastic) GetAll() ([]domain.Product, error) {
	// Ambil dari MySQL
	products, err := uc.RepoMysql.GetAll()
	if err != nil {
		return nil, err
	}

	return products, nil
}

func (uc *ProductElastic) SearchProductByKeyword(keyword string) ([]domain.Product, error) {
	return uc.RepoElastic.SearchByKeyword(keyword)
}

func (uc *ProductElastic) Delete(id int) error {
	// Perbarui di MySQL
	err := uc.RepoMysql.Delete(id)
	if err != nil {
		return err
	}

	// Perbarui di Elasticsearch
	err = uc.RepoElastic.Delete(id)
	if err != nil {
		return err
	}

	return nil
}
