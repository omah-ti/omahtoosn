# TO OSN Backend

Panduan cepat untuk tim frontend menjalankan backend lokal.

## Cara Paling Cepat

Jalankan dari folder `backend`:

```sh
go run ./cmd/dev_setup --run
```

Command ini otomatis:

- membuat `.env` dari `.env.example` jika belum ada
- menjalankan PostgreSQL via Docker Compose
- menunggu database siap
- menjalankan migration
- mengisi seed demo
- menjalankan API di `http://localhost:8081`

Command ini hanya mau memakai database localhost. Ini sengaja agar migration/seed tidak tidak sengaja menyentuh database remote.

Swagger tersedia di:

```txt
http://localhost:8081/swagger/index.html
```

Jika tidak memakai Docker dan sudah punya PostgreSQL lokal:

```sh
go run ./cmd/dev_setup --skip-docker --run
```

Pastikan `DATABASE_URL` di `.env` sesuai database lokalmu.

## Prasyarat

Minimal:

- Go 1.25+
- Docker Desktop

Opsional:

- PostgreSQL lokal, jika tidak pakai Docker
- Postman, untuk import collection
- `make`, `psql`, dan `migrate` CLI jika ingin menjalankan langkah manual

Semua OS bisa memakai command `go run` yang sama. Di Windows, gunakan PowerShell atau terminal bawaan IDE.

## Setup Manual

Pakai bagian ini hanya jika command cepat bermasalah.

### 1. Env

macOS/Linux:

```sh
cp .env.example .env
```

Windows PowerShell:

```powershell
Copy-Item .env.example .env
```

Default lokal yang penting:

```env
APP_PORT=8081
DATABASE_URL=postgres://postgres:postgres@localhost:5433/to_osn?sslmode=disable
FRONTEND_URL=http://localhost:3000
PASSWORD_RESET_PATH=/reset-password
```

Untuk fitur lupa password:

```env
RESEND_API_KEY=your_resend_api_key
EMAIL_FROM="TO OSN <onboarding@resend.dev>"
```

Produksi wajib memakai domain pengirim yang sudah verified di Resend.

### 2. Database

```sh
docker compose up -d
```

Default database Docker:

```txt
host: localhost
port: 5433
database: to_osn
user: postgres
password: postgres
```

### 3. Migration

Dengan `make`:

```sh
make migrate-up
```

Tanpa `make`:

```sh
migrate -path ./migrations -database "postgres://postgres:postgres@localhost:5433/to_osn?sslmode=disable" up
```

### 4. Seed

Seed demo:

```sh
make seed-demo
```

Tanpa `make`:

```sh
psql "postgres://postgres:postgres@localhost:5433/to_osn?sslmode=disable" -f ./seeds/demo_seed.sql
```

Seed soal OmahTOOSN:

```sh
go run ./cmd/seed_questions
```

Seed soal langsung aktif:

macOS/Linux:

```sh
TRYOUT_STATUS=ongoing go run ./cmd/seed_questions
```

Windows PowerShell:

```powershell
$env:TRYOUT_STATUS="ongoing"; go run ./cmd/seed_questions
```

### 5. Run API

```sh
go run ./cmd/api
```

Cek:

```sh
curl http://localhost:8081/health
```

## Postman

Import dua file ini:

- `postman/to-osn-v1.postman_collection.json`
- `postman/to-osn-local.postman_environment.json`

## Reset Password Lokal

1. Register user.
2. Panggil `POST /api/v1/auth/forgot-password`.
3. Cek email user.
4. Ambil token dari link `?token=...`.
5. Panggil `POST /api/v1/auth/reset-password`.

Selama halaman frontend reset password belum ada, token dari email bisa dipakai langsung di Postman.

## Command Cepat Lain

Setup tanpa menjalankan API:

```sh
go run ./cmd/dev_setup
```

Setup dengan seed soal OmahTOOSN:

```sh
go run ./cmd/dev_setup --seed=omahtoosn --run
```

Setup tanpa seed:

```sh
go run ./cmd/dev_setup --seed=none --run
```

Izinkan DB non-lokal hanya jika benar-benar sengaja:

```sh
go run ./cmd/dev_setup --allow-nonlocal-db --run
```

## Troubleshooting

- Docker tidak jalan: buka Docker Desktop, lalu ulangi command.
- Tidak pakai Docker: set `DATABASE_URL`, lalu pakai `--skip-docker`.
- Port database bentrok: ubah port di `docker-compose.yml` dan `DATABASE_URL`.
- Email tidak terkirim: cek `RESEND_API_KEY`, `EMAIL_FROM`, dan verified domain Resend.
- Cookie localhost: gunakan `COOKIE_SECURE=false` dan `COOKIE_SAME_SITE=Lax`.
- Go cache `Access is denied`: set cache ke folder repo, contoh PowerShell: `$env:GOCACHE="$PWD\bin\go-cache"; go run ./cmd/dev_setup --run`.
