# EcoServe Backend Engine

EcoServe adalah infrastruktur perangkat lunak *Enterprise-Ready* yang dirancang untuk mendigitalkan ekonomi sirkular dan memperpanjang usia perangkat elektronik. API ini menjembatani konsumen dan teknisi perbaikan melalui lapisan kecerdasan buatan, pemetaan geospasial presisi, dan pelacakan jejak karbon (EPA WARM v15).

## 🚀 Core Architecture

Sistem ini dibangun di atas fondasi **Domain-Driven Design** menggunakan bahasa **Go (Golang)**, memastikan skalabilitas tinggi, kemudahan pengujian, dan keterisolasian logika bisnis.

* **Digital Product Passport (DPP):** Manajemen siklus hidup elektronik (*Electronic Lifecycle Management*) dengan rekam jejak kepemilikan dan spesifikasi material.
* **AI-Driven Triage (Gemini 2.5 Flash):** Diagnosis otonom berbasis NLP dengan *Confidence Gate* algoritmik untuk menentukan kelayakan perbaikan mandiri (DIY) vs eskalasi ke teknisi ahli.
* **Geospatial Matching Engine:** Kueri radius presisi menggunakan ekstensi **PostGIS** (`ST_DWithin`, `ST_MakePoint`) yang mengeleminasi beban komputasi Haversine di level aplikasi.
* **ACID Transactional Integrity:** Operasi relasional kompleks (seperti *Order Completion*, kalkulasi *Impact Tracker*, dan agregasi *Review*) dibungkus dalam *Database Transactions* murni untuk menjamin konsistensi data (Zero Data Corruption).
* **Industrial Security Protocol:**
  * Implementasi otorisasi *Stateless* menggunakan **JWT (JSON Web Tokens)** dengan perlindungan **HTTPOnly & Secure Cookies** untuk memitigasi serangan XSS.
  * Proteksi manipulasi data dengan **Centralized Error Masking** untuk mencegah kebocoran skema *database*.
  * **Rate Limiting** tingkat *endpoint* untuk mencegah serangan *brute-force* dan *spamming*.
* **Zero-Downtime Reliability:** Dilengkapi dengan rutinitas *Graceful Shutdown* untuk memastikan semua proses yang berjalan diselesaikan sebelum siklus hidup kontainer diakhiri.

## 🛠 Tech Stack

* **Language:** Go 1.25+
* **Web Framework:** Fiber v2 (Express-like, dialihkan untuk performa tinggi)
* **Database:** PostgreSQL terskala awan (Supabase) + PostGIS
* **ORM:** GORM v1.31
* **AI & NLP:** Google Gemini 2.5 Flash REST API
* **Communication:** Mailjet API (Transactional Email)
* **Documentation:** Swaggo (OpenAPI/Swagger 2.0)

## 📖 API Documentation

Referensi interaktif dan kontrak data tersedia secara publik melalui antarmuka Swagger UI:
👉 **[EcoServe API Documentation](https://ecoserve-api.onrender.com/swagger/index.html)**

## ⚙️ Panduan Inisialisasi Lingkungan Lokal

Infrastruktur pendukung disediakan menggunakan Docker untuk standarisasi pengembangan.

1. **Persiapan Repositori**
   ```bash
   git clone https://github.com/noireveil/ecoserve-backend.git
   cd ecoserve-backend
   ```
2. **Konfigurasi Environment**
   Salin fail referensi lingkungan dan sesuaikan variabel kunci.
   ```bash
   cp .env.example .env
   ```
3. **Instalasi Dependensi**
   ```bash
   go mod tidy
   ```
4. **Membangun Dokumentasi Swagger**
   ```bash
   swag init -g cmd/api/main.go --parseDependency --parseInternal
   ```
5. **Menjalankan Engine (Development)**
   ```bash
   go run cmd/api/main.go
   ```

## 🌍 Analisis Dampak Lingkungan

Sistem ini terintegrasi langsung dengan standar pengukur emisi. Setiap perbaikan yang tervalidasi secara otomatis dikonversi menjadi metrik penyelamatan CO2e menggunakan algoritma dari *Environmental Protection Agency* (EPA) Waste Reduction Model (WARM) versi 15.

---
*Dikembangkan secara eksklusif untuk mendorong transisi ke ekonomi sirkular berbasis teknologi.*