APP_NAME=to-osn-backend
DB_URL?=postgres://postgres:postgres@localhost:5433/to_osn?sslmode=disable
MIGRATIONS_PATH?=./migrations

run:
	go run ./cmd/api

migrate-up:
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" up

migrate-down:
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" down 1

migrate-force:
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" force $(VERSION)

seed-demo:
	psql "$(DB_URL)" -f ./seeds/demo_seed.sql

seed-omahtoosn:
	DATABASE_URL="$(DB_URL)" go run ./cmd/seed_questions
