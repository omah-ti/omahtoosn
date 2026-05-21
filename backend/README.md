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

## Deploy Backend ke Railway

Backend ini siap dideploy sebagai service terpisah dari folder `backend`. File `railway.json` mengatur Docker build, pre-deploy migration, start command, healthcheck, dan restart policy. File `Dockerfile` membangun tiga binary:

- `api` untuk service utama
- `migrate` untuk menjalankan migration SQL sebelum deploy aktif
- `seed_questions` untuk seed soal OmahTOOSN sekali saja jika diperlukan

### 1. Buat service

Di Railway:

1. Buat project dari repository GitHub monorepo ini.
2. Tambahkan PostgreSQL service: `+ New` -> `Database` -> `PostgreSQL`.
3. Tambahkan backend service dari repository yang sama.
4. Pada backend service, buka `Settings` lalu set:

```txt
Root Directory: /backend
Railway Config File: /backend/railway.json
Watch Paths: /backend/**
```

`Root Directory` wajib supaya Railway tidak membangun Next.js frontend di root repo. `Railway Config File` wajib karena config file Railway tidak otomatis mengikuti root directory.

5. Pada tab `Networking`, buat public domain untuk backend agar endpoint `/health` dan `/swagger/index.html` bisa diakses dari internet.

### 2. Set variable backend

Di backend service, buka `Variables` lalu isi:

```env
APP_ENV=production
APP_NAME=to-osn-backend
APP_VERSION=1.0.0
DATABASE_URL=${{Postgres.DATABASE_URL}}
JWT_SECRET=isi_dengan_secret_panjang_random
ACCESS_TOKEN_TTL_MINUTES=15
REFRESH_TOKEN_TTL_HOURS=168
CORS_ALLOW_ORIGINS=https://domain-frontend-kamu
FRONTEND_URL=https://domain-frontend-kamu
PASSWORD_RESET_PATH=/reset-password
PASSWORD_RESET_TTL_MINUTES=30
COOKIE_SECURE=true
COOKIE_SAME_SITE=None
COOKIE_DOMAIN=
RESEND_API_KEY=isi_jika_fitur_lupa_password_dipakai
EMAIL_FROM=TO OSN <noreply@domain-terverifikasi-kamu>
EMAIL_REPLY_TO=
```

Catatan variable:

- `APP_PORT` tidak perlu diisi. Railway menyediakan `PORT`, dan backend sudah membaca `PORT` otomatis.
- `DATABASE_URL` harus memakai reference variable dari PostgreSQL service, biasanya `${{Postgres.DATABASE_URL}}`. Jika nama service PostgreSQL berbeda, sesuaikan `Postgres`.
- `JWT_SECRET` jangan pakai default. Gunakan string random panjang.
- Jika frontend dan backend berada di domain berbeda, gunakan `COOKIE_SECURE=true` dan `COOKIE_SAME_SITE=None`.
- Jika memakai custom domain satu parent domain, misalnya `app.example.com` dan `api.example.com`, `COOKIE_DOMAIN` boleh diisi `.example.com`. Kalau memakai domain Railway/Vercel terpisah, biarkan kosong.
- `RESEND_API_KEY` dan `EMAIL_FROM` wajib hanya untuk fitur forgot password. Domain pengirim harus sudah verified di Resend untuk production.

### 3. Deploy dan cek

Setelah variable disimpan, deploy backend. Pada deploy pertama:

1. Docker image dibangun dari `backend/Dockerfile`.
2. Railway menjalankan pre-deploy command `./migrate`.
3. Jika migration sukses, Railway menjalankan `./api`.
4. Healthcheck memanggil `/health`.

Endpoint yang perlu dicek:

```txt
https://domain-backend-kamu/health
https://domain-backend-kamu/swagger/index.html
```

### 4. Seed soal OmahTOOSN

Migration hanya membuat schema. Untuk mengisi soal, jalankan sekali dari shell/one-off command backend service setelah deploy sukses:

```sh
TRYOUT_STATUS=ongoing ./seed_questions
```

Jika belum mau tryout aktif, jalankan tanpa status ongoing:

```sh
./seed_questions
```

Jangan jadikan seed sebagai pre-deploy command permanen, karena setelah user mulai mengerjakan tryout, seed ulang bisa ditolak oleh proteksi `ALLOW_SEED_WITH_ATTEMPTS`.

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
