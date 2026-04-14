# EcoServe Backend

EcoServe adalah platform manajemen siklus hidup elektronik yang dirancang untuk memperpanjang usia perangkat dan menekan pertumbuhan limbah elektronik (*e-waste*) dalam kerangka ekonomi sirkular. API ini dibangun menggunakan **Go** dengan fokus pada performa tinggi dan akurasi data geospasial.

## Overview Fitur

* **Digital Product Passport (DPP):** Registrasi dan pencatatan riwayat perangkat elektronik konsumen secara digital untuk memantau status garansi dan rekam jejak perbaikan.
* **AI Triage System:** Integrasi **Gemini 2.5 Flash** untuk mendiagnosis kerusakan perangkat secara otomatis dengan mekanisme *Confidence Gate* untuk meminimalkan halusinasi AI.
* **Geospatial Search:** Memanfaatkan ekstensi **PostGIS** untuk melakukan pencarian teknisi terdekat berdasarkan radius geografis secara akurat.
* **Automated E-Waste Tracker:** Kalkulasi otomatis penghematan emisi dan limbah menggunakan algoritma **EPA WARM v15** setiap kali transaksi selesai.
* **Anti-Fraud & Security Layer:** Protokol keamanan yang mencakup *GPS Locking*, bukti visual (*Photo Proof*), serta **Fiber Rate Limiter** untuk perlindungan terhadap serangan DDoS dan spam API.
* **Secure Authentication:** Autentikasi nir-sandi melalui **Email OTP** yang diintegrasikan via **Mailjet API** dan otorisasi berbasis **JWT**.

## Tumpukan Teknologi (Tech Stack)

* **Language:** Go (Golang) 1.25+
* **Web Framework:** Fiber v2
* **Database:** PostgreSQL (Hosted on Supabase) + PostGIS
* **AI Engine:** Google Gemini 2.5 Flash
* **Email Service:** Mailjet API (HTTP Protocol)
* **Deployment:** Render (Backend) & Vercel (Frontend)

## Arsitektur Proyek

Mengadopsi *Layered Architecture* untuk memastikan pemisahan logika bisnis dari infrastruktur:

* `cmd/api/`: Titik masuk utama aplikasi (*entry point*) dan konfigurasi *middleware*.
* `internal/domain/`: Definisi entitas, model basis data, dan Digital Product Passport (DPP).
* `internal/usecase/`: Logika bisnis inti (Kalkulasi emisi, *AI Prompt Engineering*, manajemen OTP).
* `internal/repository/`: Lapisan akses data (GORM & SQL Spasial).
* `internal/delivery/http/`: Handler REST API dan *middleware* keamanan (*Rate Limiting*, CORS).
* `pkg/utils/`: Utilitas pembantu (JWT, Mailjet API Integration, OTP Generator).

## Dokumentasi API (OpenAPI)

Sesuai standar SRS, dokumentasi API tersedia secara interaktif melalui Swagger UI:
`http://localhost:3000/swagger/`

## Panduan Setup Lokal

1.  **Clone Repositori:**
    ```bash
    git clone https://github.com/noireveil/ecoserve-backend.git
    cd ecoserve-backend
    ```
2.  **Konfigurasi Environment:**
    Salin fail `.env.example` menjadi `.env` dan isi variabel yang diperlukan:
    * `MAILJET_API_KEY` & `MAILJET_SECRET_KEY`
    * `GEMINI_API_KEY`
    * Database Credentials (PostgreSQL)
3.  **Sinkronisasi Dependensi:**
    ```bash
    go mod tidy
    ```
4.  **Jalankan Server:**
    ```bash
    go run cmd/api/main.go
    ```
5.  **Verifikasi:**
    Akses `http://localhost:3000/health` untuk memastikan konektivitas API dan basis data telah aktif.