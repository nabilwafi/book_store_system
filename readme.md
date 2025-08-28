# Book Store System

Sistem manajemen toko buku berbasis gRPC yang dibangun dengan Go, menggunakan arsitektur clean architecture dan PostgreSQL sebagai database.

## 🚀 Fitur Utama

- **Manajemen Pengguna**: Registrasi, login, dan autentikasi berbasis JWT
- **Manajemen Kategori**: CRUD operasi untuk kategori buku
- **Manajemen Buku**: CRUD operasi untuk buku dengan kategori
- **Sistem Pemesanan**: Pembuatan order, manajemen status, dan pembayaran
- **Laporan**: Laporan penjualan, buku terlaris, dan statistik harga
- **Autentikasi & Autorisasi**: Role-based access control (Admin/User)

## 🏗️ Arsitektur

Proyek ini menggunakan **Clean Architecture** dengan struktur sebagai berikut:

```
book-store-system/
├── cmd/server/           # Entry point aplikasi
├── config/              # Konfigurasi aplikasi
├── internal/
│   ├── entity/          # Domain entities
│   ├── repository/      # Data access layer
│   ├── service/         # Business logic layer
│   └── transport/       # Presentation layer
│       ├── dto/         # Data Transfer Objects
│       ├── grpc/        # gRPC handlers
│       └── http/        # HTTP handlers (future)
├── pkg/                 # Shared utilities
│   ├── database/        # Database connection
│   ├── helpers/         # Helper functions
│   ├── logger/          # Logging utilities
│   └── middleware/      # Middleware functions
├── proto/               # Protocol Buffer definitions
└── test/                # Test clients
```

## 🛠️ Teknologi yang Digunakan

- **Go 1.21+**: Bahasa pemrograman utama
- **gRPC**: Framework komunikasi antar service
- **Protocol Buffers**: Serialisasi data
- **PostgreSQL**: Database relational
- **GORM**: ORM untuk Go
- **JWT**: Autentikasi dan autorisasi
- **Docker**: Containerization
- **Docker Compose**: Orchestration

## 📋 Prasyarat

- Go 1.21 atau lebih baru
- PostgreSQL 13+
- Docker & Docker Compose (opsional)
- Protocol Buffer Compiler (protoc)

## 🚀 Instalasi & Setup

### 1. Clone Repository

```bash
git clone <repository-url>
cd book-store-system
```

### 2. Setup Environment

Salin file environment dan sesuaikan konfigurasi:

```bash
cp .env.example .env
```

Edit file `.env` sesuai dengan konfigurasi database Anda:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=your_username
DB_PASSWORD=your_password
DB_NAME=bookstore
JWT_SECRET=your_jwt_secret
SERVER_PORT=50051
```

### 3. Install Dependencies

```bash
go mod download
```

### 4. Setup Database

Buat database PostgreSQL:

```sql
CREATE DATABASE bookstore;
```

### 5. Generate Protocol Buffers (jika diperlukan)

```bash
protoc --go_out=. --go-grpc_out=. proto/bookstore.proto
```

### 6. Jalankan Aplikasi

```bash
go run cmd/server/main.go
```

Atau menggunakan Docker Compose:

```bash
docker-compose up -d
```

## 📖 Penggunaan

### gRPC Services

Aplikasi menyediakan beberapa service gRPC:

#### 1. User Service
- `Register`: Registrasi pengguna baru
- `Login`: Autentikasi pengguna
- `GetProfile`: Mendapatkan profil pengguna

#### 2. Category Service
- `CreateCategory`: Membuat kategori baru (Admin only)
- `GetCategories`: Mendapatkan daftar kategori
- `GetCategory`: Mendapatkan detail kategori
- `UpdateCategory`: Memperbarui kategori (Admin only)
- `DeleteCategory`: Menghapus kategori (Admin only)

#### 3. Book Service
- `CreateBook`: Membuat buku baru (Admin only)
- `GetBooks`: Mendapatkan daftar buku dengan pagination
- `GetBook`: Mendapatkan detail buku
- `GetBooksByCategory`: Mendapatkan buku berdasarkan kategori
- `UpdateBook`: Memperbarui buku (Admin only)
- `DeleteBook`: Menghapus buku (Admin only)

#### 4. Order Service
- `CreateOrder`: Membuat pesanan baru
- `GetOrders`: Mendapatkan daftar pesanan
- `GetOrder`: Mendapatkan detail pesanan
- `UpdateOrderStatus`: Memperbarui status pesanan (Admin only)
- `ProcessPayment`: Memproses pembayaran

#### 5. Report Service
- `GetSalesReport`: Laporan penjualan berdasarkan periode
- `GetTopBooks`: Laporan buku terlaris
- `GetBookPriceStatistics`: Statistik harga buku (min, max, rata-rata)

### Test Clients

Proyek ini menyediakan beberapa test client untuk pengujian:

```bash
# Test basic functionality
go run test/grpc_client.go

# Test CRUD operations
go run test/grpc_crud_client.go

# Test order functionality
go run test/grpc_order_client.go

# Test price statistics
go run test/grpc_price_stats_client.go
```

## 🔐 Autentikasi

Sistem menggunakan JWT untuk autentikasi. Setiap request yang memerlukan autentikasi harus menyertakan token JWT dalam field `token`.

### Role-based Access Control

- **Admin**: Akses penuh ke semua operasi
- **User**: Akses terbatas (tidak bisa CRUD kategori/buku, tidak bisa mengubah status order)

## 📊 Database Schema

### Users
- `id`: Primary key
- `username`: Unique username
- `email`: Unique email
- `password`: Hashed password
- `role`: User role (admin/user)
- `created_at`, `updated_at`, `deleted_at`: Timestamps

### Categories
- `id`: Primary key
- `name`: Category name
- `description`: Category description
- `created_at`, `updated_at`, `deleted_at`: Timestamps

### Books
- `id`: Primary key
- `title`: Book title
- `author`: Book author
- `isbn`: Unique ISBN
- `price`: Book price
- `stock`: Available stock
- `category_id`: Foreign key to categories
- `created_at`, `updated_at`, `deleted_at`: Timestamps

### Orders
- `id`: Primary key
- `user_id`: Foreign key to users
- `total_amount`: Total order amount
- `status`: Order status
- `payment_status`: Payment status
- `created_at`, `updated_at`, `deleted_at`: Timestamps

### Order Items
- `id`: Primary key
- `order_id`: Foreign key to orders
- `book_id`: Foreign key to books
- `quantity`: Item quantity
- `price`: Item price at time of order
- `created_at`, `updated_at`, `deleted_at`: Timestamps

## 🧪 Testing

Untuk menjalankan test:

```bash
go test ./...
```

Untuk test dengan coverage:

```bash
go test -cover ./...
```

## 📝 API Documentation

Dokumentasi lengkap API dapat ditemukan dalam file Protocol Buffer di `proto/bookstore.proto`.

### Contoh Request/Response

#### Login
```protobuf
// Request
{
  "username": "admin",
  "password": "password123"
}

// Response
{
  "success": true,
  "message": "Login successful",
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "user": {
    "id": 1,
    "username": "admin",
    "email": "admin@example.com",
    "role": "admin"
  }
}
```

#### Create Book
```protobuf
// Request
{
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "title": "The Go Programming Language",
  "author": "Alan Donovan",
  "isbn": "978-0134190440",
  "price": 45.99,
  "stock": 100,
  "category_id": 1
}

// Response
{
  "success": true,
  "message": "Book created successfully",
  "book": {
    "id": 1,
    "title": "The Go Programming Language",
    "author": "Alan Donovan",
    "isbn": "978-0134190440",
    "price": 45.99,
    "stock": 100,
    "category_id": 1
  }
}
```

## 🐳 Docker Deployment

Untuk deployment menggunakan Docker:

```bash
# Build dan jalankan dengan Docker Compose
docker-compose up -d

# Lihat logs
docker-compose logs -f

# Stop services
docker-compose down
```

## 🤝 Contributing

1. Fork repository
2. Buat feature branch (`git checkout -b feature/amazing-feature`)
3. Commit perubahan (`git commit -m 'Add amazing feature'`)
4. Push ke branch (`git push origin feature/amazing-feature`)
5. Buat Pull Request

## 📄 License

Proyek ini dilisensikan di bawah MIT License - lihat file [LICENSE](LICENSE) untuk detail.

## 📞 Support

Jika Anda memiliki pertanyaan atau masalah, silakan buat issue di repository ini.

---

**Dibuat dengan ❤️ menggunakan Go dan gRPC**