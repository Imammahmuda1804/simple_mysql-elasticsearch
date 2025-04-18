package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"simple_mysql_elasticsearch/internal/domain"
	"strconv"

	"github.com/elastic/go-elasticsearch/v8"
)

type ProductElastic struct {
	Client *elasticsearch.Client
}

func (p *ProductElastic) Create(product domain.Product) error {
	data, err := json.Marshal(product)
	if err != nil {
		return err
	}

	fmt.Println("SKU   ", product.SKU)
	fmt.Println("ID   ", product.ID)

	res, err := p.Client.Index(
		"products",
		bytes.NewReader(data),
		p.Client.Index.WithDocumentID(strconv.Itoa(product.ID)),
		p.Client.Index.WithRefresh("true"),
	)

	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.IsError() {
		body, _ := io.ReadAll(res.Body)
		return fmt.Errorf("error indexing document: %s", string(body))
	}

	return nil
}

func (p *ProductElastic) Update(product domain.Product) error {
	// Re-index document
	return p.Create(product)
}

func (p *ProductElastic) Delete(id int) error {
	res, err := p.Client.Delete(
		"products",
		strconv.Itoa(id),
		p.Client.Delete.WithRefresh("true"),
	)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		body, _ := io.ReadAll(res.Body)
		return fmt.Errorf("error deleting document: %s", string(body))
	}

	return nil
}

func (p *ProductElastic) GetAll() ([]domain.Product, error) {
	var buf bytes.Buffer
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match_all": map[string]interface{}{},
		},
	}

	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return nil, err
	}

	res, err := p.Client.Search(
		p.Client.Search.WithContext(context.Background()),
		p.Client.Search.WithIndex("products"),
		p.Client.Search.WithBody(&buf),
		p.Client.Search.WithTrackTotalHits(true),
		p.Client.Search.WithPretty(),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		body, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("error searching documents: %s", string(body))
	}

	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, err
	}

	hits := r["hits"].(map[string]interface{})["hits"].([]interface{})
	var products []domain.Product

	for _, hit := range hits {
		source := hit.(map[string]interface{})["_source"]
		sourceBytes, _ := json.Marshal(source)

		var product domain.Product
		if err := json.Unmarshal(sourceBytes, &product); err != nil {
			return nil, err
		}

		products = append(products, product)
	}

	return products, nil
}

func (p *ProductElastic) SearchByKeyword(keyword string) ([]domain.Product, error) {
	fmt.Println("keyword", keyword)
	var buf bytes.Buffer

	// Elasticsearch query: mencari pada name, description, dan category
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"should": []interface{}{
					// Pencarian di name dan description dengan fuzziness AUTO
					map[string]interface{}{
						"multi_match": map[string]interface{}{
							"query":     keyword,
							"fields":    []string{"name", "description"},
							"fuzziness": "AUTO",
						},
					},
					// Pencarian partial match di category
					map[string]interface{}{
						"match_phrase_prefix": map[string]interface{}{
							"category": keyword,
						},
					},
				},
			},
		},
	}

	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return nil, err
	}

	// Perform the search request
	res, err := p.Client.Search(
		p.Client.Search.WithIndex("products"),
		p.Client.Search.WithBody(&buf),
		p.Client.Search.WithTrackTotalHits(true),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// Parse response
	if res.IsError() {
		return nil, errors.New("error from Elasticsearch: " + res.String())
	}

	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, err
	}

	// Extract hits
	hits := r["hits"].(map[string]interface{})["hits"].([]interface{})
	var products []domain.Product
	for _, hit := range hits {
		source := hit.(map[string]interface{})["_source"]
		sourceBytes, _ := json.Marshal(source)

		var product domain.Product
		if err := json.Unmarshal(sourceBytes, &product); err != nil {
			continue
		}
		products = append(products, product)
	}

	return products, nil
}
