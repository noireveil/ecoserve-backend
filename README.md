# EcoServe Backend API

Layanan backend untuk platform EcoServe. API ini dirancang untuk melayani aplikasi klien (PWA) dalam menghubungkan masyarakat dengan teknisi elektronik terdekat guna mendukung ekonomi sirkular dan pengurangan limbah elektronik (E-Waste).

## Base URL (Production)
API telah diluncurkan dan dapat diakses secara publik melalui:
`https://ecoserve-api.onrender.com`

## Tech Stack
- **Language**: Go (Golang) 1.25
- **Framework**: Fiber v2
- **Database**: PostgreSQL (Supabase) dengan ekstensi PostGIS (untuk kueri spasial radius)
- **ORM**: GORM
- **Infrastruktur**: Docker & Docker Compose
- **Deployment**: Render (Web Service)

## Struktur Proyek
Aplikasi ini menggunakan pendekatan arsitektur modular (*Standard Go Layout*):
- `internal/domain`: Entitas inti dan representasi tabel database.
- `internal/repository`: Lapisan abstraksi untuk interaksi langsung dengan database.
- `internal/usecase`: Tempat berkumpulnya logika bisnis (seperti kalkulasi reduksi E-Waste).
- `internal/delivery`: Menangani *routing* dan HTTP *request/response* (termasuk *Middleware* otentikasi).

## Cara Menjalankan Aplikasi (Local Development)

1. Pastikan Docker dan Git sudah terinstal di sistem.
2. Lakukan *clone* repositori:
   ```bash
   git clone https://github.com/noireveil/ecoserve-backend.git
   cd ecoserve-backend
   ```
3. Salin berkas *environment*:
   ```bash
   cp .env.example .env
   ```
4. Jalankan infrastruktur melalui Docker Compose:
   ```bash
   docker compose up -d --build
   ```
5. API lokal akan berjalan di `http://localhost:3000`.

## Status Pengembangan Terkini
- Struktur inti basis data dan *routing* dasar telah selesai.
- Integrasi otentikasi berbasis OTP surel (SMTP) telah beroperasi.
- Pengamanan *endpoint* privat menggunakan otorisasi token JWT.
- Peluncuran infrastruktur awan terhubung secara berkelanjutan (*Continuous Deployment*) ke platform Render dan Supabase.
