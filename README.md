# EcoServe Backend Engine

EcoServe adalah infrastruktur perangkat lunak yang dirancang untuk mendigitalkan ekonomi sirkular dan memperpanjang usia perangkat elektronik. API ini menjembatani konsumen dan teknisi perbaikan melalui lapisan kecerdasan buatan, pemetaan geospasial presisi, dan pelacakan jejak karbon (EPA WARM v15).

## Table of Contents
- [EcoServe Backend Engine](#ecoserve-backend-engine)
  - [Table of Contents](#table-of-contents)
  - [Core Architecture](#core-architecture)
  - [Tech Stack](#tech-stack)
  - [API Documentation](#api-documentation)
  - [Local Environment Setup](#local-environment-setup)
    - [Option 1: Docker (Recommended)](#option-1-docker-recommended)
    - [Option 2: Manual Installation](#option-2-manual-installation)
  - [Environmental Impact Analysis](#environmental-impact-analysis)

## Core Architecture

Sistem ini dibangun di atas fondasi Domain-Driven Design menggunakan bahasa Go (Golang), memastikan skalabilitas tinggi, kemudahan pengujian, dan keterisolasian logika bisnis.

* **Digital Product Passport (DPP):** Manajemen siklus hidup elektronik (*Electronic Lifecycle Management*) dengan rekam jejak kepemilikan dan spesifikasi material.
* **AI-Driven Triage (Gemini 2.5 Flash):** Diagnosis otonom berbasis NLP dengan *Confidence Gate* algoritmik untuk menentukan kelayakan perbaikan mandiri (DIY) vs eskalasi ke teknisi ahli.
* **Geospatial Matching Engine:** Kueri radius presisi menggunakan ekstensi PostGIS (`ST_DWithin`, `ST_MakePoint`) yang mengeliminasi beban komputasi Haversine di level aplikasi.
* **Technician Control Center:** Sistem kendali ketersediaan *real-time* (Online/Offline) yang tersinkronisasi langsung dengan filter PostGIS, dilengkapi agregasi metrik performa dan pendapatan teknisi.
* **State-Machine Order Lifecycle & Anti-Fraud:** Transisi status pesanan transaksional yang ketat (*Pending, Accepted, Cancelled, Completed*) dengan kewajiban verifikasi geolokasi dan bukti foto pada titik akhir penyelesaian.
* **Dynamic Marketplace Revenue Model:** Implementasi model bisnis transparan dengan sistem *Take Rate* 10% sebagai komisi platform. Sistem secara otomatis menghitung *Total Fee*, *Platform Fee*, dan *Net Technician Earnings*.
* **ACID Transactional Integrity & Data Retention:** Operasi relasional kompleks dibungkus dalam *Database Transactions* murni untuk menjamin konsistensi data, didukung oleh mekanisme *Soft Delete* tersinkronisasi untuk retensi dan pemulihan akun tanpa anomali.
* **Industrial Security Protocol:** Implementasi keamanan berlapis mencakup *Stateless Authorization* via JWT dengan perlindungan HTTPOnly & Secure Cookies, *Centralized Error Masking* untuk mencegah kebocoran skema database, serta *Rate Limiting* tingkat endpoint untuk mencegah serangan brute-force dan spamming.

## Tech Stack

| Component            | Technology                      |
| :------------------- | :------------------------------ |
| **Language**         | Go 1.25+                        |
| **Web Framework**    | Fiber v2                        |
| **Database**         | PostgreSQL (Supabase) + PostGIS |
| **ORM**              | GORM v1.31                      |
| **Containerization** | Docker & Docker Compose         |
| **Documentation**    | Swaggo (OpenAPI / Swagger 2.0)  |

## API Documentation

Referensi interaktif dan kontrak data tersedia melalui Swagger UI. [Klik di sini untuk membuka dokumentasi API](https://ecoserve-api.onrender.com/swagger/index.html).

## Local Environment Setup

### Option 1: Docker (Recommended)
Gunakan metode ini untuk lingkungan yang terisolasi dan identik dengan server produksi (Windows / Linux / macOS). Pastikan Docker dan Docker Compose sudah terinstal.

```bash
docker-compose up --build
```
Layanan akan tersedia di `http://localhost:3000`.

### Option 2: Manual Installation

**1. Clone repository**
```bash
git clone https://github.com/noireveil/ecoserve-backend.git
```

**2. Konfigurasi environment**
```bash
cp .env.example .env
```

**3. Install dependensi**
```bash
go mod tidy
```

**4. Build dokumentasi Swagger**
```bash
swag init -g cmd/api/main.go --parseDependency --parseInternal
```

**5. Jalankan development server**
```bash
go run cmd/api/main.go
```

## Environmental Impact Analysis

Sistem ini terintegrasi langsung dengan standar pengukur emisi industri. Setiap perbaikan yang tervalidasi secara otomatis dikonversi menjadi metrik penyelamatan CO2e menggunakan algoritma dari *Environmental Protection Agency* (EPA) Waste Reduction Model (WARM) versi 15.

---
*Dikembangkan secara eksklusif untuk mendorong transisi ke ekonomi sirkular berbasis teknologi.*
