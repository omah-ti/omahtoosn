package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/omah-ti/omahtoosn/backend/internal/platform/config"
	"github.com/omah-ti/omahtoosn/backend/internal/platform/db"
)

const (
	defaultQuestionPath = "./seeds/omahtoosn/soal.json"
	defaultAssetDir     = "./assets/questions"
	defaultTryoutSlug   = "omahtoosn-informatika-2026"
	defaultTryoutTitle  = "Try Out OSN-K Informatika"
)

var imageSrcPattern = regexp.MustCompile(`src="([^"]+\.png)"`)

type seedFile struct {
	Questions []seedQuestion `json:"questions"`
}

type seedQuestion struct {
	Code                string         `json:"code"`
	QuestionType        string         `json:"question_type"`
	PromptHTML          string         `json:"prompt_html"`
	DisplayOrder        int            `json:"display_order"`
	Points              float64        `json:"points"`
	Options             []seedOption   `json:"options"`
	AnswerKey           *seedAnswerKey `json:"answer_key"`
	ShortAnswerVariants []string       `json:"short_answer_variants"`
}

type seedOption struct {
	Key          string `json:"key"`
	Text         string `json:"text"`
	DisplayOrder int    `json:"display_order"`
}

type seedAnswerKey struct {
	CorrectOptionKey string `json:"correct_option_key"`
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	questionPath := envString("QUESTION_SEED_FILE", defaultQuestionPath)
	assetDir := envString("QUESTION_ASSET_DIR", defaultAssetDir)
	assetBaseURL := strings.TrimRight(envString("QUESTION_ASSET_BASE_URL", "/question-assets/"), "/") + "/"

	payload, err := loadSeedFile(questionPath)
	if err != nil {
		log.Fatal(err)
	}
	if err := validateSeed(payload, assetDir); err != nil {
		log.Fatal(err)
	}
	rewriteImageRefs(payload, assetBaseURL)

	if envBool("SEED_VALIDATE_ONLY", false) {
		log.Printf("validated %d questions from %s", len(payload.Questions), questionPath)
		return
	}

	cfg := config.Load()
	pool, err := db.New(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("db open failed: %v", err)
	}
	defer pool.Close()

	if err := seedTryout(ctx, pool, payload); err != nil {
		log.Fatal(err)
	}
	log.Printf("seeded %d questions into tryout slug=%s", len(payload.Questions), envString("TRYOUT_SLUG", defaultTryoutSlug))
}

func loadSeedFile(path string) (*seedFile, error) {
	content, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		return nil, fmt.Errorf("read question seed file: %w", err)
	}
	var payload seedFile
	if err := json.Unmarshal(content, &payload); err != nil {
		return nil, fmt.Errorf("parse question seed file: %w", err)
	}
	return &payload, nil
}

func validateSeed(payload *seedFile, assetDir string) error {
	if len(payload.Questions) == 0 {
		return errors.New("question seed file has no questions")
	}

	codes := map[string]bool{}
	orders := map[int]bool{}
	imageRefs := map[string]bool{}
	for _, question := range payload.Questions {
		code := strings.TrimSpace(question.Code)
		if code == "" {
			return errors.New("question code is required")
		}
		if codes[code] {
			return fmt.Errorf("duplicate question code: %s", code)
		}
		codes[code] = true

		if question.DisplayOrder <= 0 {
			return fmt.Errorf("question %s has invalid display_order", code)
		}
		if orders[question.DisplayOrder] {
			return fmt.Errorf("duplicate display_order: %d", question.DisplayOrder)
		}
		orders[question.DisplayOrder] = true

		if strings.TrimSpace(question.PromptHTML) == "" {
			return fmt.Errorf("question %s prompt_html is required", code)
		}
		if question.Points < 0 {
			return fmt.Errorf("question %s has negative points", code)
		}

		collectImageRefs(question.PromptHTML, imageRefs)
		switch question.QuestionType {
		case "multiple_choice":
			if len(question.Options) == 0 {
				return fmt.Errorf("multiple_choice question %s has no options", code)
			}
			if question.AnswerKey == nil || strings.TrimSpace(question.AnswerKey.CorrectOptionKey) == "" {
				return fmt.Errorf("multiple_choice question %s has no answer_key", code)
			}
			optionKeys := map[string]bool{}
			optionOrders := map[int]bool{}
			for _, option := range question.Options {
				key := strings.ToUpper(strings.TrimSpace(option.Key))
				if key == "" {
					return fmt.Errorf("question %s has option with empty key", code)
				}
				if optionKeys[key] {
					return fmt.Errorf("question %s has duplicate option key %s", code, key)
				}
				optionKeys[key] = true
				if option.DisplayOrder <= 0 {
					return fmt.Errorf("question %s option %s has invalid display_order", code, key)
				}
				if optionOrders[option.DisplayOrder] {
					return fmt.Errorf("question %s has duplicate option display_order %d", code, option.DisplayOrder)
				}
				optionOrders[option.DisplayOrder] = true
				if strings.TrimSpace(option.Text) == "" {
					return fmt.Errorf("question %s option %s text is required", code, key)
				}
				collectImageRefs(option.Text, imageRefs)
			}
			answerKey := strings.ToUpper(strings.TrimSpace(question.AnswerKey.CorrectOptionKey))
			if !optionKeys[answerKey] {
				return fmt.Errorf("question %s answer_key %s is not a valid option", code, answerKey)
			}
		case "short_text":
			if len(question.ShortAnswerVariants) == 0 {
				return fmt.Errorf("short_text question %s has no short_answer_variants", code)
			}
		default:
			return fmt.Errorf("question %s has unsupported question_type %s", code, question.QuestionType)
		}
	}

	missing := []string{}
	for ref := range imageRefs {
		if _, err := os.Stat(filepath.Join(assetDir, ref)); err != nil {
			if os.IsNotExist(err) {
				missing = append(missing, ref)
				continue
			}
			return fmt.Errorf("check asset %s: %w", ref, err)
		}
	}
	if len(missing) > 0 {
		sort.Strings(missing)
		return fmt.Errorf("missing question image assets: %s", strings.Join(missing, ", "))
	}
	return nil
}

func collectImageRefs(input string, refs map[string]bool) {
	for _, match := range imageSrcPattern.FindAllStringSubmatch(input, -1) {
		if len(match) == 2 && !strings.Contains(match[1], "/") {
			refs[match[1]] = true
		}
	}
}

func rewriteImageRefs(payload *seedFile, assetBaseURL string) {
	for i := range payload.Questions {
		payload.Questions[i].PromptHTML = rewriteImageSrc(payload.Questions[i].PromptHTML, assetBaseURL)
		for j := range payload.Questions[i].Options {
			payload.Questions[i].Options[j].Text = rewriteImageSrc(payload.Questions[i].Options[j].Text, assetBaseURL)
		}
	}
}

func rewriteImageSrc(input, assetBaseURL string) string {
	return imageSrcPattern.ReplaceAllStringFunc(input, func(match string) string {
		parts := imageSrcPattern.FindStringSubmatch(match)
		if len(parts) != 2 || strings.Contains(parts[1], "/") {
			return match
		}
		return `src="` + assetBaseURL + parts[1] + `"`
	})
}

func seedTryout(ctx context.Context, pool *pgxpool.Pool, payload *seedFile) error {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	slug := envString("TRYOUT_SLUG", defaultTryoutSlug)
	status := envString("TRYOUT_STATUS", "draft")
	tryoutID, existed, err := upsertTryout(ctx, tx, slug, status)
	if err != nil {
		return err
	}
	if existed && !envBool("ALLOW_SEED_WITH_ATTEMPTS", false) {
		hasAttempts, err := tryoutHasAttempts(ctx, tx, tryoutID)
		if err != nil {
			return err
		}
		if hasAttempts {
			return fmt.Errorf("tryout %s already has attempts; set ALLOW_SEED_WITH_ATTEMPTS=true only if you intentionally want to update its questions", slug)
		}
	}

	for _, question := range payload.Questions {
		questionID, err := upsertQuestion(ctx, tx, tryoutID, question)
		if err != nil {
			return err
		}
		if err := syncQuestionAnswers(ctx, tx, questionID, question); err != nil {
			return err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}
	return nil
}

func upsertTryout(ctx context.Context, tx pgx.Tx, slug, status string) (string, bool, error) {
	existingID, err := getTryoutID(ctx, tx, slug)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return "", false, fmt.Errorf("check tryout: %w", err)
	}
	existed := err == nil

	query := `
		INSERT INTO tryouts (slug, title, description, instructions, status, duration_minutes, show_leaderboard)
		VALUES ($1, $2, $3, $4, $5, $6, TRUE)
		ON CONFLICT (slug) DO UPDATE SET
			title = EXCLUDED.title,
			description = EXCLUDED.description,
			instructions = EXCLUDED.instructions,
			status = EXCLUDED.status,
			duration_minutes = EXCLUDED.duration_minutes,
			show_leaderboard = EXCLUDED.show_leaderboard
		RETURNING id
	`
	var id string
	err = tx.QueryRow(ctx, query,
		slug,
		envString("TRYOUT_TITLE", defaultTryoutTitle),
		envString("TRYOUT_DESCRIPTION", "Seed soal OmahTI TO OSN Informatika."),
		envString("TRYOUT_INSTRUCTIONS", "Baca setiap soal dengan teliti. Jawaban tersimpan di server. Submit sebelum waktu habis."),
		status,
		envInt("TRYOUT_DURATION_MINUTES", 150),
	).Scan(&id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" && status == "ongoing" {
			return "", existed, errors.New("another ongoing tryout already exists; seed with TRYOUT_STATUS=draft or archive the existing ongoing tryout first")
		}
		return "", existed, fmt.Errorf("upsert tryout: %w", err)
	}
	if existingID != "" && existingID != id {
		return "", existed, fmt.Errorf("tryout id changed unexpectedly for slug %s", slug)
	}
	return id, existed, nil
}

func getTryoutID(ctx context.Context, tx pgx.Tx, slug string) (string, error) {
	var id string
	err := tx.QueryRow(ctx, `SELECT id FROM tryouts WHERE slug = $1`, slug).Scan(&id)
	return id, err
}

func tryoutHasAttempts(ctx context.Context, tx pgx.Tx, tryoutID string) (bool, error) {
	var exists bool
	if err := tx.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM attempts WHERE tryout_id = $1)`, tryoutID).Scan(&exists); err != nil {
		return false, fmt.Errorf("check attempts: %w", err)
	}
	return exists, nil
}

func upsertQuestion(ctx context.Context, tx pgx.Tx, tryoutID string, question seedQuestion) (string, error) {
	query := `
		INSERT INTO questions (tryout_id, code, question_type, prompt_html, display_order, points)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (tryout_id, code) DO UPDATE SET
			question_type = EXCLUDED.question_type,
			prompt_html = EXCLUDED.prompt_html,
			display_order = EXCLUDED.display_order,
			points = EXCLUDED.points
		RETURNING id
	`
	var id string
	if err := tx.QueryRow(ctx, query,
		tryoutID,
		strings.TrimSpace(question.Code),
		question.QuestionType,
		question.PromptHTML,
		question.DisplayOrder,
		question.Points,
	).Scan(&id); err != nil {
		return "", fmt.Errorf("upsert question %s: %w", question.Code, err)
	}
	return id, nil
}

func syncQuestionAnswers(ctx context.Context, tx pgx.Tx, questionID string, question seedQuestion) error {
	switch question.QuestionType {
	case "multiple_choice":
		keys := make([]string, 0, len(question.Options))
		for _, option := range question.Options {
			key := strings.ToUpper(strings.TrimSpace(option.Key))
			keys = append(keys, key)
			if _, err := tx.Exec(ctx, `
				INSERT INTO question_options (question_id, option_key, option_text, display_order)
				VALUES ($1, $2, $3, $4)
				ON CONFLICT (question_id, option_key) DO UPDATE SET
					option_text = EXCLUDED.option_text,
					display_order = EXCLUDED.display_order
			`, questionID, key, option.Text, option.DisplayOrder); err != nil {
				return fmt.Errorf("upsert option %s/%s: %w", question.Code, key, err)
			}
		}
		if _, err := tx.Exec(ctx, `DELETE FROM question_options WHERE question_id = $1 AND NOT (option_key = ANY($2))`, questionID, keys); err != nil {
			return fmt.Errorf("delete stale options for %s: %w", question.Code, err)
		}
		answerKey := strings.ToUpper(strings.TrimSpace(question.AnswerKey.CorrectOptionKey))
		if _, err := tx.Exec(ctx, `
			INSERT INTO question_answer_keys (question_id, correct_option_key)
			VALUES ($1, $2)
			ON CONFLICT (question_id) DO UPDATE SET correct_option_key = EXCLUDED.correct_option_key
		`, questionID, answerKey); err != nil {
			return fmt.Errorf("upsert answer key for %s: %w", question.Code, err)
		}
		if _, err := tx.Exec(ctx, `DELETE FROM question_short_answer_variants WHERE question_id = $1`, questionID); err != nil {
			return fmt.Errorf("delete short answers for multiple choice %s: %w", question.Code, err)
		}
	case "short_text":
		if _, err := tx.Exec(ctx, `DELETE FROM question_options WHERE question_id = $1`, questionID); err != nil {
			return fmt.Errorf("delete options for short answer %s: %w", question.Code, err)
		}
		if _, err := tx.Exec(ctx, `DELETE FROM question_answer_keys WHERE question_id = $1`, questionID); err != nil {
			return fmt.Errorf("delete answer key for short answer %s: %w", question.Code, err)
		}
		normalizedValues := make([]string, 0, len(question.ShortAnswerVariants))
		for _, variant := range question.ShortAnswerVariants {
			answerText := strings.TrimSpace(html.UnescapeString(variant))
			normalized := normalizeAnswer(answerText)
			normalizedValues = append(normalizedValues, normalized)
			if _, err := tx.Exec(ctx, `
				INSERT INTO question_short_answer_variants (question_id, answer_text, normalized_text)
				VALUES ($1, $2, $3)
				ON CONFLICT (question_id, normalized_text) DO UPDATE SET answer_text = EXCLUDED.answer_text
			`, questionID, answerText, normalized); err != nil {
				return fmt.Errorf("upsert short answer for %s: %w", question.Code, err)
			}
		}
		if _, err := tx.Exec(ctx, `DELETE FROM question_short_answer_variants WHERE question_id = $1 AND NOT (normalized_text = ANY($2))`, questionID, normalizedValues); err != nil {
			return fmt.Errorf("delete stale short answers for %s: %w", question.Code, err)
		}
	}
	return nil
}

func normalizeAnswer(input string) string {
	parts := strings.Fields(strings.ToLower(strings.TrimSpace(input)))
	return strings.Join(parts, " ")
}

func envString(key, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}
	return fallback
}

func envInt(key string, fallback int) int {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	var parsed int
	if _, err := fmt.Sscanf(value, "%d", &parsed); err != nil {
		return fallback
	}
	return parsed
}

func envBool(key string, fallback bool) bool {
	switch strings.ToLower(strings.TrimSpace(os.Getenv(key))) {
	case "true", "1", "yes", "y":
		return true
	case "false", "0", "no", "n":
		return false
	default:
		return fallback
	}
}
