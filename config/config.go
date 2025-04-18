package config

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	_ "github.com/go-sql-driver/mysql"
	"github.com/kelseyhightower/envconfig"
)

type Config struct { //uploads
	BaseUrl       string `envconfig:"BASE_URL" required:"true" default:"http://localhost:8085"`
	ImageDir      string `envconfig:"IMAGE_DIR" required:"true" default:"uploads"`
	MySQLUser     string `envconfig:"MYSQL_USER" required:"true" default:"root"`
	MySQLPassword string `envconfig:"MYSQL_PASSWORD" required:"true" default:"root"`
	MySQLHost     string `envconfig:"MYSQL_HOST" default:"localhost"`
	MySQLPort     int    `envconfig:"MYSQL_PORT" default:"3306"`
	MySQLDBName   string `envconfig:"MYSQL_DBNAME" required:"true" default:"sales"`

	ElasticAddresses string `envconfig:"ELASTIC_ADDRESSES" default:"http://localhost:9200"`
	ElasticUsername  string `envconfig:"ELASTIC_USERNAME"`
	ElasticPassword  string `envconfig:"ELASTIC_PASSWORD"`
}

func LoadConfig() *Config {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatalf("Failed to load env config: %v", err)
	}
	return &cfg
}

func InitDBAndElastic(cfg *Config) (*sql.DB, *elasticsearch.Client) {
	// Init MySQL
	mysqlDSN := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
		cfg.MySQLUser,
		cfg.MySQLPassword,
		cfg.MySQLHost,
		cfg.MySQLPort,
		cfg.MySQLDBName,
	)

	db, err := sql.Open("mysql", mysqlDSN)
	if err != nil {
		log.Fatalf("Failed to connect to MySQL: %v", err)
	}

	// Ping buat validasi koneksi MySQL
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping MySQL: %v", err)
	}

	// Init Elasticsearch
	esAddresses := []string{cfg.ElasticAddresses}
	esClient, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: esAddresses,
		Username:  cfg.ElasticUsername,
		Password:  cfg.ElasticPassword,
		Transport: &http.Transport{
			ResponseHeaderTimeout: time.Second * 10,
		},
	})

	if err != nil {
		log.Fatalf("Failed to create Elasticsearch client: %v", err)
	}

	// Test koneksi Elasticsearch
	res, err := esClient.Info()
	if err != nil {
		log.Fatalf("Failed to connect to Elasticsearch: %v", err)
	}
	defer res.Body.Close()

	log.Println("Connected to MySQL and Elasticsearch successfully")
	return db, esClient
}
