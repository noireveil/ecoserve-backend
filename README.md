# EcoServe Backend API

Layanan backend untuk platform EcoServe. API ini dirancang untuk melayani aplikasi klien (PWA) dalam menghubungkan masyarakat dengan teknisi elektronik terdekat guna mendukung ekonomi sirkular dan pengurangan limbah elektronik (E-Waste).

## Tech Stack
- **Language**: Go (Golang) 1.25
- **Framework**: Fiber v2
- **Database**: PostgreSQL 15 dengan ekstensi PostGIS 3.3 (untuk query spasial radius)
- **ORM**: GORM
- **Infrastruktur**: Docker & Docker Compose

## Struktur Proyek
Aplikasi ini menggunakan pendekatan arsitektur modular (*Standard Go Layout*):
- `internal/domain`: Entitas inti dan representasi tabel database.
- `internal/repository`: Lapisan abstraksi untuk interaksi langsung dengan database.
- `internal/usecase`: Tempat berkumpulnya logika bisnis (seperti kalkulasi reduksi E-Waste).
- `internal/delivery`: Menangani *routing* dan HTTP *request/response*.

## Cara Menjalankan Aplikasi (Local Development)

1. Pastikan Docker dan Git sudah terinstal di sistem Anda.
2. Lakukan *clone* repositori:
   ```bash
   git clone https://github.com/noireveil/ecoserve-backend.git
   cd ecoserve-backend
   ```
3. Salin file *environment*:
   ```bash
   cp .env.example .env
   ```
4. Jalankan infrastruktur melalui Docker Compose:
   ```bash
   docker compose up -d --build
   ```
5. API akan berjalan di `http://localhost:3000`.

## Catatan Pengembangan
Fase saat ini memuat struktur inti *database* dan *routing* dasar. Integrasi otentikasi OTP dan pengamanan *endpoint* (Middleware) akan diimplementasikan pada fase berikutnya sesuai kebutuhan integrasi dengan antarmuka klien.
