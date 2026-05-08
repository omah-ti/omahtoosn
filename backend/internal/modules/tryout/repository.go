package tryout

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/omah-ti/omahtoosn/backend/internal/platform/store"
)

type Repository struct{}

func NewRepository() *Repository {
	return &Repository{}
}

func (r *Repository) Begin(ctx context.Context, pool *pgxpool.Pool) (pgx.Tx, error) {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	return tx, nil
}

func (r *Repository) GetCurrentTryout(ctx context.Context, db store.DBTX) (Tryout, error) {
	query := `
		SELECT
			t.id,
			t.slug,
			t.title,
			COALESCE(t.description, ''),
			COALESCE(t.instructions, ''),
			t.status,
			t.duration_minutes,
			t.starts_at,
			t.ends_at,
			t.results_published_at,
			t.show_leaderboard,
			COUNT(q.id)::int AS question_count
		FROM tryouts t
		LEFT JOIN questions q ON q.tryout_id = t.id
		WHERE t.status = 'ongoing'
		GROUP BY t.id
		LIMIT 1
	`
	var out Tryout
	var startsAt sql.NullTime
	var endsAt sql.NullTime
	var resultsPublishedAt sql.NullTime
	if err := db.QueryRow(ctx, query).Scan(
		&out.ID,
		&out.Slug,
		&out.Title,
		&out.Description,
		&out.Instructions,
		&out.Status,
		&out.DurationMinutes,
		&startsAt,
		&endsAt,
		&resultsPublishedAt,
		&out.ShowLeaderboard,
		&out.QuestionCount,
	); err != nil {
		return Tryout{}, err
	}
	if startsAt.Valid {
		out.StartsAt = &startsAt.Time
	}
	if endsAt.Valid {
		out.EndsAt = &endsAt.Time
	}
	if resultsPublishedAt.Valid {
		out.ResultsPublishedAt = &resultsPublishedAt.Time
	}
	return out, nil
}

func (r *Repository) GetAttemptByUserTryout(ctx context.Context, db store.DBTX, userID, tryoutID string) (Attempt, error) {
	query := `
		SELECT id, user_id, tryout_id, status, started_at, expires_at, submitted_at, last_synced_at, version, total_questions, answered_questions, flagged_questions, correct_count, wrong_count, unanswered_count, final_score, created_at, updated_at
		FROM attempts
		WHERE user_id = $1 AND tryout_id = $2
		LIMIT 1
	`
	return scanAttempt(db.QueryRow(ctx, query, userID, tryoutID))
}

func (r *Repository) LockAttemptByUserTryout(ctx context.Context, db store.DBTX, userID, tryoutID string) (Attempt, error) {
	query := `
		SELECT id, user_id, tryout_id, status, started_at, expires_at, submitted_at, last_synced_at, version, total_questions, answered_questions, flagged_questions, correct_count, wrong_count, unanswered_count, final_score, created_at, updated_at
		FROM attempts
		WHERE user_id = $1 AND tryout_id = $2
		LIMIT 1
		FOR UPDATE
	`
	return scanAttempt(db.QueryRow(ctx, query, userID, tryoutID))
}

func (r *Repository) CreateAttempt(ctx context.Context, db store.DBTX, userID, tryoutID string, expiresAt time.Time, totalQuestions int) (Attempt, error) {
	query := `
		INSERT INTO attempts (user_id, tryout_id, status, expires_at, total_questions, unanswered_count)
		VALUES ($1, $2, 'ongoing', $3, $4, $4)
		RETURNING id, user_id, tryout_id, status, started_at, expires_at, submitted_at, last_synced_at, version, total_questions, answered_questions, flagged_questions, correct_count, wrong_count, unanswered_count, final_score, created_at, updated_at
	`
	return scanAttempt(db.QueryRow(ctx, query, userID, tryoutID, expiresAt, totalQuestions))
}

func (r *Repository) ListQuestions(ctx context.Context, db store.DBTX, tryoutID string) ([]Question, error) {
	query := `
		SELECT
			q.id,
			q.code,
			q.question_type,
			q.prompt_html,
			COALESCE(q.image_url, ''),
			q.display_order,
			q.points,
			COALESCE(qo.option_key, ''),
			COALESCE(qo.option_text, ''),
			COALESCE(qo.display_order, 0)
		FROM questions q
		LEFT JOIN question_options qo ON qo.question_id = q.id
		WHERE q.tryout_id = $1
		ORDER BY q.display_order ASC, qo.display_order ASC
	`
	rows, err := db.Query(ctx, query, tryoutID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	orderedIDs := []string{}
	questionMap := map[string]*Question{}
	for rows.Next() {
		var questionID, code, questionType, promptHTML, imageURL, optionKey, optionText string
		var displayOrder, optionOrder int
		var points float64
		if err := rows.Scan(&questionID, &code, &questionType, &promptHTML, &imageURL, &displayOrder, &points, &optionKey, &optionText, &optionOrder); err != nil {
			return nil, err
		}
		question, exists := questionMap[questionID]
		if !exists {
			question = &Question{
				ID:           questionID,
				Code:         code,
				QuestionType: questionType,
				PromptHTML:   promptHTML,
				ImageURL:     imageURL,
				DisplayOrder: displayOrder,
				Points:       points,
				Options:      []QuestionOption{},
			}
			questionMap[questionID] = question
			orderedIDs = append(orderedIDs, questionID)
		}
		if optionKey != "" {
			question.Options = append(question.Options, QuestionOption{Key: optionKey, Text: optionText, DisplayOrder: optionOrder})
		}
	}
	questions := make([]Question, 0, len(orderedIDs))
	for _, id := range orderedIDs {
		questions = append(questions, *questionMap[id])
	}
	return questions, rows.Err()
}

func (r *Repository) ListAnswers(ctx context.Context, db store.DBTX, attemptID string) ([]AttemptAnswer, error) {
	query := `
		SELECT id, attempt_id, question_id, selected_option_key, answer_text, is_flagged, answered_at, last_saved_at, client_updated_at, server_version, is_correct, awarded_points
		FROM attempt_answers
		WHERE attempt_id = $1
		ORDER BY created_at ASC
	`
	rows, err := db.Query(ctx, query, attemptID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	answers := []AttemptAnswer{}
	for rows.Next() {
		var item AttemptAnswer
		var selectedOptionKey sql.NullString
		var answerText sql.NullString
		var answeredAt sql.NullTime
		var clientUpdatedAt sql.NullTime
		var isCorrect sql.NullBool
		var awardedPoints sql.NullFloat64
		if err := rows.Scan(
			&item.ID,
			&item.AttemptID,
			&item.QuestionID,
			&selectedOptionKey,
			&answerText,
			&item.IsFlagged,
			&answeredAt,
			&item.LastSavedAt,
			&clientUpdatedAt,
			&item.ServerVersion,
			&isCorrect,
			&awardedPoints,
		); err != nil {
			return nil, err
		}
		if selectedOptionKey.Valid {
			item.SelectedOptionKey = &selectedOptionKey.String
		}
		if answerText.Valid {
			item.AnswerText = &answerText.String
		}
		if answeredAt.Valid {
			item.AnsweredAt = &answeredAt.Time
		}
		if clientUpdatedAt.Valid {
			item.ClientUpdatedAt = &clientUpdatedAt.Time
		}
		if isCorrect.Valid {
			value := isCorrect.Bool
			item.IsCorrect = &value
		}
		if awardedPoints.Valid {
			value := awardedPoints.Float64
			item.AwardedPoints = &value
		}
		answers = append(answers, item)
	}
	return answers, rows.Err()
}

func (r *Repository) GetQuestionMetaByIDs(ctx context.Context, db store.DBTX, tryoutID string, questionIDs []string) (map[string]QuestionMeta, error) {
	query := `
		SELECT q.id, q.question_type, COALESCE(array_agg(qo.option_key) FILTER (WHERE qo.option_key IS NOT NULL), '{}')
		FROM questions q
		LEFT JOIN question_options qo ON qo.question_id = q.id
		WHERE q.tryout_id = $1 AND q.id = ANY($2)
		GROUP BY q.id, q.question_type
	`
	rows, err := db.Query(ctx, query, tryoutID, questionIDs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	metas := map[string]QuestionMeta{}
	for rows.Next() {
		var questionID, questionType string
		var optionKeys []string
		if err := rows.Scan(&questionID, &questionType, &optionKeys); err != nil {
			return nil, err
		}
		meta := QuestionMeta{ID: questionID, QuestionType: questionType, ValidOptionKeys: map[string]bool{}}
		for _, optionKey := range optionKeys {
			meta.ValidOptionKeys[optionKey] = true
		}
		metas[questionID] = meta
	}
	return metas, rows.Err()
}

func (r *Repository) UpsertAnswer(ctx context.Context, db store.DBTX, attemptID, questionID string, selectedOptionKey, answerText, normalizedAnswer *string, isFlagged bool, answeredAt, clientUpdatedAt *time.Time) error {
	query := `
		INSERT INTO attempt_answers (attempt_id, question_id, selected_option_key, answer_text, normalized_answer, is_flagged, answered_at, last_saved_at, client_updated_at, is_correct, awarded_points)
		VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), $8, NULL, NULL)
		ON CONFLICT (attempt_id, question_id)
		DO UPDATE SET
			selected_option_key = EXCLUDED.selected_option_key,
			answer_text = EXCLUDED.answer_text,
			normalized_answer = EXCLUDED.normalized_answer,
			is_flagged = EXCLUDED.is_flagged,
			answered_at = EXCLUDED.answered_at,
			last_saved_at = EXCLUDED.last_saved_at,
			client_updated_at = EXCLUDED.client_updated_at,
			is_correct = NULL,
			awarded_points = NULL,
			server_version = attempt_answers.server_version + 1
		WHERE attempt_answers.client_updated_at IS NULL OR EXCLUDED.client_updated_at >= attempt_answers.client_updated_at
	`
	_, err := db.Exec(ctx, query, attemptID, questionID, selectedOptionKey, answerText, normalizedAnswer, isFlagged, answeredAt, clientUpdatedAt)
	return err
}

func (r *Repository) RefreshAttemptStats(ctx context.Context, db store.DBTX, attemptID string) (Attempt, error) {
	query := `
		WITH stats AS (
			SELECT
				COUNT(*) FILTER (WHERE selected_option_key IS NOT NULL OR COALESCE(normalized_answer, '') <> '')::int AS answered_questions,
				COUNT(*) FILTER (WHERE is_flagged)::int AS flagged_questions
			FROM attempt_answers
			WHERE attempt_id = $1
		)
		UPDATE attempts a
		SET
			answered_questions = COALESCE(stats.answered_questions, 0),
			flagged_questions = COALESCE(stats.flagged_questions, 0),
			unanswered_count = GREATEST(a.total_questions - COALESCE(stats.answered_questions, 0), 0),
			last_synced_at = NOW(),
			version = a.version + 1
		FROM stats
		WHERE a.id = $1
		RETURNING a.id, a.user_id, a.tryout_id, a.status, a.started_at, a.expires_at, a.submitted_at, a.last_synced_at, a.version, a.total_questions, a.answered_questions, a.flagged_questions, a.correct_count, a.wrong_count, a.unanswered_count, a.final_score, a.created_at, a.updated_at
	`
	return scanAttempt(db.QueryRow(ctx, query, attemptID))
}

func (r *Repository) GetScoringQuestions(ctx context.Context, db store.DBTX, tryoutID string) ([]ScoringQuestion, error) {
	query := `
		SELECT q.id, q.question_type, q.points, COALESCE(qak.correct_option_key, ''), COALESCE(array_agg(qsav.normalized_text) FILTER (WHERE qsav.id IS NOT NULL), '{}')
		FROM questions q
		LEFT JOIN question_answer_keys qak ON qak.question_id = q.id
		LEFT JOIN question_short_answer_variants qsav ON qsav.question_id = q.id
		WHERE q.tryout_id = $1
		GROUP BY q.id, q.question_type, q.points, qak.correct_option_key, q.display_order
		ORDER BY q.display_order ASC
	`
	rows, err := db.Query(ctx, query, tryoutID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []ScoringQuestion{}
	for rows.Next() {
		var item ScoringQuestion
		var variants []string
		if err := rows.Scan(&item.ID, &item.QuestionType, &item.Points, &item.CorrectOptionKey, &variants); err != nil {
			return nil, err
		}
		item.Variants = map[string]bool{}
		for _, variant := range variants {
			item.Variants[variant] = true
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *Repository) GetScoringAnswers(ctx context.Context, db store.DBTX, attemptID string) (map[string]ScoringAnswer, error) {
	query := `
		SELECT question_id, selected_option_key, normalized_answer, is_flagged
		FROM attempt_answers
		WHERE attempt_id = $1
	`
	rows, err := db.Query(ctx, query, attemptID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := map[string]ScoringAnswer{}
	for rows.Next() {
		var questionID string
		var selectedOptionKey sql.NullString
		var normalizedAnswer sql.NullString
		var isFlagged bool
		if err := rows.Scan(&questionID, &selectedOptionKey, &normalizedAnswer, &isFlagged); err != nil {
			return nil, err
		}
		item := ScoringAnswer{QuestionID: questionID, IsFlagged: isFlagged}
		if selectedOptionKey.Valid {
			item.SelectedOptionKey = &selectedOptionKey.String
		}
		if normalizedAnswer.Valid {
			item.NormalizedAnswer = &normalizedAnswer.String
		}
		items[questionID] = item
	}
	return items, rows.Err()
}

func (r *Repository) UpdateAnswerScore(ctx context.Context, db store.DBTX, attemptID, questionID string, isCorrect *bool, awardedPoints float64) error {
	query := `
		UPDATE attempt_answers
		SET is_correct = $3, awarded_points = $4
		WHERE attempt_id = $1 AND question_id = $2
	`
	_, err := db.Exec(ctx, query, attemptID, questionID, isCorrect, awardedPoints)
	return err
}

func (r *Repository) FinalizeAttempt(ctx context.Context, db store.DBTX, attemptID, status string, totalQuestions, answeredQuestions, correctCount, wrongCount, unansweredCount int, finalScore float64) (Attempt, error) {
	query := `
		UPDATE attempts
		SET
			status = $2,
			submitted_at = NOW(),
			last_synced_at = NOW(),
			version = version + 1,
			total_questions = $3,
			answered_questions = $4,
			flagged_questions = COALESCE((SELECT COUNT(*) FROM attempt_answers WHERE attempt_id = $1 AND is_flagged), 0),
			correct_count = $5,
			wrong_count = $6,
			unanswered_count = $7,
			final_score = $8
		WHERE id = $1
		RETURNING id, user_id, tryout_id, status, started_at, expires_at, submitted_at, last_synced_at, version, total_questions, answered_questions, flagged_questions, correct_count, wrong_count, unanswered_count, final_score, created_at, updated_at
	`
	return scanAttempt(db.QueryRow(ctx, query, attemptID, status, totalQuestions, answeredQuestions, correctCount, wrongCount, unansweredCount, finalScore))
}

func (r *Repository) GetResultByUserTryout(ctx context.Context, db store.DBTX, userID, tryoutID string) (ResultSummary, error) {
	query := `
		WITH ranked AS (
			SELECT
				a.id,
				a.user_id,
				a.status,
				a.submitted_at,
				a.total_questions,
				a.correct_count,
				a.wrong_count,
				a.unanswered_count,
				COALESCE(a.final_score, 0) AS final_score,
				ROW_NUMBER() OVER (ORDER BY a.final_score DESC NULLS LAST, a.submitted_at ASC NULLS LAST, a.id ASC) AS rank,
				COUNT(*) OVER () AS total_participants
			FROM attempts a
			WHERE a.tryout_id = $1 AND a.status IN ('submitted', 'auto_submitted')
		)
		SELECT id, status, submitted_at, total_questions, correct_count, wrong_count, unanswered_count, final_score, rank, total_participants
		FROM ranked
		WHERE user_id = $2
		LIMIT 1
	`
	var out ResultSummary
	var submittedAt sql.NullTime
	if err := db.QueryRow(ctx, query, tryoutID, userID).Scan(
		&out.AttemptID,
		&out.Status,
		&submittedAt,
		&out.TotalQuestions,
		&out.CorrectCount,
		&out.WrongCount,
		&out.UnansweredCount,
		&out.FinalScore,
		&out.Rank,
		&out.TotalParticipants,
	); err != nil {
		return ResultSummary{}, err
	}
	if submittedAt.Valid {
		out.SubmittedAt = &submittedAt.Time
	}
	return out, nil
}

func (r *Repository) ListLeaderboard(ctx context.Context, db store.DBTX, tryoutID string, limit, offset int) ([]LeaderboardEntry, int, error) {
	query := `
		WITH ranked AS (
			SELECT
				a.user_id,
				u.full_name,
				COALESCE(a.final_score, 0) AS final_score,
				a.submitted_at,
				ROW_NUMBER() OVER (ORDER BY a.final_score DESC NULLS LAST, a.submitted_at ASC NULLS LAST, a.id ASC) AS rank,
				COUNT(*) OVER () AS total_rows
			FROM attempts a
			JOIN users u ON u.id = a.user_id
			WHERE a.tryout_id = $1 AND a.status IN ('submitted', 'auto_submitted')
		)
		SELECT rank, user_id, full_name, final_score, submitted_at, total_rows
		FROM ranked
		ORDER BY rank ASC
		LIMIT $2 OFFSET $3
	`
	rows, err := db.Query(ctx, query, tryoutID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	entries := []LeaderboardEntry{}
	total := 0
	for rows.Next() {
		var entry LeaderboardEntry
		var submittedAt sql.NullTime
		if err := rows.Scan(&entry.Rank, &entry.UserID, &entry.FullName, &entry.FinalScore, &submittedAt, &total); err != nil {
			return nil, 0, err
		}
		if submittedAt.Valid {
			entry.SubmittedAt = &submittedAt.Time
		}
		entries = append(entries, entry)
	}
	return entries, total, rows.Err()
}

func scanAttempt(row pgx.Row) (Attempt, error) {
	var out Attempt
	var submittedAt sql.NullTime
	var lastSyncedAt sql.NullTime
	var finalScore sql.NullFloat64
	if err := row.Scan(
		&out.ID,
		&out.UserID,
		&out.TryoutID,
		&out.Status,
		&out.StartedAt,
		&out.ExpiresAt,
		&submittedAt,
		&lastSyncedAt,
		&out.Version,
		&out.TotalQuestions,
		&out.AnsweredQuestions,
		&out.FlaggedQuestions,
		&out.CorrectCount,
		&out.WrongCount,
		&out.UnansweredCount,
		&finalScore,
		&out.CreatedAt,
		&out.UpdatedAt,
	); err != nil {
		return Attempt{}, err
	}
	if submittedAt.Valid {
		out.SubmittedAt = &submittedAt.Time
	}
	if lastSyncedAt.Valid {
		out.LastSyncedAt = &lastSyncedAt.Time
	}
	if finalScore.Valid {
		value := finalScore.Float64
		out.FinalScore = &value
	}
	return out, nil
}
