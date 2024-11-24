# Ecommerce API

Ini adalah API ecommerce yang dibangun dengan Go menggunakan paket routing `go-chi`, yang memungkinkan Anda untuk mengelola produk, pengguna, dan pesanan. API ini menyediakan endpoint untuk otentikasi pengguna, manajemen produk, keranjang belanja, dan pembuatan pesanan.

## Fitur

- Otentikasi pengguna (login, registrasi, logout)
- Manajemen akun pengguna (manajemen alamat, pembaruan detail pengguna)
- Katalog produk (produk terlaris, promosi, rekomendasi)
- Fungsi keranjang belanja (menambah, memperbarui, menghapus item)
- Manajemen daftar keinginan
- Pembuatan pesanan dan proses checkout
- Rute API dilindungi dengan middleware otentikasi

## Teknologi yang Digunakan

- Go (Golang)
- Chi Router
- Zap Logger
- Middleware Kustom
- Uber Zap untuk pencatatan terstruktur

## Setup

1. Clone repository ini ke mesin lokal Anda:

   ```bash
   git clone https://github.com/usernameanda/ecommerce-api.git
   cd ecommerce-api
   ```

2. Install dependensi:

   ```bash
   go mod tidy
   ```

3. Buat file konfigurasi `config.env` atau atur variabel lingkungan untuk nilai-nilai konfigurasi yang dibutuhkan oleh aplikasi.

4. Jalankan aplikasi:

   ```bash
   go run main.go
   ```

5. API akan dijalankan secara default pada port `8080`. Anda dapat mengubahnya di konfigurasi.

## Endpoint

### Endpoint Otentikasi

- **POST** `/login` - Login pengguna
- **POST** `/register` - Registrasi pengguna
- **POST** `/logout` - Logout pengguna

### Endpoint Akun (Dilindungi)

- **GET** `/api/account/address` - Mendapatkan semua alamat pengguna
- **GET** `/api/account/detail-user` - Mendapatkan detail pengguna
- **PUT** `/api/account/update-user` - Memperbarui informasi pengguna
- **POST** `/api/account/address` - Membuat alamat baru

### Endpoint Produk

- **GET** `/api/products/` - Mendapatkan semua produk
- **GET** `/api/products/best-selling` - Mendapatkan produk terlaris
- **GET** `/api/products/weekly-promotion` - Mendapatkan produk dengan promosi mingguan
- **GET** `/api/products/recomments` - Mendapatkan produk rekomendasi
- **GET** `/api/products/{id}` - Mendapatkan produk berdasarkan ID

### Endpoint Keranjang (Dilindungi)

- **GET** `/api/products/carts` - Mendapatkan semua item di keranjang
- **POST** `/api/products/carts` - Menambah item ke keranjang
- **PUT** `/api/products/carts/{id}` - Memperbarui jumlah item di keranjang
- **DELETE** `/api/products/carts/{id}` - Menghapus item dari keranjang
- **GET** `/api/products/total-carts` - Mendapatkan jumlah total item di keranjang

### Endpoint Daftar Keinginan (Dilindungi)

- **POST** `/api/products/wishlist` - Menambah item ke daftar keinginan
- **DELETE** `/api/products/wishlist/{id}` - Menghapus item dari daftar keinginan

### Endpoint Kategori

- **GET** `/api/categories` - Mendapatkan semua kategori produk

### Endpoint Banner

- **GET** `/api/banners` - Mendapatkan semua banner

## Middleware

API ini menggunakan middleware otentikasi kustom (`authMiddleware`) untuk semua endpoint yang berhubungan dengan akun pengguna, keranjang, dan pesanan. Middleware ini memastikan hanya pengguna yang terotentikasi yang dapat mengakses endpoint-endpoint tersebut.

## Logging

Pencatatan terstruktur ditangani dengan logger `zap` dari Uber. Aplikasi mencatat hal-hal berikut:

- Informasi tentang konfigurasi saat aplikasi dimulai
- Detail permintaan dan respons HTTP, termasuk metode, URL, dan durasi

## Konfigurasi

Konfigurasi API dimuat dari file konfigurasi (misalnya, `config.env`). Anda dapat menyesuaikan konfigurasi untuk lingkungan yang berbeda seperti pengembangan atau produksi.
