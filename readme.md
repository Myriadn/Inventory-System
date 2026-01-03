# Inventory System

Aplikasi Inventory Management System berbasis RESTful API menggunakan bahasa pemrograman Golang, bertujuan untuk membantu pengguna dalam mengelola data barang, kategori, rak, gudang, serta proses penjualan barang.

## Cara Menjalankan

1. Buat database PostgreSQL bernama `db_inventory_system`.
2. Jalankan query di file `db_inventory_system.sql`.
3. Sesuaikan koneksi DB di `.env`.
4. Jalankan: `go run cmd/main.go` atau `make run`
5. Lakukan pengujian di Postman atau `http://localhost:8080/ping`

## Cara Testing Coverage

Jalankan perintah di terminal. `go test -coverprofile=coverage.out ./internal/repository/... ./internal/service/...` kemudian `go tool cover -html=coverage.out -o coverage.html`

Atau sebenarnya juga udah siapin scriptnya dengan: `make coverage`

## Video

[Summary Video](github.com)
