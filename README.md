# Evermos Ecommerce API

REST API e-commerce berbasis Go untuk mengelola akun, toko, alamat, produk, dan transaksi. Aplikasi menggunakan MySQL untuk data utama dan Redis untuk sesi autentikasi serta cache produk.

## Fitur

- Registrasi, login, refresh token, dan logout.
- Sesi multi-perangkat yang disimpan di Redis dengan TTL.
- Pengelolaan profil, toko, alamat, produk, dan transaksi.
- Pengelolaan akun dan kategori khusus admin.
- Cache Redis untuk daftar dan detail produk.
- Endpoint wilayah provinsi dan kota.
- Upload file.

## Tech Stack

<p align="center">
  <a href="https://go.dev/"><img alt="Go" src="https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white"></a>
  <a href="https://www.mysql.com/"><img alt="MySQL" src="https://img.shields.io/badge/MySQL-4479A1?style=for-the-badge&logo=mysql&logoColor=white"></a>
  <a href="https://redis.io/"><img alt="Redis" src="https://img.shields.io/badge/Redis-FF4438?style=for-the-badge&logo=redis&logoColor=white"></a>
  <a href="https://www.docker.com/"><img alt="Docker" src="https://img.shields.io/badge/Docker-2496ED?style=for-the-badge&logo=docker&logoColor=white"></a>
  <a href="https://k6.io/"><img alt="k6" src="https://img.shields.io/badge/k6-7D64FF?style=for-the-badge&logo=k6&logoColor=white"></a>
  <a href="https://gorm.io/"><img alt="GORM" src="https://img.shields.io/badge/GORM-00ADD8?style=for-the-badge&logo=go&logoColor=white"></a>
  <a href="https://gofiber.io/"><img alt="GoFiber" src="https://img.shields.io/badge/GoFiber-00ADD8?style=for-the-badge&logo=go&logoColor=white"></a>
</p>

Backend menggunakan Go 1.26+, Fiber v2, GORM, dan MySQL 8. Redis menyediakan sesi dan cache; Docker Compose menjalankan layanan; k6 digunakan untuk load testing.

## Menjalankan dengan Docker Compose

1. Buat file `.env` di root proyek:

   ```env
   PORT=8000
   DB_HOST=localhost
   DB_PORT=3306
   DB_USER=root
   DB_PASS=ganti_password_ini
   DB_NAME=evermos
   REDIS_HOST=localhost
   REDIS_PORT=6379
   APP_SECRET=ganti_dengan_secret_panjang_dan_acak
   APP_NAME=evermos-ecommerce
   ```

   Compose mengatur `DB_HOST=db` dan `REDIS_HOST=redis` untuk container aplikasi. Nilai host di atas dipakai saat aplikasi dijalankan langsung di komputer.

2. Jalankan layanan:

   ```bash
   docker compose up --build
   ```

   API tersedia di `http://localhost:8000`. Adminer tersedia di `http://localhost:8080`.

Untuk menjalankan aplikasi Go langsung, pastikan MySQL dan Redis aktif, atur `DB_HOST` dan `REDIS_HOST` di `.env` ke `localhost`, lalu jalankan:

```bash
go run .
```

## Autentikasi

Login mengembalikan access token dan refresh token. Sertakan access token untuk endpoint terlindungi:

```http
Authorization: Bearer <access_token>
```

Sesi disimpan dengan key `session:{user_id}:{token_id}` dan kedaluwarsa mengikuti refresh token. Login dari beberapa perangkat membuat sesi terpisah. Refresh token digunakan pada `POST /auth/refresh`, sedangkan `POST /auth/logout` menghapus sesi aktif tersebut.

## Endpoint

| Method | Endpoint | Akses |
| --- | --- | --- |
| `POST` | `/auth/register` | Publik |
| `POST` | `/auth/login` | Publik |
| `POST` | `/auth/refresh` | Refresh token |
| `POST` | `/auth/logout` | Access token |
| `GET`, `PUT` | `/auth/profile` | Access token |
| `GET`, `POST` | `/api/accounts` | Admin |
| `GET`, `PUT`, `DELETE` | `/api/accounts/:id` | Admin |
| `POST` | `/api/categories` | Admin |
| `PUT`, `DELETE` | `/api/categories/:id` | Admin |
| `GET`, `POST` | `/api/products` | Access token, daftar sesuai kepemilikan toko |
| `GET`, `PUT`, `DELETE` | `/api/products/:id` | Access token, sesuai kepemilikan toko |
| `GET` | `/api/addresses` | Access token, data milik sendiri |
| `POST` | `/api/address` | Access token |
| `PUT`, `DELETE` | `/api/addresses/:id` | Access token, data milik sendiri |
| `GET`, `POST`, `PUT` | `/api/store` | Access token, toko milik sendiri |
| `GET`, `POST` | `/api/transactions` | Access token, data transaksi milik sendiri |
| `GET` | `/provcity/listprovincies` | Publik |
| `GET` | `/provcity/listcities/:prov_id` | Publik |
| `GET` | `/provcity/detailprovince/:prov_id` | Publik |
| `GET` | `/provcity/detailcity/:city_id` | Publik |
| `POST` | `/api/upload` | Access token |

Operasi `GET /api/products` dan `GET /api/products/:id` menggunakan cache Redis selama 10 detik. Query daftar produk menerima parameter `search`, `sort`, `limit`, dan `offset`.

## Load testing

Script di `k6/load-test.js` mencakup request yang membuat atau mengubah data, termasuk upload dan transaksi. Jalankan hanya terhadap lingkungan uji dengan akun serta data khusus testing:

```bash
k6 run k6/load-test.js
```

Konfigurasi URL, kredensial, tahapan beban, dan data uji berada di bagian atas script dan fungsi `setup()`.

## Struktur Proyek

```text
.
├── cmd/                    # Routing dan dependency injection
├── internal/
│   ├── config/             # Setup database, Redis, dan seeder
│   ├── controllers/        # HTTP handlers
│   ├── interfaces/         # Kontrak service dan repository
│   ├── models/             # Model data dan token
│   ├── repository/         # Akses database
│   └── services/           # Logika bisnis
├── k6/                     # Script load test dan fixture upload
├── middleware/             # Autentikasi dan otorisasi
├── utils/                  # Helper umum dan JWT
├── docker-compose.yml
├── Dockerfile
├── go.mod
└── main.go
```
