# Simple Elasticsearch & mysql

API sederhana untuk mengelola produk (create, update, get, search, delete) menggunakan Go, MySQL, dan Elasticsearch.

## Prasyarat

Pastikan Anda sudah menginstal:

- **MySQL**
- **Elasticsearch**
- **Go** (minimal versi 1.18)

## Instalasi

1. Clone repository ini:
```bash
git clone https://github.com/Imammahmuda1804/simple_mysql-elasticsearch.git
cd simple_mysql_elasticsearch
```

2. Buat database di MySQL:
```sql
CREATE DATABASE sales;
```

3. Buat tabel `products`:
```sql
CREATE TABLE products (
    id          INT AUTO_INCREMENT PRIMARY KEY,
    sku         VARCHAR(100) NOT NULL UNIQUE,
    name        VARCHAR(255) NOT NULL,
    description TEXT,
    price       DECIMAL(10,2) NOT NULL,
    category    VARCHAR(100),
    image_path  VARCHAR(255),
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
```

4. Jalankan aplikasinya:
```bash
go mod init simple_mysql_elasticsearch
go mod tidy
go mod vendor
go run cmd/main.go
```

## Endpoint API

### Create Product

```bash
curl --location 'http://localhost:8085/products' \
--header 'Content-Type: application/json' \
--data '{
    "sku" :"SAM-2",
    "name": "Laptop Samsung gamers",
    "description": "Laptop dengan spesifikasi tinggi untuk gaming.",
    "price": 1500,
    "category": "Elektronik",
    "images":""
}'
```

### Update Product

```bash
curl --location --request PUT 'http://localhost:8085/products' \
--header 'Content-Type: application/json' \
--data '{
    "id" :1,
    "sku" :"SAM-2",
    "name": "Laptop Samsung untuk gamers ja",
    "description": "Laptop dengan spesifikasi tinggi untuk gaming.",
    "price": 1000000,
    "category": "Elektronik",
     "images":""
}'
```

### Get All Products

```bash
curl --location --request GET 'http://localhost:8085/products' \
--header 'Content-Type: application/json'
```

### Search Product

```bash
curl --location 'http://localhost:8085/products/search?keyword=laptop'
```

### Delete Product

```bash
curl --location --request DELETE 'http://localhost:8085/products/1'
```
