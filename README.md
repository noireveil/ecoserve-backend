# EcoServe Backend

EcoServe adalah platform yang dirancang untuk mempermudah akses servis elektronik guna memperpanjang usia perangkat dan menekan pertumbuhan limbah elektronik (*e-waste*). API ini dibangun menggunakan **Go** dengan fokus pada performa, skalabilitas, dan akurasi data geospasial.

## Overview Fitur

* **AI Triage System:** Integrasi **Gemini 2.5 Flash** untuk mendiagnosis kerusakan perangkat secara otomatis berdasarkan input pengguna, lengkap dengan mitigasi bahaya dan estimasi biaya.
* **Geospatial Search:** Menggunakan **PostGIS** untuk pencarian teknisi terdekat berdasarkan radius koordinat (Longitude/Latitude).
* **Automated E-Waste Tracker:** Kalkulasi otomatis penghematan limbah elektronik dalam satuan kilogram (Kg) setiap kali sebuah servis selesai dilakukan.
* **Secure Authentication:** Sistem login tanpa kata sandi menggunakan **Email OTP** (SMTP) dan otorisasi berbasis **JWT**.

## Tech Stack

* **Language:** Go (Golang) 1.25
* **Web Framework:** Fiber v2
* **Database:** PostgreSQL (Supabase) + PostGIS Extension
* **ORM:** GORM
* **AI Engine:** Google Gemini 2.5 Flash
* **Containerization:** Docker & Docker Compose
* **Deployment:** Render

## Arsitektur

Proyek ini menerapkan **Clean Architecture** (Modular Layout) untuk memisahkan logika bisnis dari detail infrastruktur:

* `cmd/`: Entry point aplikasi.
* `internal/domain/`: Definisi struct dan model database.
* `internal/usecase/`: Logika bisnis (E-waste calculation, AI prompt engineering, dll).
* `internal/repository/`: Interaksi data (SQL Query, GORM).
* `internal/delivery/`: HTTP Handlers dan routing.
* `pkg/`: Utility/Helper (JWT, Mailer, Logger).

## Setup Lokal

1.  Clone repositori:
    ```bash
    git clone https://github.com/noireveil/ecoserve-backend.git
    cd ecoserve-backend
    ```
2.  Konfigurasi environment:
    ```bash
    cp .env.example .env
    # Isi variabel seperti DB_URL, SMTP_PASS, dan GEMINI_API_KEY
    ```
3.  Jalankan dengan Docker:
    ```bash
    docker compose up -d --build
    ```
4.  Akses API melalui `http://localhost:3000`.