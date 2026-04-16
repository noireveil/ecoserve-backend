# EcoServe Backend

EcoServe adalah platform manajemen siklus hidup elektronik yang dirancang untuk memperpanjang usia perangkat dan menekan pertumbuhan limbah elektronik (*e-waste*) dalam ekosistem ekonomi sirkular. API ini dibangun menggunakan **Go** dengan fokus pada efisiensi performa, skalabilitas, dan akurasi data geospasial.

## Fitur Utama

* **Digital Product Passport (DPP):** Transparansi siklus hidup perangkat melalui pencatatan riwayat digital yang komprehensif.
* **AI-Driven Triage System:** Diagnosis kerusakan otomatis berbasis **Gemini 2.5 Flash** dengan integrasi *Confidence Gate* untuk menjamin akurasi hasil analisis.
* **Geospatial Service Engine:** Optimasi pencarian teknisi terdekat memanfaatkan **PostGIS** dan sistem koordinat presisi untuk kebutuhan navigasi.
* **Advanced Order Management:** Manajemen alur kerja layanan perbaikan yang terintegrasi, mencakup fitur *real-time tracking* koordinat konsumen dan sistem *Accept Order* untuk teknisi.
* **Automated Impact Tracking:** Kalkulasi metrik lingkungan otomatis berdasarkan algoritma **EPA WARM v15** untuk mengukur penghematan limbah elektronik.
* **Industrial Security Layer:** Perlindungan data melalui autentikasi **JWT**, enkripsi sesi, dan mitigasi *fraud* menggunakan verifikasi geospasial (*GPS Locking*).

## Tumpukan Teknologi

* **Core:** Go (Golang) 1.25+ & Fiber v2 Framework
* **Database:** PostgreSQL with PostGIS Extension (Hosted on Supabase)
* **Artificial Intelligence:** Google Gemini 2.5 Flash API
* **Authentication:** JWT-based Secure Authorization
* **Documentation:** OpenAPI (Swagger) Standard

## Arsitektur Sistem

Proyek ini mengimplementasikan *Layered Architecture* untuk memisahkan logika bisnis dari lapisan infrastruktur, memastikan kode mudah diuji dan dikembangkan:

* `cmd/api/`: *Entry point* aplikasi dan konfigurasi *middleware*.
* `internal/domain/`: Definisi entitas dan skema data inti.
* `internal/usecase/`: Implementasi logika bisnis dan aturan sistem.
* `internal/repository/`: Abstraksi akses basis data dan kueri geospasial.
* `internal/delivery/`: Lapisan komunikasi HTTP dan API *handlers*.

## Dokumentasi API

Dokumentasi API tersedia secara interaktif dan dapat diakses melalui Swagger UI:
`https://ecoserve-api.onrender.com/swagger/index.html`

## Panduan Instalasi Lokal

1.  **Persiapan Repositori:**
    ```bash
    git clone https://github.com/noireveil/ecoserve-backend.git
    cd ecoserve-backend
    ```
2.  **Konfigurasi Environment:**
    Lengkapi fail `.env` pada direktori *root* dengan kredensial yang diperlukan (Database, JWT Secret, dan Gemini API Key).
3.  **Manajemen Dependensi:**
    ```bash
    go mod tidy
    ```
4.  **Menjalankan Aplikasi:**
    ```bash
    go run cmd/api/main.go
    ```
5.  **Verifikasi Sistem:**
    Pastikan layanan aktif dengan mengakses *endpoint* kesehatan sistem di `http://localhost:3000/health`.