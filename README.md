# Proyek E-Commerce gRPC

Proyek ini merupakan **implementasi backend untuk platform e-commerce** yang dirancang untuk mendukung seluruh proses bisnis utama dalam perdagangan digital, mulai dari **pengelolaan pengguna** hingga **transaksi pembayaran**. Sistem ini dibangun dengan pendekatan modular sehingga setiap layanan memiliki tanggung jawab yang spesifik, namun tetap terhubung secara efisien melalui komunikasi antar layanan.

Arsitektur proyek memanfaatkan **gRPC sebagai protokol komunikasi antar layanan**. Dengan gRPC dan Protobuf, pertukaran data antar komponen menjadi lebih cepat, ringan, serta mudah diperluas. Hal ini memungkinkan sistem untuk **skalabel** ketika jumlah pengguna, produk, maupun transaksi meningkat secara signifikan.

Semua data penting — mulai dari detail pengguna, katalog produk, hingga catatan transaksi — disimpan secara konsisten dalam basis data, sehingga platform ini dapat menjadi fondasi yang **andal, aman, dan siap produksi** untuk aplikasi e-commerce modern.

### 🎯 Lingkup & Fitur Utama

* **🔐 Autentikasi & Manajemen Pengguna**
  Mendukung registrasi, login, serta otorisasi berbasis **JWT**. Sistem juga memungkinkan pengaturan role pengguna (misalnya customer, merchant, atau admin).

* **📦 Manajemen Produk & Kategori**
  Admin maupun merchant dapat menambahkan, mengubah, atau menghapus produk. Produk dikategorikan untuk mempermudah pencarian dan pengelolaan stok.

* **🏬 Manajemen Merchant**
  Merchant dapat membuat akun, mengelola toko mereka, serta menambahkan katalog produk yang dimiliki.

* **🛒 Keranjang Belanja (Shopping Cart)**
  Pengguna dapat menambahkan produk ke keranjang, memperbarui jumlah, dan menyimpan daftar item sebelum melakukan pemesanan.

* **📑 Proses Pemesanan (Order)**
  Mendukung alur penuh dari checkout hingga pencatatan pesanan. Setiap order dicatat dengan detail produk, jumlah, harga, status, dan informasi pengiriman.

* **💳 Transaksi & Pembayaran**
  Sistem mencatat pembayaran yang dilakukan oleh pengguna terhadap pesanan mereka. Proses ini dirancang agar **terintegrasi dengan merchant** dan tercatat secara transparan.

* **⭐ Ulasan Produk**
  Setelah pesanan selesai, pengguna dapat memberikan ulasan dan rating terhadap produk yang dibeli, sehingga platform lebih interaktif dan mendukung kepercayaan konsumen.

---

## 🧰 Tech Teknologi

- 🐹 **Go (Golang)** — Bahasa implementasi.
- 🌐 **Echo** — Kerangka kerja web minimalis untuk membangun REST API.
- 🪵 **Zap Logger** — Pencatatan terstruktur untuk aplikasi berkinerja tinggi.
- 📦 **SQLC** — Menghasilkan kode Go yang aman dari tipe dari kueri SQL.
- 🚀 **gRPC** — RPC berkinerja tinggi untuk komunikasi layanan internal.
- 🧳 **Goose** — Alat migrasi untuk mengelola perubahan skema database.
- 🐳 **Docker** — Platform kontainerisasi untuk lingkungan pengembangan yang konsisten.
- 📄 **Swago** — Menghasilkan dokumentasi Swagger 2.0 untuk rute Echo.
- 🔗 **Docker Compose** — Mengelola aplikasi Docker multi-kontainer.

---

## Entity-Relationship Diagram (ERD)

Berikut adalah ERD yang menggambarkan skema database dari proyek ini, dirender menggunakan Mermaid.js.

```mermaid
erDiagram
    users {
        INT user_id PK
        VARCHAR firstname
        VARCHAR lastname
        VARCHAR email UK
        VARCHAR password
    }

    roles {
        INT role_id PK
        VARCHAR role_name UK
    }

    user_roles {
        INT user_role_id PK
        INT user_id FK
        INT role_id FK
    }

    refresh_tokens {
        INT refresh_token_id PK
        INT user_id FK
        VARCHAR token UK
        TIMESTAMP expiration
    }

    merchants {
        INT merchant_id PK
        INT user_id FK
        VARCHAR name
        TEXT description
        VARCHAR status
    }

    categories {
        INT category_id PK
        VARCHAR name
        VARCHAR slug_category UK
    }

    products {
        INT product_id PK
        INT merchant_id FK
        INT category_id FK
        VARCHAR name
        INT price
        INT count_in_stock
        VARCHAR slug_product UK
    }

    carts {
        INT cart_id PK
        INT user_id FK
        INT product_id FK
        INT quantity
    }

    orders {
        INT order_id PK
        INT user_id FK
        INT merchant_id FK
        INT total_price
    }

    order_items {
        INT order_item_id PK
        INT order_id FK
        INT product_id FK
        INT quantity
        INT price
    }

    shipping_addresses {
        INT shipping_address_id PK
        INT order_id FK
        TEXT alamat
        VARCHAR courier
        DECIMAL shipping_cost
    }

    transactions {
        INT transaction_id PK
        INT order_id FK
        INT merchant_id FK
        VARCHAR payment_method
        INT amount
        VARCHAR payment_status
    }

    reviews {
        INT review_id PK
        INT user_id FK
        INT product_id FK
        TEXT comment
        INT rating
    }

    review_details {
        INT review_detail_id PK
        INT review_id FK
        VARCHAR type
        TEXT url
    }

    merchant_business_information {
        INT merchant_business_info_id PK
        INT merchant_id FK
        VARCHAR business_type
        VARCHAR tax_id
    }

    merchant_certifications_and_awards {
        INT merchant_certification_id PK
        INT merchant_id FK
        VARCHAR title
        VARCHAR issued_by
    }

    merchant_policies {
        INT merchant_policy_id PK
        INT merchant_id FK
        VARCHAR policy_type
        VARCHAR title
    }

    merchant_details {
        INT merchant_detail_id PK
        INT merchant_id FK
        VARCHAR display_name
    }

    merchant_social_media_links {
        INT merchant_social_id PK
        INT merchant_detail_id FK
        VARCHAR platform
        TEXT url
    }

    sliders {
        INT slider_id PK
        VARCHAR name
        VARCHAR image
    }

    banners {
        INT banner_id PK
        VARCHAR name
        DATE start_date
        DATE end_date
        BOOLEAN is_active
    }

    users ||--o{ user_roles : "has"
    roles ||--o{ user_roles : "has"
    users ||--o{ refresh_tokens : "has"
    users ||--o{ merchants : "owns"
    users ||--o{ carts : "has"
    users ||--o{ orders : "places"
    users ||--o{ reviews : "writes"

    merchants ||--o{ products : "sells"
    merchants ||--o{ orders : "receives"
    merchants ||--o{ transactions : "processes"
    merchants ||--o{ merchant_business_information : "has"
    merchants ||--o{ merchant_certifications_and_awards : "has"
    merchants ||--o{ merchant_policies : "defines"
    merchants ||--o{ merchant_details : "has"

    merchant_details ||--o{ merchant_social_media_links : "has"

    categories ||--o{ products : "contains"

    products ||--o{ carts : "added to"
    products ||--o{ order_items : "included in"
    products ||--o{ reviews : "has"

    orders ||--o{ order_items : "contains"
    orders ||--o{ shipping_addresses : "ships to"
    orders ||--o{ transactions : "has"

    reviews ||--o{ review_details : "has"
```

## Memulai

Untuk memulai proyek ini, ikuti langkah-langkah berikut:

### 1. Clone Repositori

```bash
git clone https://github.com/MamangRust/ecommerce-grpc.git
cd ecommerce-grpc
```

### 2. Prasyarat

- Go (versi 1.20+)
- Docker & Docker Compose
- `make`
- `protoc`

### 3. Konfigurasi

Salin file `.env.example` menjadi `.env` dan sesuaikan variabel lingkungan jika diperlukan. Untuk menjalankan dengan Docker, gunakan `docker.env`.

## Cara Menjalankan Proyek

Anda bisa menjalankan proyek ini menggunakan Docker (direkomendasikan) atau secara lokal dengan Go.

### 1. Menjalankan dengan Docker

Cara termudah untuk memulai adalah dengan Docker Compose. Perintah ini akan membangun image, menjalankan database, menerapkan migrasi, dan memulai server gRPC serta client HTTP.

```bash
# Menjalankan semua layanan di background
make docker-up

# Menghentikan semua layanan
make docker-down
```

Layanan yang akan berjalan:
- `postgres`: Database PostgreSQL di port `5432`.
- `server`: Server gRPC di port `50051`.
- `client`: Klien HTTP (Gateway) di port `5000`.

### 2. Menjalankan Secara Lokal

Jika Anda tidak menggunakan Docker, Anda bisa menjalankan setiap bagian secara manual.

```bash
# 1. Terapkan migrasi database
make migrate

# 2. Jalankan server gRPC
make run-server

# 3. (Di terminal lain) Jalankan klien/gateway HTTP
make run-client
```

### Perintah `make` Lainnya

- `make generate-proto`: Membuat ulang kode Go dari file `.proto`.
- `make lint`: Menjalankan linter pada kode.
- `make test`: Menjalankan unit test.
- `make sqlc-generate`: Membuat ulang kode dari query SQL menggunakan `sqlc`.
