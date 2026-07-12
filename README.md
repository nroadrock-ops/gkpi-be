# GKPI Cimahi Backend API

Proyek backend REST API untuk Gereja GKPI Cimahi menggunakan Golang, Fiber, GORM, dan PostgreSQL (via Supabase).

## Cara Menjalankan Server Lokal

1. Pastikan Anda sudah menginstal Go dan memiliki koneksi internet.
2. Salin file `.env.example` menjadi `.env` dan isi variabel yang dibutuhkan, terutama `SUPABASE_DB_URL` dan `JWT_SECRET`.
   ```bash
   cp .env.example .env
   ```
3. Jalankan server secara lokal:
   ```bash
   go run cmd/api/main.go
   ```
   Server akan berjalan di `http://localhost:3000`.

## Cara Menjalankan Seed Data

Untuk mengisi database dengan data *dummy* awal (Jemaat, Jadwal Ibadah, Artikel, dll), jalankan perintah berikut:

```bash
go run cmd/api/main.go -seed
```
*Catatan: Menjalankan flag `-seed` akan mengeksekusi migrasi tabel dan mengisi data, lalu program akan langsung berhenti (exit) tanpa menjalankan server HTTP.*

## Import Postman Collection & Environment

Terdapat dua file Postman di root proyek ini untuk mempermudah testing:

1. `gkpi-cimahi.postman_collection.json`
2. `gkpi-cimahi.postman_environment.json`

### Cara Import:
1. Buka aplikasi Postman.
2. Klik tombol **Import** di kiri atas.
3. *Drag and drop* atau pilih kedua file JSON tersebut.
4. Di bagian kanan atas Postman, pilih environment **GKPI Cimahi Local**.
5. Gunakan endpoint **Auth -> Login** untuk mendapatkan token.
6. Salin token dari *response* login dan masukkan ke tab *Variables* di Environment **GKPI Cimahi Local** pada kolom `access_token` (kolom *Current Value*). Semua endpoint yang *protected* akan otomatis menggunakan token tersebut.
