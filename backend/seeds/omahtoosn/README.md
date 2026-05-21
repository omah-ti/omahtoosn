# OmahTOOSN question seed

`soal.json` is the source data for the real tryout questions. The PNG files referenced by the question HTML live in `../../assets/questions` and are served by the backend at `/question-assets`.

Run from `backend`:

```sh
go run ./cmd/seed_questions
```

The seed is idempotent and defaults to `TRYOUT_STATUS=draft` so it does not conflict with another ongoing tryout. To make this tryout active, archive the current ongoing tryout first, then run:

```sh
TRYOUT_STATUS=ongoing go run ./cmd/seed_questions
```

Useful overrides:

- `DATABASE_URL`
- `TRYOUT_SLUG`
- `TRYOUT_TITLE`
- `TRYOUT_DURATION_MINUTES`
- `QUESTION_ASSET_BASE_URL`
