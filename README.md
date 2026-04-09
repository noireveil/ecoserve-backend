# EcoServe Backend

EcoServe adalah platform manajemen siklus hidup elektronik yang dirancang untuk memperpanjang usia perangkat dan menekan pertumbuhan limbah elektronik (*e-waste*) dalam kerangka ekonomi sirkular. API ini dibangun menggunakan **Go** dengan fokus pada performa tinggi dan akurasi data geospasial.

## Overview Fitur

* **AI Triage System:** Integrasi **Gemini 2.5 Flash** untuk mendiagnosis kerusakan perangkat secara otomatis dengan mekanisme *Confidence Gate* untuk meminimalkan halusinasi AI.
* **Geospatial Search:** Memanfaatkan ekstensi **PostGIS** untuk melakukan pencarian teknisi terdekat berdasarkan radius geografis secara akurat.
* **Automated E-Waste Tracker:** Kalkulasi otomatis penghematan emisi dan limbah menggunakan algoritma **EPA WARM v15** setiap kali transaksi selesai.
* **Anti-Fraud Layer:** Protokol keamanan yang mencakup *GPS Locking*, bukti visual (*Photo Proof*), dan konfirmasi ganda (*Dual Confirmation*) antara konsumen dan teknisi.
* **Secure Authentication:** Autentikasi nir-sandi melalui **Email OTP** dan otorisasi berbasis **JWT** yang dikelola secara aman.

## Tumpukan Teknologi (Tech Stack)

* **Language:** Go (Golang) 1.25
* **Web Framework:** Fiber v2
* **Database:** PostgreSQL (Hosted on Supabase) + PostGIS
* **AI Engine:** Google Gemini 2.5 Flash
* **Auth:** Supabase Auth (Email OTP)
* **Deployment:** Render & Vercel Edge Network

## Arsitektur Proyek

Mengadopsi *Layered Microservices Architecture* untuk memastikan pemisahan logika bisnis dari infrastruktur:

* `cmd/api/`: Titik masuk utama aplikasi (*entry point*).
* `internal/domain/`: Definisi entitas, model basis data, dan Digital Product Passport (DPP).
* `internal/usecase/`: Logika bisnis inti (Kalkulasi emisi, *AI Prompt Engineering*).
* `internal/repository/`: Lapisan akses data (GORM & SQL Spasial).
* `internal/delivery/http/`: Handler REST API dan *middleware* keamanan.
* `pkg/utils/`: Utilitas pembantu (JWT, SMTP Mailer, OTP Generator).

## Dokumentasi API (OpenAPI)

Sesuai standar SRS, dokumentasi API tersedia secara interaktif melalui Swagger UI. Anda dapat mengaksesnya saat server berjalan di:
`http://localhost:3000/swagger/`

## Panduan Setup Lokal

1.  **Clone Repositori:**
    ```bash
    git clone https://github.com/noireveil/ecoserve-backend.git
    cd ecoserve-backend
    ```
2.  **Konfigurasi Environment:**
    Salin fail `.env.example` menjadi `.env` dan isi variabel yang diperlukan (Database, SMTP, dan Gemini API Key).
3.  **Jalankan dengan Docker:**
    Sistem akan otomatis menjalankan PostgreSQL dengan ekstensi PostGIS.
    ```bash
    docker compose up -d --build
    ```
4.  **Verifikasi:**
    Akses `http://localhost:3000/health` untuk memastikan konektivitas API dan basis data telah aktif.